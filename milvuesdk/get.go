package milvuesdk

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/rronan/gonetdicom/dicomutil"
	"github.com/rronan/gonetdicom/dicomweb"
	"github.com/suyashkumar/dicom"
)

func WaitDone(api_url, study_instance_uid string, token string, interval int, timeout int) (GetStudyStatusResponseV3, error) {
	t1 := time.Now().Add(time.Duration(timeout * 1e9))
	var status_response GetStudyStatusResponseV3
	for time.Now().Before(t1) {
		status_response, err := GetStatus(api_url, study_instance_uid, token)
		if err != nil {
			return GetStudyStatusResponseV3{}, err
		}
		if status_response.Status != "running" {
			return status_response, nil
		}
		time.Sleep(time.Duration(interval * 1e9))
	}
	return status_response, nil
}

func GetStatus(api_url, study_instance_uid string, token string) (GetStudyStatusResponseV3, error) {
	url := fmt.Sprintf("%s/v3/studies/%s/status", api_url, study_instance_uid)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return GetStudyStatusResponseV3{}, err
	}
	req.Header.Set("x-goog-meta-owner", token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return GetStudyStatusResponseV3{}, err
	}
	defer resp.Body.Close()
	status_response := GetStudyStatusResponseV3{}
	json.NewDecoder(resp.Body).Decode(&status_response)
	return status_response, nil
}

func Get(api_url, study_instance_uid string, inference_command string, token string) ([]*dicom.Dataset, error) {
	url := fmt.Sprintf(
		"%s/v3/studies/%s?inference_command=%s&signed_url=false",
		api_url,
		study_instance_uid,
		inference_command,
	)
	headers := map[string]string{"x-goog-meta-owner": token, "Content-Type": "multipart/related; type=application/dicom"}
	dcm_slice, _, err := dicomweb.Wado(url, headers)
	if err != nil {
		return []*dicom.Dataset{}, err
	}
	return dcm_slice, nil
}

func GetToFile(api_url, study_instance_uid string, inference_command string, token string, folder string) ([]string, error) {
	url := fmt.Sprintf(
		"%s/v3/studies/%s?inference_command=%s&signed_url=false",
		api_url,
		study_instance_uid,
		inference_command,
	)
	headers := map[string]string{"x-goog-meta-owner": token, "Content-Type": "multipart/related; type=application/dicom"}
	dcm_path_slice, _, err := dicomweb.WadoToFile(url, headers, folder)
	if err != nil {
		return []string{}, err
	}
	return dcm_path_slice, nil
}

func downloadSignedUrl(signed_url string, token string) (*dicom.Dataset, error) {
	req, err := http.NewRequest("GET", signed_url, nil)
	if err != nil {
		return &dicom.Dataset{}, err
	}
	req.Header.Set("x-goog-meta-owner", token)
	req.Header.Set("Content-Type", "application/dicom")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &dicom.Dataset{}, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &dicom.Dataset{}, err
	}
	return dicomutil.Bytes2Dicom(data)
}

func GetSignedUrl(api_url, study_instance_uid string, inference_command string, token string) ([]*dicom.Dataset, error) {
	url := fmt.Sprintf(
		"%s/v3/studies/%s?inference_command=%s&signed_url=true",
		api_url,
		study_instance_uid,
		inference_command,
	)
	headers := map[string]string{
		"x-goog-meta-owner": token,
		"Accept":            "application/json",
	}
	resp, err := dicomweb.Get(url, headers)
	if err != nil {
		return []*dicom.Dataset{}, err
	}
	defer resp.Body.Close()
	get_response := GetStudyResponseV3{}
	json.NewDecoder(resp.Body).Decode(&get_response)
	if get_response.SignedUrls == nil || len(*get_response.SignedUrls) == 0 {
		return []*dicom.Dataset{}, nil
	}
	dcm_slice := []*dicom.Dataset{}
	for _, signed_url := range *get_response.SignedUrls {
		dcm, err := downloadSignedUrl(signed_url, token)
		if err != nil {
			return []*dicom.Dataset{}, err
		}
		dcm_slice = append(dcm_slice, dcm)
	}
	get_response.SignedUrls = nil
	return dcm_slice, nil
}
func downloadSignedUrlToFile(signed_url string, token string, dcm_path string) error {
	req, err := http.NewRequest("GET", signed_url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("x-goog-meta-owner", token)
	req.Header.Set("Content-Type", "application/dicom")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	f, err := os.Create(dcm_path)
	if err != nil {
		return err
	}
	defer f.Close()
	io.Copy(f, resp.Body)
	return nil
}

func GetSignedUrlToFile(api_url, study_instance_uid string, inference_command string, token string, folder string) ([]string, error) {
	res := []string{}
	url := fmt.Sprintf(
		"%s/v3/studies/%s?inference_command=%s&signed_url=true",
		api_url,
		study_instance_uid,
		inference_command,
	)
	headers := map[string]string{
		"x-goog-meta-owner": token,
		"Accept":            "application/json",
	}
	resp, err := dicomweb.Get(url, headers)
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()
	get_response := GetStudyResponseV3{}
	json.NewDecoder(resp.Body).Decode(&get_response)
	if get_response.SignedUrls == nil || len(*get_response.SignedUrls) == 0 {
		return res, nil
	}
	for _, signed_url := range *get_response.SignedUrls {
		dcm_path := fmt.Sprintf("%s/%x", folder, dicomutil.RandomDicomName())
		err := downloadSignedUrlToFile(signed_url, token, dcm_path)
		if err != nil {
			return res, err
		}
		res = append(res, dcm_path)
	}
	return res, nil
}

func GetSmarturgences(api_url, study_instance_uid string, token string) (GetSmarturgencesResponseV3, error) {
	url := fmt.Sprintf("%s/v3/smarturgences/%s", api_url, study_instance_uid)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return GetSmarturgencesResponseV3{}, err
	}
	req.Header.Set("x-goog-meta-owner", token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return GetSmarturgencesResponseV3{}, err
	}
	defer resp.Body.Close()
	smarturgences_response := GetSmarturgencesResponseV3{}
	json.NewDecoder(resp.Body).Decode(&smarturgences_response)
	return smarturgences_response, nil
}

func GetSmartxpert(api_url, study_instance_uid string, token string) (GetSmartxpertResponseV3, error) {
	url := fmt.Sprintf("%s/v3/smartxpert/%s", api_url, study_instance_uid)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return GetSmartxpertResponseV3{}, err
	}
	req.Header.Set("x-goog-meta-owner", token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return GetSmartxpertResponseV3{}, err
	}
	defer resp.Body.Close()
	smartxpert_response := GetSmartxpertResponseV3{}
	json.NewDecoder(resp.Body).Decode(&smartxpert_response)
	return smartxpert_response, nil
}

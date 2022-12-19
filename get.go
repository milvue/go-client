package goclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/rronan/gonetdicom/dicomutil"
	"github.com/rronan/gonetdicom/dicomweb"
	"github.com/suyashkumar/dicom"
)

func Get(api_url, study_instance_uid string, inference_command string, token string) ([]*dicom.Dataset, GetResponse, error) {
	url := fmt.Sprintf(
		"%s/v3/studies/%s?inference_command=%s&signed_url=false",
		api_url,
		study_instance_uid,
		inference_command,
	)
	headers := map[string]string{
		"x-goog-meta-owner": token,
		"Content-Type":      "multipart/related; type=application/dicom",
		"Accept":            "multipart/related, application/json",
	}
	dcm_slice, resp, err := dicomweb.Wado(url, headers)
	if err != nil {
		return []*dicom.Dataset{}, GetResponse{}, err
	}
	var get_response GetResponse
	json.NewDecoder(resp.Body).Decode(&get_response)
	return dcm_slice, get_response, nil
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

func GetSignedUrl(api_url, study_instance_uid string, inference_command string, token string) ([]*dicom.Dataset, GetResponse, error) {
	url := fmt.Sprintf(
		"%s/v3/studies/%s?inference_command=%s&signed_url=true",
		api_url,
		study_instance_uid,
		inference_command,
	)
	headers := map[string]string{
		"x-goog-meta-owner": token,
		"Content-Type":      "application/json",
		"Accept":            "application/json",
	}
	resp, err := dicomweb.Get(url, headers)
	if err != nil {
		return []*dicom.Dataset{}, GetResponse{}, err
	}
	defer resp.Body.Close()
	get_response := GetResponse{}
	json.NewDecoder(resp.Body).Decode(&get_response)
	dcm_slice := []*dicom.Dataset{}
	for _, signed_url := range get_response.SignedUrls {
		dcm, err := downloadSignedUrl(signed_url, token)
		if err != nil {
			return []*dicom.Dataset{}, get_response, err
		}
		dcm_slice = append(dcm_slice, dcm)
	}
	get_response.SignedUrls = nil
	return dcm_slice, get_response, nil
}

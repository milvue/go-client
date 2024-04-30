package milvuesdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/rronan/gonetdicom/dicomutil"
	"github.com/rronan/gonetdicom/dicomweb"
	"github.com/suyashkumar/dicom"
)

var MILVUE_API_URL = getenv("MILVUE_API_URL", "")

func Post(api_url string, dcm_slice []*dicom.Dataset, token string, client_timeout int) error {
	url := api_url + "/v3/studies?signed_url=false"
	headers := map[string]string{
		"x-goog-meta-owner": token,
		"Content-Type":      "multipart/related; type=application/dicom",
		"Accept":            "application/json",
	}
	resp, err := dicomweb.Stow(url, dcm_slice, headers, client_timeout)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func PostFromFile(api_url string, dcm_path_slice []string, token string, client_timeout int) error {
	url := api_url + "/v3/studies?signed_url=false"
	headers := map[string]string{
		"x-goog-meta-owner": token,
		"Content-Type":      "multipart/related; type=application/dicom",
		"Accept":            "application/json",
	}
	resp, err := dicomweb.StowFromFile(url, dcm_path_slice, headers, client_timeout)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func PostSignedUrl(api_url string, dcm_slice []*dicom.Dataset, token string, client_timeout int) error {
	url := fmt.Sprintf("%s/v3/studies?signed_url=true", api_url)
	headers := map[string]string{
		"x-goog-meta-owner": token,
		"Content-Type":      "application/dicom",
		"Accept":            "application/json",
	}
	r, err := dicomweb.Stow(url, pruneDicomSlice(dcm_slice), headers, client_timeout)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	post_signed_url_response := PostSignedUrlResponseV3{}
	json.NewDecoder(r.Body).Decode(&post_signed_url_response)
	for _, dcm := range dcm_slice {
		_, _, sop_instance_uid, err := dicomutil.GetUIDs(dcm)
		if err != nil {
			return err
		}
		// don't know why but I need to redeclare this
		headers = map[string]string{
			"x-goog-meta-owner": token,
			"Content-Type":      "application/dicom",
		}
		signed_url := post_signed_url_response.SignedUrls[sop_instance_uid]
		err = dicomweb.Put(signed_url, dcm, headers, client_timeout)
		if err != nil {
			return err
		}
	}
	return nil
}

func PostSignedUrlFromFile(api_url string, dcm_path_slice []string, token string, client_timeout int) error {
	url := fmt.Sprintf("%s/v3/studies?signed_url=true", api_url)
	headers := map[string]string{
		"x-goog-meta-owner": token,
		"Content-Type":      "application/dicom",
		"Accept":            "application/json",
	}
	dcm_slice := []*dicom.Dataset{}
	for _, dcm_path := range dcm_path_slice {
		dcm, err := dicom.ParseFile(dcm_path, nil, dicom.SkipPixelData())
		if err != nil {
			return err
		}
		dcm_slice = append(dcm_slice, &dcm)
	}
	resp, err := dicomweb.Stow(url, pruneDicomSlice(dcm_slice), headers, client_timeout)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	post_signed_url_response := PostSignedUrlResponseV3{}
	json.NewDecoder(resp.Body).Decode(&post_signed_url_response)
	for i, dcm_path := range dcm_path_slice {
		_, _, sop_instance_uid, err := dicomutil.GetUIDs(dcm_slice[i])
		if err != nil {
			return err
		}
		// don't know why but I need to redeclare this
		headers = map[string]string{
			"x-goog-meta-owner": token,
			"Content-Type":      "application/dicom",
		}
		signed_url := post_signed_url_response.SignedUrls[sop_instance_uid]
		err = dicomweb.PutFromFile(signed_url, dcm_path, headers, client_timeout)
		if err != nil {
			return err
		}
	}
	return nil
}

func PostInteresting(api_url string, study_instance_uid, token string, client_timeout int) (*http.Response, error) {
	url := api_url + "/v3/interesting"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return &http.Response{}, err
	}
	req.Header.Set("x-goog-meta-owner", token)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: time.Duration(client_timeout * 1e9)}
	resp, err := client.Do(req)
	if err != nil {
		return &http.Response{}, err
	}
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		return &http.Response{}, &dicomweb.RequestError{StatusCode: resp.StatusCode, Err: errors.New(resp.Status)}
	}
	return resp, nil
}

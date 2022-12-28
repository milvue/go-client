package milvuesdk

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rronan/gonetdicom/dicomutil"
	"github.com/rronan/gonetdicom/dicomweb"
	"github.com/suyashkumar/dicom"
)

var MILVUE_API_URL = getenv("MILVUE_API_URL", "")

func Post(api_url string, dcm_slice []*dicom.Dataset, token string) error {
	url := api_url + "/v3/studies?signed_url=false"
	headers := map[string]string{
		"x-goog-meta-owner": token,
		"Content-Type":      "multipart/related; type=application/dicom",
		"Accept":            "application/json",
	}
	_, err := dicomweb.Post(url, dcm_slice, headers)
	return err
}

func PostSignedUrl(api_url string, dcm_slice []*dicom.Dataset, token string) error {
	url := fmt.Sprintf("%s/v3/studies?signed_url=true", api_url)
	pruned_dcm_slice, err := pruneDicoms(dcm_slice)
	if err != nil {
		return err
	}
	headers := map[string]string{
		"x-goog-meta-owner": token,
		"Content-Type":      "application/dicom",
		"Accept":            "application/json",
	}
	r, err := dicomweb.Post(url, pruned_dcm_slice, headers)
	defer r.Body.Close()
	post_signed_url_response := PostSignedUrlResponseV3{}
	json.NewDecoder(r.Body).Decode(&post_signed_url_response)
	for _, dcm := range dcm_slice {
		_, sop_instance_uid, err := dicomutil.GetUIDs(dcm)
		if err != nil {
			return err
		}
		// don't know why but I need to redeclare this
		headers = map[string]string{
			"x-goog-meta-owner": token,
			"Content-Type":      "application/dicom",
		}
		signed_url := post_signed_url_response.SignedUrls[sop_instance_uid]
		err = dicomweb.Put(signed_url, dcm, headers)
		if err != nil {
			return err
		}
	}
	return nil
}

func PostInteresting(api_url string, study_instance_uid, token string) (*http.Response, error) {
	url := api_url + "/v3/interesting"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return &http.Response{}, err
	}
	req.Header.Set("x-goog-meta-owner", token)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &http.Response{}, err
	}
	return resp, nil
}

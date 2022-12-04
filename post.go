package goclient

import (
	"fmt"

	"github.com/rronan/gonetdicom/dicomweb"
	"github.com/suyashkumar/dicom"
)

var MILVUE_API_URL = getenv("MILVUE_API_URL", "")

func Post(api_url string, dcm_slice []*dicom.Dataset, token string) error {
	url := api_url + "/v3/studies?signed_url=false"
	headers := map[string]string{"x-goog-meta-owner": token, "Content-Type": "multipart/related; type=application/dicom"}
	return dicomweb.Post(url, dcm_slice, headers)
}

func PostSignedUrl(api_url string, dcm_slice []*dicom.Dataset, token string) error {
	url := fmt.Sprintf("%s/v3/studies?signed_url=true", api_url)
	pruned_dcm_slice, err := pruneDicoms(dcm_slice)
	if err != nil {
		return err
	}
	headers := map[string]string{"x-goog-meta-owner": token, "Content-Type": "application/json"}
	dicomweb.Post(url, pruned_dcm_slice, headers)
	signed_url := "should come from Post"
	headers["Content-Type"] = "application/dicom"
	for _, dcm := range dcm_slice {
		err = dicomweb.Put(signed_url, dcm, headers)
		if err != nil {
			return err
		}
	}
	return nil
}

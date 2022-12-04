package goclient

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/rronan/gonetdicom/dicomutil"
	"github.com/rronan/gonetdicom/dicomweb"
	"github.com/suyashkumar/dicom"
)

func Get(api_url, study_instance_uid string, token string) ([]*dicom.Dataset, []byte, error) {
	url := fmt.Sprintf("%s/v3/studies/%s?signed_url=false", api_url, study_instance_uid)
	headers := map[string]string{"x-goog-meta-owner": token, "Content-Type": "multipart/related; type=application/dicom"}
	dcm_slice, msg, err := dicomweb.GetDicomWeb(url, headers)
	if err != nil {
		return []*dicom.Dataset{}, []byte{}, err
	}
	// Unmarshal msg as open, see https://github.com/deepmap/oapi-codegen
	return dcm_slice, msg, nil
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

func GetSignedUrl(api_url, study_instance_uid string, token string) ([]*dicom.Dataset, []byte, error) {
	url := fmt.Sprintf("%s/v3/studies/%s?signed_url=true", api_url, study_instance_uid)
	headers := map[string]string{"x-goog-meta-owner": token, "Content-Type": "application/json"}
	_, msg, err := dicomweb.GetDicomWeb(url, headers)
	if err != nil {
		return []*dicom.Dataset{}, []byte{}, err
	}
	// Unmarshal msg as open, see https://github.com/deepmap/oapi-codegen
	signed_url_slice := []string{"foo", "bar"}
	dcm_slice := []*dicom.Dataset{}
	for _, signed_url := range signed_url_slice {
		dcm, err := downloadSignedUrl(signed_url, token)
		if err != nil {
			return []*dicom.Dataset{}, []byte{}, err
		}
		dcm_slice = append(dcm_slice, dcm)
	}
	return dcm_slice, msg, nil
}

package goclient

import (
	"github.com/rronan/gonetdicom/dicomutil"
	"github.com/rronan/gonetdicom/dicomweb"
	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/tag"
)

func Post(url string, dcm_slice []*dicom.Dataset, token string) error {
	headers := map[string]string{"x-goog-meta-owner": token, "Content-Type": "multipart/related; type=application/dicom"}
	return dicomweb.Post(url, dcm_slice, headers)
}

func pruneDicoms(dcm_slice []*dicom.Dataset) ([]*dicom.Dataset, error) {
	res := []*dicom.Dataset{}
	for _, dcm := range dcm_slice {
		pruned_dcm := dcm
		element, err := pruned_dcm.FindElementByTag(tag.PixelData)
		if err != nil {
			return []*dicom.Dataset{}, err
		}
		*element = *dicomutil.NULL_PIXEL
		res = append(res, pruned_dcm)
	}
	return res, nil
}

func PostSignedUrl(url string, dcm_slice []*dicom.Dataset, token string) error {
	pruned_dcm_slice, err := pruneDicoms(dcm_slice)
	if err != nil {
		return err
	}
	headers := map[string]string{"x-goog-meta-owner": token, "Content-Type": "multipart/related; type=application/dicom"}
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

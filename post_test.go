package goclient

import (
	"testing"

	"github.com/suyashkumar/dicom"
)

func Test_Post(t *testing.T) {
	dcm_slice := []*dicom.Dataset{}
	for _, path := range DICOM_PATH_SLICE {
		dcm, err := dicom.ParseFile(path, nil)
		if err != nil {
			panic(err)
		}
		dcm_slice = append(dcm_slice, &dcm)
	}
	err := Post(API_URL, dcm_slice, TOKEN)
	if err != nil {
		panic(err)
	}
}

func Test_PostSignedUrl(t *testing.T) {
	dcm_slice := []*dicom.Dataset{}
	for _, path := range DICOM_PATH_SLICE {
		dcm, err := dicom.ParseFile(path, nil)
		if err != nil {
			panic(err)
		}
		dcm_slice = append(dcm_slice, &dcm)
	}
	err := Post(API_URL, dcm_slice, TOKEN)
	if err != nil {
		panic(err)
	}
}

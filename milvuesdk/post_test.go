package milvuesdk

import (
	"fmt"
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

func Test_PostUrlFromFile(t *testing.T) {
	err := PostFromFile(API_URL, DICOM_PATH_SLICE, TOKEN)
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
	err := PostSignedUrl(API_URL, dcm_slice, TOKEN)
	if err != nil {
		panic(err)
	}
}

func Test_PostSignedUrlFromFile(t *testing.T) {
	err := PostSignedUrlFromFile(API_URL, DICOM_PATH_SLICE, TOKEN)
	if err != nil {
		panic(err)
	}
}

func Test_PostInteresting(t *testing.T) {
	status_code, err := PostInteresting(API_URL, StudyInstanceUID, TOKEN)
	if err != nil {
		panic(err)
	}
	fmt.Println(status_code)
}

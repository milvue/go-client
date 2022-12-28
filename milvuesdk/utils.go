package milvuesdk

import (
	"os"

	"github.com/rronan/gonetdicom/dicomutil"
	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/tag"
)

var API_URL = getenv("MILVUE_API_URL", "")
var TOKEN = getenv("MILVUE_TOKEN", "")
var StudyInstanceUID = "1.2.826.0.1.3680044.0.0.0.20221228121333.16387"
var DICOM_PATH_SLICE = []string{
	"../data/study/1.2.276.0.7230010.3.1.4.0.78767.1672226121.633599.dcm",
	"../data/study/1.2.276.0.7230010.3.1.4.0.78767.1672226121.633601.dcm",
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
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

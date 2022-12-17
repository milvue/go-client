package goclient

import (
	"os"

	"github.com/rronan/gonetdicom/dicomutil"
	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/tag"
)

var API_URL = getenv("MILVUE_API_URL", "")
var TOKEN = getenv("MILVUE_TOKEN", "")
var DICOM_PATH_SLICE = []string{
	"data/study/DX.1.2.840.113564.1721633113.20201101205936468650.3003000225002.dcm",
	"data/study/DX.1.2.840.113564.1721633113.20201101205936484660.3003000225002.dcm",
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

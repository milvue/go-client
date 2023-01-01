package milvuesdk

import (
	"errors"
	"testing"

	"github.com/rronan/gonetdicom/dicomutil"
	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/tag"
)

func Test_pruneDicom(t *testing.T) {
	dcm, err := dicom.ParseFile(DICOM_PATH_SLICE[0], nil)
	if err != nil {
		panic(err)
	}
	pruned_dcm := pruneDicom(&dcm)
	study_instance_uid, sop_instance_uid, err := dicomutil.GetUIDs(&dcm)
	if err != nil {
		panic(err)
	}
	p_study_instance_uid, p_sop_instance_uid, err := dicomutil.GetUIDs(&pruned_dcm)
	if err != nil {
		panic(err)
	}
	if study_instance_uid != p_study_instance_uid {
		panic(errors.New("Inconsistent study_instance_uid"))
	}
	if sop_instance_uid != p_sop_instance_uid {
		panic(errors.New("Inconsistent study_instance_uid"))
	}
	_ = dicomutil.Dicom2Bytes(&dcm)
	_ = dicomutil.Dicom2Bytes(&pruned_dcm)
	e, err := dcm.FindElementByTag(tag.PixelData)
	if err != nil {
		panic(err)
	}
	if e.ValueLength == 0 {
		panic(errors.New("dcm.PixelData has ValueLength == 0"))
	}
	pruned_e, err := pruned_dcm.FindElementByTag(tag.PixelData)
	if err != nil {
		panic(err)
	}
	if pruned_e.ValueLength > 0 {
		panic(errors.New("pruned_dcm.PixelData has ValueLength > 0"))
	}
}

func Test_pruneDicomSlice(t *testing.T) {
	dcm_slice := []*dicom.Dataset{}
	for _, path := range DICOM_PATH_SLICE {
		dcm, err := dicom.ParseFile(path, nil)
		if err != nil {
			panic(err)
		}
		dcm_slice = append(dcm_slice, &dcm)
	}
	pruned_dcm_slice := pruneDicomSlice(dcm_slice)
	for _, pruned_dcm := range pruned_dcm_slice {
		_ = dicomutil.Dicom2Bytes(pruned_dcm)
	}
}

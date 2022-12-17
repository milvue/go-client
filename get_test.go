package goclient

import (
	"fmt"
	"testing"

	"github.com/rronan/gonetdicom/dicomutil"
)

func Test_Get(t *testing.T) {
	for _, inference_command := range []string{"smarturgences", "smartxpert"} {
		fmt.Println(inference_command)
		dcm_slice, get_response, err := Get(API_URL, "1.2.840.113970.1.2.840.113970.6418804.20201101.1205635", inference_command, TOKEN)
		fmt.Println(get_response)
		if err != nil {
			panic(err)
		}
		for _, dcm := range dcm_slice {
			study_instance_uid, sop_instance_uid, err := dicomutil.GetUIDs(dcm)
			if err != nil {
				panic(err)
			}
			fmt.Printf("%s,%s\n", study_instance_uid, sop_instance_uid)
		}
	}
}

func Test_GetSignedUrl(t *testing.T) {
	for _, inference_command := range []string{"smarturgences", "smartxpert"} {
		fmt.Println(inference_command)
		dcm_slice, get_response, err := GetSignedUrl(API_URL, "1.2.840.113970.1.2.840.113970.6418804.20201101.1205635", inference_command, TOKEN)
		fmt.Println(get_response)
		if err != nil {
			panic(err)
		}
		for _, dcm := range dcm_slice {
			study_instance_uid, sop_instance_uid, err := dicomutil.GetUIDs(dcm)
			if err != nil {
				panic(err)
			}
			fmt.Printf("%s,%s\n", study_instance_uid, sop_instance_uid)
		}
	}
}

module github.com/milvue/go-client/milvuesdk

go 1.19

require (
	github.com/deepmap/oapi-codegen v1.12.4
	github.com/rronan/gonetdicom v0.0.0-20230218192732-0eaa6bda31ac
	github.com/suyashkumar/dicom v1.0.5
)

require (
	github.com/apapsch/go-jsonmerge/v2 v2.0.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	golang.org/x/text v0.4.0 // indirect
)

replace github.com/rronan/gonetdicom v0.0.0-20230218192732-0eaa6bda31ac => ../../gonetdicom

replace github.com/suyashkumar/dicom v1.0.5 => github.com/tcheever/suyashkumar-dicom v1.0.6-0.20220603162441-1fb3be6fbd88

module github.com/milvue/go-client/milvuesdk

go 1.19

require (
	github.com/deepmap/oapi-codegen v1.12.4
	github.com/rronan/gonetdicom v0.0.0-20230101161711-4dd179eefd7b
	github.com/suyashkumar/dicom v1.0.5
)

require (
	github.com/apapsch/go-jsonmerge/v2 v2.0.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	golang.org/x/text v0.4.0 // indirect
)

// this allows to modify gonetdicom locally to fix things
// replace github.com/rronan/gonetdicom => ../../gonetdicom

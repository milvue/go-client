module github.com/milvue/go-client/milvuesdk

go 1.21

require (
	github.com/deepmap/oapi-codegen v1.12.4
	github.com/rronan/gonetdicom v0.0.0-20240423141507-f0c0b6643c1e // TODO revise version when PR merged in gonetdicom
	github.com/suyashkumar/dicom v1.0.6
)

require (
	github.com/apapsch/go-jsonmerge/v2 v2.0.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	golang.org/x/text v0.7.0 // indirect
)

// replace github.com/rronan/gonetdicom v0.0.0-20231120170418-33702d88ae85 => ../../gonetdicom

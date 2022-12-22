# A Go Client for Milvue API v3

This package allows to POST studies and GET results to/from Milvue API (https://milvue.com)

It implements STOW and WADO protocols from the DICOMweb Standard, https://www.dicomstandard.org/using/dicomweb

In addition, structs are generated with from openapi schemas:

- `GetStatusResponseV3`: schema for *v3/studies/{study_instance_uid}/status*, the current prediction status
- `GetSmarturgencesResponseV3`: schema for *v3/smarturgences/{study_instance_uid}*, Smarturgences results
- `GetSmartxpertResponseV3`: schema for *v3/smartxpert/{study_instance_uid}*, Smartxpert results


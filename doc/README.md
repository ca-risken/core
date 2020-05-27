# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [finding/finding.proto](#finding/finding.proto)
    - [Finding](#core.finding.Finding)
  
- [finding/service.proto](#finding/service.proto)
    - [GetFindingRequest](#core.finding.GetFindingRequest)
    - [GetFindingResponse](#core.finding.GetFindingResponse)
    - [ListFindingRequest](#core.finding.ListFindingRequest)
    - [ListFindingResponse](#core.finding.ListFindingResponse)
  
    - [FindingService](#core.finding.FindingService)
  
- [iam/service.proto](#iam/service.proto)
    - [AuthnRequest](#core.iam.AuthnRequest)
    - [AuthnResponse](#core.iam.AuthnResponse)
    - [AuthzRequest](#core.iam.AuthzRequest)
    - [AuthzResponse](#core.iam.AuthzResponse)
  
    - [IAMService](#core.iam.IAMService)
  
- [iam/user.proto](#iam/user.proto)
    - [User](#core.iam.User)
  
- [Scalar Value Types](#scalar-value-types)



<a name="finding/finding.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## finding/finding.proto



<a name="core.finding.Finding"></a>

### Finding



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| fiding_id | [string](#string) |  |  |
| name | [string](#string) |  |  |
| data_source | [string](#string) |  |  |
| resource | [string](#string) |  |  |
| project_id | [string](#string) |  |  |
| data | [string](#string) |  |  |
| created_at | [int64](#int64) |  |  |
| updated_at | [int64](#int64) |  |  |





 

 

 

 



<a name="finding/service.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## finding/service.proto



<a name="core.finding.GetFindingRequest"></a>

### GetFindingRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| finding_id | [string](#string) |  |  |






<a name="core.finding.GetFindingResponse"></a>

### GetFindingResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| data | [Finding](#core.finding.Finding) |  |  |






<a name="core.finding.ListFindingRequest"></a>

### ListFindingRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [string](#string) |  |  |
| since | [string](#string) |  |  |






<a name="core.finding.ListFindingResponse"></a>

### ListFindingResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_ids | [string](#string) | repeated |  |





 

 

 


<a name="core.finding.FindingService"></a>

### FindingService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ListFinding | [ListFindingRequest](#core.finding.ListFindingRequest) | [ListFindingResponse](#core.finding.ListFindingResponse) |  |
| GetFinding | [GetFindingRequest](#core.finding.GetFindingRequest) | [GetFindingResponse](#core.finding.GetFindingResponse) |  |

 



<a name="iam/service.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## iam/service.proto



<a name="core.iam.AuthnRequest"></a>

### AuthnRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| token | [string](#string) |  |  |






<a name="core.iam.AuthnResponse"></a>

### AuthnResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user | [User](#core.iam.User) |  |  |






<a name="core.iam.AuthzRequest"></a>

### AuthzRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [string](#string) |  |  |
| action_ptn | [string](#string) |  |  |
| resource_ptn | [string](#string) |  |  |






<a name="core.iam.AuthzResponse"></a>

### AuthzResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ok | [bool](#bool) |  |  |





 

 

 


<a name="core.iam.IAMService"></a>

### IAMService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Authn | [AuthnRequest](#core.iam.AuthnRequest) | [AuthnResponse](#core.iam.AuthnResponse) |  |
| Authz | [AuthzRequest](#core.iam.AuthzRequest) | [AuthzResponse](#core.iam.AuthzResponse) |  |

 



<a name="iam/user.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## iam/user.proto



<a name="core.iam.User"></a>

### User



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [string](#string) |  |  |
| name | [string](#string) |  |  |
| activated | [bool](#bool) |  |  |
| created_at | [int64](#int64) |  |  |
| updated_at | [int64](#int64) |  |  |





 

 

 

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |


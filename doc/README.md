# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [finding/finding.proto](#finding/finding.proto)
    - [DeleteFindingRequest](#core.finding.DeleteFindingRequest)
    - [DeleteResourceRequest](#core.finding.DeleteResourceRequest)
    - [Empty](#core.finding.Empty)
    - [Finding](#core.finding.Finding)
    - [FindingForUpsert](#core.finding.FindingForUpsert)
    - [FindingTag](#core.finding.FindingTag)
    - [FindingTagForUpsert](#core.finding.FindingTagForUpsert)
    - [GetFindingRequest](#core.finding.GetFindingRequest)
    - [GetFindingResponse](#core.finding.GetFindingResponse)
    - [GetResourceRequest](#core.finding.GetResourceRequest)
    - [GetResourceResponse](#core.finding.GetResourceResponse)
    - [ListFindingRequest](#core.finding.ListFindingRequest)
    - [ListFindingResponse](#core.finding.ListFindingResponse)
    - [ListFindingTagRequest](#core.finding.ListFindingTagRequest)
    - [ListFindingTagResponse](#core.finding.ListFindingTagResponse)
    - [ListResourceRequest](#core.finding.ListResourceRequest)
    - [ListResourceResponse](#core.finding.ListResourceResponse)
    - [ListResourceTagRequest](#core.finding.ListResourceTagRequest)
    - [ListResourceTagResponse](#core.finding.ListResourceTagResponse)
    - [PutFindingRequest](#core.finding.PutFindingRequest)
    - [PutFindingResponse](#core.finding.PutFindingResponse)
    - [PutResourceRequest](#core.finding.PutResourceRequest)
    - [PutResourceResponse](#core.finding.PutResourceResponse)
    - [Resource](#core.finding.Resource)
    - [ResourceForUpsert](#core.finding.ResourceForUpsert)
    - [ResourceTag](#core.finding.ResourceTag)
    - [ResourceTagForUpsert](#core.finding.ResourceTagForUpsert)
    - [TagFindingRequest](#core.finding.TagFindingRequest)
    - [TagFindingResponse](#core.finding.TagFindingResponse)
    - [TagResourceRequest](#core.finding.TagResourceRequest)
    - [TagResourceResponse](#core.finding.TagResourceResponse)
    - [UntagFindingRequest](#core.finding.UntagFindingRequest)
    - [UntagResourceRequest](#core.finding.UntagResourceRequest)
  
    - [FindingService](#core.finding.FindingService)
  
- [iam/iam.proto](#iam/iam.proto)
    - [AuthnRequest](#core.iam.AuthnRequest)
    - [AuthnResponse](#core.iam.AuthnResponse)
    - [AuthzRequest](#core.iam.AuthzRequest)
    - [AuthzResponse](#core.iam.AuthzResponse)
    - [User](#core.iam.User)
  
    - [IAMService](#core.iam.IAMService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="finding/finding.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## finding/finding.proto



<a name="core.finding.DeleteFindingRequest"></a>

### DeleteFindingRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [uint32](#uint32) |  |  |
| finding_id | [uint64](#uint64) |  |  |






<a name="core.finding.DeleteResourceRequest"></a>

### DeleteResourceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [uint32](#uint32) |  |  |
| resource_id | [uint64](#uint64) |  |  |






<a name="core.finding.Empty"></a>

### Empty
Empty 空メッセージ






<a name="core.finding.Finding"></a>

### Finding
Finding エンティティ


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| finding_id | [uint64](#uint64) |  |  |
| description | [string](#string) |  |  |
| data_source | [string](#string) |  |  |
| resource_name | [string](#string) |  |  |
| project_id | [uint32](#uint32) |  |  |
| original_score | [float](#float) |  |  |
| original_max_score | [float](#float) |  |  |
| score | [float](#float) |  |  |
| data | [string](#string) |  |  |
| created_at | [int64](#int64) |  |  |
| updated_at | [int64](#int64) |  |  |






<a name="core.finding.FindingForUpsert"></a>

### FindingForUpsert
Finding エンティティ（登録・更新用）


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| description | [string](#string) |  |  |
| data_source | [string](#string) |  |  |
| resource_name | [string](#string) |  |  |
| project_id | [uint32](#uint32) |  |  |
| original_score | [float](#float) |  |  |
| original_max_score | [float](#float) |  |  |
| score | [float](#float) |  |  |
| data | [string](#string) |  |  |






<a name="core.finding.FindingTag"></a>

### FindingTag
FindingTag エンティティ


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| finding_tag_id | [uint64](#uint64) |  |  |
| finding_id | [uint64](#uint64) |  |  |
| tag_key | [string](#string) |  |  |
| tag_value | [string](#string) |  |  |
| created_at | [int64](#int64) |  |  |
| updated_at | [int64](#int64) |  |  |






<a name="core.finding.FindingTagForUpsert"></a>

### FindingTagForUpsert
FindingTag エンティティ（登録・更新用）


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| finding_id | [uint64](#uint64) |  |  |
| tag_key | [string](#string) |  |  |
| tag_value | [string](#string) |  |  |






<a name="core.finding.GetFindingRequest"></a>

### GetFindingRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [uint32](#uint32) |  |  |
| finding_id | [uint64](#uint64) |  |  |






<a name="core.finding.GetFindingResponse"></a>

### GetFindingResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| finding | [Finding](#core.finding.Finding) |  |  |






<a name="core.finding.GetResourceRequest"></a>

### GetResourceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [uint32](#uint32) |  |  |
| resource_id | [uint64](#uint64) |  |  |






<a name="core.finding.GetResourceResponse"></a>

### GetResourceResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| resource | [Resource](#core.finding.Resource) |  |  |






<a name="core.finding.ListFindingRequest"></a>

### ListFindingRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [uint32](#uint32) |  |  |
| project_id | [uint32](#uint32) | repeated |  |
| data_source | [string](#string) | repeated |  |
| resource_name | [string](#string) | repeated |  |
| from_score | [float](#float) |  |  |
| to_score | [float](#float) |  |  |
| from_at | [int64](#int64) |  |  |
| to_at | [int64](#int64) |  |  |






<a name="core.finding.ListFindingResponse"></a>

### ListFindingResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| finding_id | [uint64](#uint64) | repeated |  |






<a name="core.finding.ListFindingTagRequest"></a>

### ListFindingTagRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [uint32](#uint32) |  |  |
| finding_id | [uint64](#uint64) |  |  |






<a name="core.finding.ListFindingTagResponse"></a>

### ListFindingTagResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| tag | [FindingTag](#core.finding.FindingTag) | repeated |  |






<a name="core.finding.ListResourceRequest"></a>

### ListResourceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [uint32](#uint32) |  |  |
| project_id | [uint32](#uint32) | repeated |  |
| resource_name | [string](#string) | repeated |  |
| from_sum_score | [float](#float) |  |  |
| to_sum_score | [float](#float) |  |  |
| from_at | [int64](#int64) |  |  |
| to_at | [int64](#int64) |  |  |






<a name="core.finding.ListResourceResponse"></a>

### ListResourceResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| resource_id | [uint64](#uint64) | repeated |  |






<a name="core.finding.ListResourceTagRequest"></a>

### ListResourceTagRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [uint32](#uint32) |  |  |
| resource_id | [uint64](#uint64) |  |  |






<a name="core.finding.ListResourceTagResponse"></a>

### ListResourceTagResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| tag | [ResourceTag](#core.finding.ResourceTag) | repeated |  |






<a name="core.finding.PutFindingRequest"></a>

### PutFindingRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [uint32](#uint32) |  |  |
| finding | [FindingForUpsert](#core.finding.FindingForUpsert) |  |  |






<a name="core.finding.PutFindingResponse"></a>

### PutFindingResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| finding | [Finding](#core.finding.Finding) |  |  |






<a name="core.finding.PutResourceRequest"></a>

### PutResourceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [uint32](#uint32) |  |  |
| resource | [ResourceForUpsert](#core.finding.ResourceForUpsert) |  |  |






<a name="core.finding.PutResourceResponse"></a>

### PutResourceResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| resource | [Resource](#core.finding.Resource) |  |  |






<a name="core.finding.Resource"></a>

### Resource
Resource エンティティ


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| resource_id | [uint64](#uint64) |  |  |
| resource_name | [string](#string) |  |  |
| project_id | [uint32](#uint32) |  |  |
| created_at | [int64](#int64) |  |  |
| updated_at | [int64](#int64) |  |  |






<a name="core.finding.ResourceForUpsert"></a>

### ResourceForUpsert
Resource エンティティ（登録・更新用）


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| resource_name | [string](#string) |  |  |
| project_id | [uint32](#uint32) |  |  |






<a name="core.finding.ResourceTag"></a>

### ResourceTag
ResourceTag エンティティ


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| resource_tag_id | [uint64](#uint64) |  |  |
| resource_id | [uint64](#uint64) |  |  |
| tag_key | [string](#string) |  |  |
| tag_value | [string](#string) |  |  |
| created_at | [int64](#int64) |  |  |
| updated_at | [int64](#int64) |  |  |






<a name="core.finding.ResourceTagForUpsert"></a>

### ResourceTagForUpsert
ResourceTag エンティティ（登録・更新用）


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| resource_id | [uint64](#uint64) |  |  |
| tag_key | [string](#string) |  |  |
| tag_value | [string](#string) |  |  |






<a name="core.finding.TagFindingRequest"></a>

### TagFindingRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [uint32](#uint32) |  |  |
| tag | [FindingTagForUpsert](#core.finding.FindingTagForUpsert) |  |  |






<a name="core.finding.TagFindingResponse"></a>

### TagFindingResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| tag | [FindingTag](#core.finding.FindingTag) |  |  |






<a name="core.finding.TagResourceRequest"></a>

### TagResourceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [uint32](#uint32) |  |  |
| tag | [ResourceTagForUpsert](#core.finding.ResourceTagForUpsert) |  |  |






<a name="core.finding.TagResourceResponse"></a>

### TagResourceResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| tag | [ResourceTag](#core.finding.ResourceTag) |  |  |






<a name="core.finding.UntagFindingRequest"></a>

### UntagFindingRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [uint32](#uint32) |  |  |
| finding_tag_id | [uint64](#uint64) |  |  |






<a name="core.finding.UntagResourceRequest"></a>

### UntagResourceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [uint32](#uint32) |  |  |
| resource_tag_id | [uint64](#uint64) |  |  |





 

 

 


<a name="core.finding.FindingService"></a>

### FindingService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ListFinding | [ListFindingRequest](#core.finding.ListFindingRequest) | [ListFindingResponse](#core.finding.ListFindingResponse) | fiding |
| GetFinding | [GetFindingRequest](#core.finding.GetFindingRequest) | [GetFindingResponse](#core.finding.GetFindingResponse) |  |
| PutFinding | [PutFindingRequest](#core.finding.PutFindingRequest) | [PutFindingResponse](#core.finding.PutFindingResponse) |  |
| DeleteFinding | [DeleteFindingRequest](#core.finding.DeleteFindingRequest) | [Empty](#core.finding.Empty) |  |
| ListFindingTag | [ListFindingTagRequest](#core.finding.ListFindingTagRequest) | [ListFindingTagResponse](#core.finding.ListFindingTagResponse) |  |
| TagFinding | [TagFindingRequest](#core.finding.TagFindingRequest) | [TagFindingResponse](#core.finding.TagFindingResponse) |  |
| UntagFinding | [UntagFindingRequest](#core.finding.UntagFindingRequest) | [Empty](#core.finding.Empty) |  |
| ListResource | [ListResourceRequest](#core.finding.ListResourceRequest) | [ListResourceResponse](#core.finding.ListResourceResponse) | resource |
| GetResource | [GetResourceRequest](#core.finding.GetResourceRequest) | [GetResourceResponse](#core.finding.GetResourceResponse) |  |
| PutResource | [PutResourceRequest](#core.finding.PutResourceRequest) | [PutResourceResponse](#core.finding.PutResourceResponse) |  |
| DeleteResource | [DeleteResourceRequest](#core.finding.DeleteResourceRequest) | [Empty](#core.finding.Empty) |  |
| ListResourceTag | [ListResourceTagRequest](#core.finding.ListResourceTagRequest) | [ListResourceTagResponse](#core.finding.ListResourceTagResponse) |  |
| TagResource | [TagResourceRequest](#core.finding.TagResourceRequest) | [TagResourceResponse](#core.finding.TagResourceResponse) |  |
| UntagResource | [UntagResourceRequest](#core.finding.UntagResourceRequest) | [Empty](#core.finding.Empty) |  |

 



<a name="iam/iam.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## iam/iam.proto



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
| user_id | [int32](#int32) |  |  |
| action_ptn | [string](#string) |  |  |
| resource_ptn | [string](#string) |  |  |






<a name="core.iam.AuthzResponse"></a>

### AuthzResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ok | [bool](#bool) |  |  |






<a name="core.iam.User"></a>

### User
Userエンティティ


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [int32](#int32) |  |  |
| name | [string](#string) |  |  |
| activated | [bool](#bool) |  |  |
| created_at | [int64](#int64) |  |  |
| updated_at | [int64](#int64) |  |  |





 

 

 


<a name="core.iam.IAMService"></a>

### IAMService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Authn | [AuthnRequest](#core.iam.AuthnRequest) | [AuthnResponse](#core.iam.AuthnResponse) |  |
| Authz | [AuthzRequest](#core.iam.AuthzRequest) | [AuthzResponse](#core.iam.AuthzResponse) |  |

 



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


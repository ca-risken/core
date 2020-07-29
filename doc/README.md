# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [finding/entity.proto](#finding/entity.proto)
    - [Finding](#core.finding.Finding)
    - [FindingForUpsert](#core.finding.FindingForUpsert)
    - [FindingTag](#core.finding.FindingTag)
    - [FindingTagForUpsert](#core.finding.FindingTagForUpsert)
    - [Resource](#core.finding.Resource)
    - [ResourceForUpsert](#core.finding.ResourceForUpsert)
    - [ResourceTag](#core.finding.ResourceTag)
    - [ResourceTagForUpsert](#core.finding.ResourceTagForUpsert)
  
- [finding/service.proto](#finding/service.proto)
    - [DeleteFindingRequest](#core.finding.DeleteFindingRequest)
    - [DeleteResourceRequest](#core.finding.DeleteResourceRequest)
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
    - [TagFindingRequest](#core.finding.TagFindingRequest)
    - [TagFindingResponse](#core.finding.TagFindingResponse)
    - [TagResourceRequest](#core.finding.TagResourceRequest)
    - [TagResourceResponse](#core.finding.TagResourceResponse)
    - [UntagFindingRequest](#core.finding.UntagFindingRequest)
    - [UntagResourceRequest](#core.finding.UntagResourceRequest)
  
    - [FindingService](#core.finding.FindingService)
  
- [iam/entity.proto](#iam/entity.proto)
    - [Policy](#core.iam.Policy)
    - [PolicyForUpsert](#core.iam.PolicyForUpsert)
    - [Role](#core.iam.Role)
    - [RoleForUpsert](#core.iam.RoleForUpsert)
    - [RolePolicy](#core.iam.RolePolicy)
    - [User](#core.iam.User)
    - [UserForUpsert](#core.iam.UserForUpsert)
    - [UserRole](#core.iam.UserRole)
  
- [iam/policy.proto](#iam/policy.proto)
    - [AttachPolicyRequest](#core.iam.AttachPolicyRequest)
    - [AttachPolicyResponse](#core.iam.AttachPolicyResponse)
    - [DeletePolicyRequest](#core.iam.DeletePolicyRequest)
    - [DetachPolicyRequest](#core.iam.DetachPolicyRequest)
    - [GetPolicyRequest](#core.iam.GetPolicyRequest)
    - [GetPolicyResponse](#core.iam.GetPolicyResponse)
    - [ListPolicyRequest](#core.iam.ListPolicyRequest)
    - [ListPolicyResponse](#core.iam.ListPolicyResponse)
    - [PutPolicyRequest](#core.iam.PutPolicyRequest)
    - [PutPolicyResponse](#core.iam.PutPolicyResponse)
  
- [iam/role.proto](#iam/role.proto)
    - [AttachRoleRequest](#core.iam.AttachRoleRequest)
    - [AttachRoleResponse](#core.iam.AttachRoleResponse)
    - [DeleteRoleRequest](#core.iam.DeleteRoleRequest)
    - [DetachRoleRequest](#core.iam.DetachRoleRequest)
    - [GetRoleRequest](#core.iam.GetRoleRequest)
    - [GetRoleResponse](#core.iam.GetRoleResponse)
    - [ListRoleRequest](#core.iam.ListRoleRequest)
    - [ListRoleResponse](#core.iam.ListRoleResponse)
    - [PutRoleRequest](#core.iam.PutRoleRequest)
    - [PutRoleResponse](#core.iam.PutRoleResponse)
  
- [iam/service.proto](#iam/service.proto)
    - [IsAuthorizedRequest](#core.iam.IsAuthorizedRequest)
    - [IsAuthorizedResponse](#core.iam.IsAuthorizedResponse)
  
    - [IAMService](#core.iam.IAMService)
  
- [iam/user.proto](#iam/user.proto)
    - [GetUserRequest](#core.iam.GetUserRequest)
    - [GetUserResponse](#core.iam.GetUserResponse)
    - [ListUserRequest](#core.iam.ListUserRequest)
    - [ListUserResponse](#core.iam.ListUserResponse)
    - [PutUserRequest](#core.iam.PutUserRequest)
    - [PutUserResponse](#core.iam.PutUserResponse)
  
- [project/entity.proto](#project/entity.proto)
    - [Project](#core.project.Project)
  
- [project/service.proto](#project/service.proto)
    - [CreateProjectRequest](#core.project.CreateProjectRequest)
    - [CreateProjectResponse](#core.project.CreateProjectResponse)
    - [DeleteProjectRequest](#core.project.DeleteProjectRequest)
    - [ListProjectRequest](#core.project.ListProjectRequest)
    - [ListProjectResponse](#core.project.ListProjectResponse)
    - [UpdateProjectRequest](#core.project.UpdateProjectRequest)
    - [UpdateProjectResponse](#core.project.UpdateProjectResponse)
  
    - [ProjectService](#core.project.ProjectService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="finding/entity.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## finding/entity.proto



<a name="core.finding.Finding"></a>

### Finding
Finding


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| finding_id | [uint64](#uint64) |  |  |
| description | [string](#string) |  |  |
| data_source | [string](#string) |  |  |
| data_source_id | [string](#string) |  |  |
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
Finding For Upsert
(Unique keys: project_id, data_source, data_source_id)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| description | [string](#string) |  |  |
| data_source | [string](#string) |  |  |
| data_source_id | [string](#string) |  |  |
| resource_name | [string](#string) |  |  |
| project_id | [uint32](#uint32) |  |  |
| original_score | [float](#float) |  |  |
| original_max_score | [float](#float) |  |  |
| data | [string](#string) |  |  |






<a name="core.finding.FindingTag"></a>

### FindingTag
FindingTag


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| finding_tag_id | [uint64](#uint64) |  |  |
| finding_id | [uint64](#uint64) |  |  |
| project_id | [uint32](#uint32) |  |  |
| tag_key | [string](#string) |  |  |
| tag_value | [string](#string) |  |  |
| created_at | [int64](#int64) |  |  |
| updated_at | [int64](#int64) |  |  |






<a name="core.finding.FindingTagForUpsert"></a>

### FindingTagForUpsert
FindingTag For Upsert
(Unique keys: finding_id, tag_key)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| finding_id | [uint64](#uint64) |  |  |
| project_id | [uint32](#uint32) |  |  |
| tag_key | [string](#string) |  |  |
| tag_value | [string](#string) |  |  |






<a name="core.finding.Resource"></a>

### Resource
Resource


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| resource_id | [uint64](#uint64) |  |  |
| resource_name | [string](#string) |  |  |
| project_id | [uint32](#uint32) |  |  |
| created_at | [int64](#int64) |  |  |
| updated_at | [int64](#int64) |  |  |






<a name="core.finding.ResourceForUpsert"></a>

### ResourceForUpsert
Resource For upsert
(Unique keys: project_id, resource_name)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| resource_name | [string](#string) |  |  |
| project_id | [uint32](#uint32) |  |  |






<a name="core.finding.ResourceTag"></a>

### ResourceTag
ResourceTag


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| resource_tag_id | [uint64](#uint64) |  |  |
| resource_id | [uint64](#uint64) |  |  |
| project_id | [uint32](#uint32) |  |  |
| tag_key | [string](#string) |  |  |
| tag_value | [string](#string) |  |  |
| created_at | [int64](#int64) |  |  |
| updated_at | [int64](#int64) |  |  |






<a name="core.finding.ResourceTagForUpsert"></a>

### ResourceTagForUpsert
ResourceTag For upsert
(Unique keys: resource_id, tag_key)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| resource_id | [uint64](#uint64) |  |  |
| project_id | [uint32](#uint32) |  |  |
| tag_key | [string](#string) |  |  |
| tag_value | [string](#string) |  |  |





 

 

 

 



<a name="finding/service.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## finding/service.proto



<a name="core.finding.DeleteFindingRequest"></a>

### DeleteFindingRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| finding_id | [uint64](#uint64) |  |  |






<a name="core.finding.DeleteResourceRequest"></a>

### DeleteResourceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| resource_id | [uint64](#uint64) |  |  |






<a name="core.finding.GetFindingRequest"></a>

### GetFindingRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
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
| project_id | [uint32](#uint32) |  |  |
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
| project_id | [uint32](#uint32) |  |  |
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
| project_id | [uint32](#uint32) |  |  |
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
| project_id | [uint32](#uint32) |  |  |
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
| project_id | [uint32](#uint32) |  |  |
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
| project_id | [uint32](#uint32) |  |  |
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
| project_id | [uint32](#uint32) |  |  |
| resource | [ResourceForUpsert](#core.finding.ResourceForUpsert) |  |  |






<a name="core.finding.PutResourceResponse"></a>

### PutResourceResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| resource | [Resource](#core.finding.Resource) |  |  |






<a name="core.finding.TagFindingRequest"></a>

### TagFindingRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
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
| project_id | [uint32](#uint32) |  |  |
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
| project_id | [uint32](#uint32) |  |  |
| finding_tag_id | [uint64](#uint64) |  |  |






<a name="core.finding.UntagResourceRequest"></a>

### UntagResourceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| resource_tag_id | [uint64](#uint64) |  |  |





 

 

 


<a name="core.finding.FindingService"></a>

### FindingService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ListFinding | [ListFindingRequest](#core.finding.ListFindingRequest) | [ListFindingResponse](#core.finding.ListFindingResponse) | fiding |
| GetFinding | [GetFindingRequest](#core.finding.GetFindingRequest) | [GetFindingResponse](#core.finding.GetFindingResponse) |  |
| PutFinding | [PutFindingRequest](#core.finding.PutFindingRequest) | [PutFindingResponse](#core.finding.PutFindingResponse) |  |
| DeleteFinding | [DeleteFindingRequest](#core.finding.DeleteFindingRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) |  |
| ListFindingTag | [ListFindingTagRequest](#core.finding.ListFindingTagRequest) | [ListFindingTagResponse](#core.finding.ListFindingTagResponse) |  |
| TagFinding | [TagFindingRequest](#core.finding.TagFindingRequest) | [TagFindingResponse](#core.finding.TagFindingResponse) |  |
| UntagFinding | [UntagFindingRequest](#core.finding.UntagFindingRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) |  |
| ListResource | [ListResourceRequest](#core.finding.ListResourceRequest) | [ListResourceResponse](#core.finding.ListResourceResponse) | resource |
| GetResource | [GetResourceRequest](#core.finding.GetResourceRequest) | [GetResourceResponse](#core.finding.GetResourceResponse) |  |
| PutResource | [PutResourceRequest](#core.finding.PutResourceRequest) | [PutResourceResponse](#core.finding.PutResourceResponse) |  |
| DeleteResource | [DeleteResourceRequest](#core.finding.DeleteResourceRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) |  |
| ListResourceTag | [ListResourceTagRequest](#core.finding.ListResourceTagRequest) | [ListResourceTagResponse](#core.finding.ListResourceTagResponse) |  |
| TagResource | [TagResourceRequest](#core.finding.TagResourceRequest) | [TagResourceResponse](#core.finding.TagResourceResponse) |  |
| UntagResource | [UntagResourceRequest](#core.finding.UntagResourceRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) |  |

 



<a name="iam/entity.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## iam/entity.proto



<a name="core.iam.Policy"></a>

### Policy
Policy


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| policy_id | [uint32](#uint32) |  |  |
| name | [string](#string) |  |  |
| project_id | [uint32](#uint32) |  |  |
| action_ptn | [string](#string) |  |  |
| resource_ptn | [string](#string) |  |  |
| created_at | [int64](#int64) |  |  |
| updated_at | [int64](#int64) |  |  |






<a name="core.iam.PolicyForUpsert"></a>

### PolicyForUpsert
PolicyForUpsert


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| project_id | [uint32](#uint32) |  |  |
| action_ptn | [string](#string) |  |  |
| resource_ptn | [string](#string) |  |  |






<a name="core.iam.Role"></a>

### Role
Role


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| role_id | [uint32](#uint32) |  |  |
| name | [string](#string) |  |  |
| project_id | [uint32](#uint32) |  |  |
| created_at | [int64](#int64) |  |  |
| updated_at | [int64](#int64) |  |  |






<a name="core.iam.RoleForUpsert"></a>

### RoleForUpsert
RoleForUpsert


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| project_id | [uint32](#uint32) |  |  |






<a name="core.iam.RolePolicy"></a>

### RolePolicy
RolePolicy


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| role_id | [uint32](#uint32) |  |  |
| policy_id | [uint32](#uint32) |  |  |
| project_id | [uint32](#uint32) |  |  |
| created_at | [int64](#int64) |  |  |
| updated_at | [int64](#int64) |  |  |






<a name="core.iam.User"></a>

### User
User


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [uint32](#uint32) |  |  |
| sub | [string](#string) |  |  |
| name | [string](#string) |  |  |
| activated | [bool](#bool) |  |  |
| created_at | [int64](#int64) |  |  |
| updated_at | [int64](#int64) |  |  |






<a name="core.iam.UserForUpsert"></a>

### UserForUpsert
UserForUpsert
(Unique keys: sub)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| sub | [string](#string) |  |  |
| name | [string](#string) |  |  |
| activated | [bool](#bool) |  |  |






<a name="core.iam.UserRole"></a>

### UserRole
UserRole


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [uint32](#uint32) |  |  |
| role_id | [uint32](#uint32) |  |  |
| project_id | [uint32](#uint32) |  |  |
| created_at | [int64](#int64) |  |  |
| updated_at | [int64](#int64) |  |  |





 

 

 

 



<a name="iam/policy.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## iam/policy.proto



<a name="core.iam.AttachPolicyRequest"></a>

### AttachPolicyRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| role_id | [uint32](#uint32) |  |  |
| policy_id | [uint32](#uint32) |  |  |






<a name="core.iam.AttachPolicyResponse"></a>

### AttachPolicyResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| role_policy | [RolePolicy](#core.iam.RolePolicy) |  |  |






<a name="core.iam.DeletePolicyRequest"></a>

### DeletePolicyRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| policy_id | [uint32](#uint32) |  |  |






<a name="core.iam.DetachPolicyRequest"></a>

### DetachPolicyRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| role_id | [uint32](#uint32) |  |  |
| policy_id | [uint32](#uint32) |  |  |






<a name="core.iam.GetPolicyRequest"></a>

### GetPolicyRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| policy_id | [uint32](#uint32) |  |  |






<a name="core.iam.GetPolicyResponse"></a>

### GetPolicyResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| policy | [Policy](#core.iam.Policy) |  |  |






<a name="core.iam.ListPolicyRequest"></a>

### ListPolicyRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| name | [string](#string) |  |  |






<a name="core.iam.ListPolicyResponse"></a>

### ListPolicyResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| policy_id | [uint32](#uint32) | repeated |  |






<a name="core.iam.PutPolicyRequest"></a>

### PutPolicyRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| policy | [PolicyForUpsert](#core.iam.PolicyForUpsert) |  |  |






<a name="core.iam.PutPolicyResponse"></a>

### PutPolicyResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| policy | [Policy](#core.iam.Policy) |  |  |





 

 

 

 



<a name="iam/role.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## iam/role.proto



<a name="core.iam.AttachRoleRequest"></a>

### AttachRoleRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| user_id | [uint32](#uint32) |  |  |
| role_id | [uint32](#uint32) |  |  |






<a name="core.iam.AttachRoleResponse"></a>

### AttachRoleResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_role | [UserRole](#core.iam.UserRole) |  |  |






<a name="core.iam.DeleteRoleRequest"></a>

### DeleteRoleRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| role_id | [uint32](#uint32) |  |  |






<a name="core.iam.DetachRoleRequest"></a>

### DetachRoleRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| user_id | [uint32](#uint32) |  |  |
| role_id | [uint32](#uint32) |  |  |






<a name="core.iam.GetRoleRequest"></a>

### GetRoleRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| role_id | [uint32](#uint32) |  |  |






<a name="core.iam.GetRoleResponse"></a>

### GetRoleResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| role | [Role](#core.iam.Role) |  |  |






<a name="core.iam.ListRoleRequest"></a>

### ListRoleRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| name | [string](#string) |  |  |






<a name="core.iam.ListRoleResponse"></a>

### ListRoleResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| role_id | [uint32](#uint32) | repeated |  |






<a name="core.iam.PutRoleRequest"></a>

### PutRoleRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| role | [RoleForUpsert](#core.iam.RoleForUpsert) |  |  |






<a name="core.iam.PutRoleResponse"></a>

### PutRoleResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| role | [Role](#core.iam.Role) |  |  |





 

 

 

 



<a name="iam/service.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## iam/service.proto



<a name="core.iam.IsAuthorizedRequest"></a>

### IsAuthorizedRequest
IsAuthorizedRequest
ユーザからのリクエストに対して、アクションやリソースへの認可を行います


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [uint32](#uint32) |  | UserID,(e.g.)111 |
| project_id | [uint32](#uint32) |  | ProjectID,(e.g.)1001 |
| action_name | [string](#string) |  | Service&amp;API_name(&lt;service_name&gt;/&lt;API&gt;format),(e.g.)`finding/GetFinding` |
| resource_name | [string](#string) |  | System_resource_name(&lt;prefix&gt;/&lt;resouorce_name&gt;format),(e.g.)`aws:accessAnalyzer/samplebucket` |






<a name="core.iam.IsAuthorizedResponse"></a>

### IsAuthorizedResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ok | [bool](#bool) |  |  |





 

 

 


<a name="core.iam.IAMService"></a>

### IAMService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ListUser | [ListUserRequest](#core.iam.ListUserRequest) | [ListUserResponse](#core.iam.ListUserResponse) | User |
| GetUser | [GetUserRequest](#core.iam.GetUserRequest) | [GetUserResponse](#core.iam.GetUserResponse) |  |
| PutUser | [PutUserRequest](#core.iam.PutUserRequest) | [PutUserResponse](#core.iam.PutUserResponse) |  |
| ListRole | [ListRoleRequest](#core.iam.ListRoleRequest) | [ListRoleResponse](#core.iam.ListRoleResponse) | Role |
| GetRole | [GetRoleRequest](#core.iam.GetRoleRequest) | [GetRoleResponse](#core.iam.GetRoleResponse) |  |
| PutRole | [PutRoleRequest](#core.iam.PutRoleRequest) | [PutRoleResponse](#core.iam.PutRoleResponse) |  |
| DeleteRole | [DeleteRoleRequest](#core.iam.DeleteRoleRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) |  |
| AttachRole | [AttachRoleRequest](#core.iam.AttachRoleRequest) | [AttachRoleResponse](#core.iam.AttachRoleResponse) |  |
| DetachRole | [DetachRoleRequest](#core.iam.DetachRoleRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) |  |
| ListPolicy | [ListPolicyRequest](#core.iam.ListPolicyRequest) | [ListPolicyResponse](#core.iam.ListPolicyResponse) | Policy |
| GetPolicy | [GetPolicyRequest](#core.iam.GetPolicyRequest) | [GetPolicyResponse](#core.iam.GetPolicyResponse) |  |
| PutPolicy | [PutPolicyRequest](#core.iam.PutPolicyRequest) | [PutPolicyResponse](#core.iam.PutPolicyResponse) |  |
| DeletePolicy | [DeletePolicyRequest](#core.iam.DeletePolicyRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) |  |
| AttachPolicy | [AttachPolicyRequest](#core.iam.AttachPolicyRequest) | [AttachPolicyResponse](#core.iam.AttachPolicyResponse) |  |
| DetachPolicy | [DetachPolicyRequest](#core.iam.DetachPolicyRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) |  |
| IsAuthorized | [IsAuthorizedRequest](#core.iam.IsAuthorizedRequest) | [IsAuthorizedResponse](#core.iam.IsAuthorizedResponse) | 認可（ユーザがリクエストしたアクションや、リソースに対しての認可を行います） |

 



<a name="iam/user.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## iam/user.proto



<a name="core.iam.GetUserRequest"></a>

### GetUserRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [uint32](#uint32) |  |  |
| sub | [string](#string) |  |  |






<a name="core.iam.GetUserResponse"></a>

### GetUserResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user | [User](#core.iam.User) |  |  |






<a name="core.iam.ListUserRequest"></a>

### ListUserRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| name | [string](#string) |  |  |
| activated | [bool](#bool) |  |  |






<a name="core.iam.ListUserResponse"></a>

### ListUserResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [uint32](#uint32) | repeated |  |






<a name="core.iam.PutUserRequest"></a>

### PutUserRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user | [UserForUpsert](#core.iam.UserForUpsert) |  |  |






<a name="core.iam.PutUserResponse"></a>

### PutUserResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user | [User](#core.iam.User) |  |  |





 

 

 

 



<a name="project/entity.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## project/entity.proto



<a name="core.project.Project"></a>

### Project
Project


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| name | [string](#string) |  |  |
| created_at | [int64](#int64) |  |  |
| updated_at | [int64](#int64) |  |  |





 

 

 

 



<a name="project/service.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## project/service.proto



<a name="core.project.CreateProjectRequest"></a>

### CreateProjectRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [uint32](#uint32) |  | project owner |
| name | [string](#string) |  |  |






<a name="core.project.CreateProjectResponse"></a>

### CreateProjectResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project | [Project](#core.project.Project) |  |  |






<a name="core.project.DeleteProjectRequest"></a>

### DeleteProjectRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |






<a name="core.project.ListProjectRequest"></a>

### ListProjectRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [uint32](#uint32) |  |  |
| project_id | [uint32](#uint32) |  |  |
| name | [string](#string) |  |  |






<a name="core.project.ListProjectResponse"></a>

### ListProjectResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project | [Project](#core.project.Project) | repeated |  |






<a name="core.project.UpdateProjectRequest"></a>

### UpdateProjectRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| name | [string](#string) |  |  |






<a name="core.project.UpdateProjectResponse"></a>

### UpdateProjectResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project | [Project](#core.project.Project) |  |  |





 

 

 


<a name="core.project.ProjectService"></a>

### ProjectService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ListProject | [ListProjectRequest](#core.project.ListProjectRequest) | [ListProjectResponse](#core.project.ListProjectResponse) | project |
| CreateProject | [CreateProjectRequest](#core.project.CreateProjectRequest) | [CreateProjectResponse](#core.project.CreateProjectResponse) |  |
| UpdateProject | [UpdateProjectRequest](#core.project.UpdateProjectRequest) | [UpdateProjectResponse](#core.project.UpdateProjectResponse) |  |
| DeleteProject | [DeleteProjectRequest](#core.project.DeleteProjectRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) |  |

 



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


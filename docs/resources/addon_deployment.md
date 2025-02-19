---
page_title: "spectrocloud_addon_deployment Resource - terraform-provider-spectrocloud"
subcategory: ""
description: |-
  
---

# spectrocloud_addon_deployment (Resource)

  

## Example Usage



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `cluster_uid` (String)

### Optional

- `apply_setting` (String)
- `cluster_profile` (Block List) (see [below for nested schema](#nestedblock--cluster_profile))
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--cluster_profile"></a>
### Nested Schema for `cluster_profile`

Required:

- `id` (String) The ID of the cluster profile.

Optional:

- `pack` (Block List) (see [below for nested schema](#nestedblock--cluster_profile--pack))

<a id="nestedblock--cluster_profile--pack"></a>
### Nested Schema for `cluster_profile.pack`

Required:

- `name` (String) The name of the pack. The name must be unique within the cluster profile.

Optional:

- `manifest` (Block List) (see [below for nested schema](#nestedblock--cluster_profile--pack--manifest))
- `registry_uid` (String) The registry UID of the pack. The registry UID is the unique identifier of the registry.
- `tag` (String) The tag of the pack. The tag is the version of the pack.
- `type` (String) The type of the pack. The default value is `spectro`.
- `uid` (String)
- `values` (String) The values of the pack. The values are the configuration values of the pack. The values are specified in YAML format.

<a id="nestedblock--cluster_profile--pack--manifest"></a>
### Nested Schema for `cluster_profile.pack.manifest`

Required:

- `content` (String) The content of the manifest. The content is the YAML content of the manifest.
- `name` (String) The name of the manifest. The name must be unique within the pack.

Read-Only:

- `uid` (String)




<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `delete` (String)
- `update` (String)
---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "spectrocloud_role Data Source - terraform-provider-spectrocloud"
subcategory: ""
description: |-
  
---

# spectrocloud_role (Data Source)



## Example Usage

```terraform
data "spectrocloud_role" "role1" {
  name = "Project Editor"

  # (alternatively)
  # id =  "5fd0ca727c411c71b55a359c"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `name` (String)

### Read-Only

- `id` (String) The ID of this resource.



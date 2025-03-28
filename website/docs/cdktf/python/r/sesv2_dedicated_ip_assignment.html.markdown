---
subcategory: "SESv2 (Simple Email V2)"
layout: "aws"
page_title: "AWS: aws_sesv2_dedicated_ip_assignment"
description: |-
  Terraform resource for managing an AWS SESv2 (Simple Email V2) Dedicated IP Assignment.
---


<!-- Please do not edit this file, it is generated. -->
# Resource: aws_sesv2_dedicated_ip_assignment

Terraform resource for managing an AWS SESv2 (Simple Email V2) Dedicated IP Assignment.

This resource is used with "Standard" dedicated IP addresses. This includes addresses [requested and relinquished manually](https://docs.aws.amazon.com/ses/latest/dg/dedicated-ip-case.html) via an AWS support case, or [Bring Your Own IP](https://docs.aws.amazon.com/ses/latest/dg/dedicated-ip-byo.html) addresses. Once no longer assigned, this resource returns the IP to the [`ses-default-dedicated-pool`](https://docs.aws.amazon.com/ses/latest/dg/managing-ip-pools.html), managed by AWS.

## Example Usage

### Basic Usage

```python
# DO NOT EDIT. Code generated by 'cdktf convert' - Please report bugs at https://cdk.tf/bug
from constructs import Construct
from cdktf import TerraformStack
#
# Provider bindings are generated by running `cdktf get`.
# See https://cdk.tf/provider-generation for more details.
#
from imports.aws.sesv2_dedicated_ip_assignment import Sesv2DedicatedIpAssignment
class MyConvertedCode(TerraformStack):
    def __init__(self, scope, name):
        super().__init__(scope, name)
        Sesv2DedicatedIpAssignment(self, "example",
            destination_pool_name="my-pool",
            ip="0.0.0.0"
        )
```

## Argument Reference

The following arguments are required:

* `ip` - (Required) Dedicated IP address.
* `destination_pool_name` - (Required) Dedicated IP address.

## Attribute Reference

This resource exports the following attributes in addition to the arguments above:

* `id` - A comma-separated string made up of `ip` and `destination_pool_name`.

## Import

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import SESv2 (Simple Email V2) Dedicated IP Assignment using the `id`, which is a comma-separated string made up of `ip` and `destination_pool_name`. For example:

```python
# DO NOT EDIT. Code generated by 'cdktf convert' - Please report bugs at https://cdk.tf/bug
from constructs import Construct
from cdktf import TerraformStack
#
# Provider bindings are generated by running `cdktf get`.
# See https://cdk.tf/provider-generation for more details.
#
from imports.aws.sesv2_dedicated_ip_assignment import Sesv2DedicatedIpAssignment
class MyConvertedCode(TerraformStack):
    def __init__(self, scope, name):
        super().__init__(scope, name)
        Sesv2DedicatedIpAssignment.generate_config_for_import(self, "example", "0.0.0.0,my-pool")
```

Using `terraform import`, import SESv2 (Simple Email V2) Dedicated IP Assignment using the `id`, which is a comma-separated string made up of `ip` and `destination_pool_name`. For example:

```console
% terraform import aws_sesv2_dedicated_ip_assignment.example "0.0.0.0,my-pool"
```

<!-- cache-key: cdktf-0.20.8 input-3bade86996d3acf65855f6a6231ce2ef82aa6ec1d5d1f14f8816370a5fbe92d7 -->
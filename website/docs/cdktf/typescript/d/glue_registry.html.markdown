---
subcategory: "Glue"
layout: "aws"
page_title: "AWS: aws_glue_registry"
description: |-
  Terraform data source for managing an AWS Glue Registry.
---


<!-- Please do not edit this file, it is generated. -->
# Data Source: aws_glue_registry

Terraform data source for managing an AWS Glue Registry.

## Example Usage

### Basic Usage

```typescript
// DO NOT EDIT. Code generated by 'cdktf convert' - Please report bugs at https://cdk.tf/bug
import { Construct } from "constructs";
import { TerraformStack } from "cdktf";
/*
 * Provider bindings are generated by running `cdktf get`.
 * See https://cdk.tf/provider-generation for more details.
 */
import { DataAwsGlueRegistry } from "./.gen/providers/aws/data-aws-glue-registry";
class MyConvertedCode extends TerraformStack {
  constructor(scope: Construct, name: string) {
    super(scope, name);
    new DataAwsGlueRegistry(this, "example", {
      name: "example",
    });
  }
}

```

## Argument Reference

The following arguments are required:

* `name` - (Required) Name of the Glue Registry.

## Attribute Reference

This data source exports the following attributes in addition to the arguments above:

* `arn` - Amazon Resource Name (ARN) of Glue Registry.
* `description` - A description of the registry.

<!-- cache-key: cdktf-0.20.8 input-e9eb8b4e6e8262dd3c80922987d4f0253afa9a1750108733d9da020bcd0e6760 -->
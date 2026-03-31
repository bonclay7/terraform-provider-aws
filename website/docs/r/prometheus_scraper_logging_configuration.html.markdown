---
subcategory: "AMP (Managed Prometheus)"
layout: "aws"
page_title: "AWS: aws_prometheus_scraper_logging_configuration"
description: |-
  Manages an Amazon Managed Service for Prometheus (AMP) Scraper Logging Configuration.
---

# Resource: aws_prometheus_scraper_logging_configuration

Manages an Amazon Managed Service for Prometheus (AMP) Scraper Logging Configuration.

## Example Usage

```terraform
resource "aws_cloudwatch_log_group" "example" {
  name = "/aws/prometheus/scraper-logs/example"
}

resource "aws_prometheus_scraper_logging_configuration" "example" {
  scraper_id = aws_prometheus_scraper.example.id

  logging_destination {
    cloudwatch_logs {
      log_group_arn = "${aws_cloudwatch_log_group.example.arn}:*"
    }
  }
}
```

### With Scraper Components

```terraform
resource "aws_cloudwatch_log_group" "example" {
  name = "/aws/prometheus/scraper-logs/example"
}

resource "aws_prometheus_scraper_logging_configuration" "example" {
  scraper_id = aws_prometheus_scraper.example.id

  logging_destination {
    cloudwatch_logs {
      log_group_arn = "${aws_cloudwatch_log_group.example.arn}:*"
    }
  }

  scraper_components {
    type = "COLLECTOR"
  }

  scraper_components {
    type = "EXPORTER"
  }

  scraper_components {
    type = "SERVICE_DISCOVERY"
  }
}
```

## Argument Reference

This resource supports the following arguments:

* `logging_destination` - (Required) Configuration block for the logging destination. See [`logging_destination`](#logging_destination).
* `scraper_id` - (Required, Forces new resource) The ID of the AMP scraper for which to configure logging.

The following arguments are optional:

* `region` - (Optional) Region where this resource will be [managed](https://docs.aws.amazon.com/general/latest/gr/rande.html#regional-endpoints). Defaults to the Region set in the [provider configuration](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#aws-configuration-reference).
* `scraper_components` - (Optional) Set of configuration blocks for the scraper components to enable logging for. See [`scraper_components`](#scraper_components). If not specified, all components are logged.

### `logging_destination`

* `cloudwatch_logs` - (Required) Configuration block for CloudWatch Logs destination. See [`cloudwatch_logs`](#cloudwatch_logs).

#### `cloudwatch_logs`

* `log_group_arn` - (Required) The ARN of the CloudWatch log group to which scraper logs will be sent. The ARN must end with `:*`.

### `scraper_components`

* `type` - (Required) The type of scraper component to configure for logging. Valid values: `COLLECTOR`, `EXPORTER`, `SERVICE_DISCOVERY`.

The following arguments are optional:

* `config` - (Optional) Configuration block for component-specific settings. See [`config`](#config).

#### `config`

* `options` - (Optional) A map of configuration options for the scraper component.

## Attribute Reference

This resource exports no additional attributes.

## Import

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import the Scraper Logging Configuration using the scraper ID. For example:

```terraform
import {
  to = aws_prometheus_scraper_logging_configuration.example
  id = "s-example1-1234-abcd-5678-ef9012abcd34"
}
```

## Timeouts

[Configuration options](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts):

- `create` - (Default `5m`)
- `update` - (Default `5m`)
- `delete` - (Default `5m`)

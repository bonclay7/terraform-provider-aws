---
subcategory: "OpsWorks"
layout: "aws"
page_title: "AWS: aws_opsworks_custom_layer"
description: |-
  Provides an OpsWorks custom layer resource.
---

# Resource: aws_opsworks_custom_layer

Provides an OpsWorks custom layer resource.

!> **ALERT:** AWS no longer supports OpsWorks Stacks. All related resources will be removed from the Terraform AWS Provider in the next major version.

## Example Usage

```terraform
resource "aws_opsworks_custom_layer" "custlayer" {
  name       = "My Awesome Custom Layer"
  short_name = "awesome"
  stack_id   = aws_opsworks_stack.main.id
}
```

## Argument Reference

This resource supports the following arguments:

* `name` - (Required) A human-readable name for the layer.
* `short_name` - (Required) A short, machine-readable name for the layer, which will be used to identify it in the Chef node JSON.
* `stack_id` - (Required) ID of the stack the layer will belong to.
* `auto_assign_elastic_ips` - (Optional) Whether to automatically assign an elastic IP address to the layer's instances.
* `auto_assign_public_ips` - (Optional) For stacks belonging to a VPC, whether to automatically assign a public IP address to each of the layer's instances.
* `cloudwatch_configuration` - (Optional) Will create an EBS volume and connect it to the layer's instances. See [Cloudwatch Configuration](#cloudwatch-configuration).
* `custom_instance_profile_arn` - (Optional) The ARN of an IAM profile that will be used for the layer's instances.
* `custom_security_group_ids` - (Optional) Ids for a set of security groups to apply to the layer's instances.
* `auto_healing` - (Optional) Whether to enable auto-healing for the layer.
* `install_updates_on_boot` - (Optional) Whether to install OS and package updates on each instance when it boots.
* `instance_shutdown_timeout` - (Optional) The time, in seconds, that OpsWorks will wait for Chef to complete after triggering the Shutdown event.
* `elastic_load_balancer` - (Optional) Name of an Elastic Load Balancer to attach to this layer
* `drain_elb_on_shutdown` - (Optional) Whether to enable Elastic Load Balancing connection draining.
* `load_based_auto_scaling` - (Optional) Load-based auto scaling configuration. See [Load Based AutoScaling](#load-based-autoscaling)
* `system_packages` - (Optional) Names of a set of system packages to install on the layer's instances.
* `use_ebs_optimized_instances` - (Optional) Whether to use EBS-optimized instances.
* `ebs_volume` - (Optional) Will create an EBS volume and connect it to the layer's instances. See [EBS Volume](#ebs-volume).
* `custom_json` - (Optional) Custom JSON attributes to apply to the layer.
* `tags` - (Optional) A map of tags to assign to the resource. If configured with a provider [`default_tags` configuration block](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block) present, tags with matching keys will overwrite those defined at the provider-level.

The following extra optional arguments, all lists of Chef recipe names, allow
custom Chef recipes to be applied to layer instances at the five different
lifecycle events, if custom cookbooks are enabled on the layer's stack:

* `custom_configure_recipes`
* `custom_deploy_recipes`
* `custom_setup_recipes`
* `custom_shutdown_recipes`
* `custom_undeploy_recipes`

### EBS Volume

* `mount_point` - (Required) The path to mount the EBS volume on the layer's instances.
* `size` - (Required) The size of the volume in gigabytes.
* `number_of_disks` - (Required) The number of disks to use for the EBS volume.
* `raid_level` - (Required) The RAID level to use for the volume.
* `type` - (Optional) The type of volume to create. This may be `standard` (the default), `io1` or `gp2`.
* `iops` - (Optional) For PIOPS volumes, the IOPS per disk.
* `encrypted` - (Optional) Encrypt the volume.

### Cloudwatch Configuration

* `enabled` - (Optional)
* `log_streams` - (Optional) A block the specifies how an opsworks logs look like. See [Log Streams](#log-streams).

#### Log Streams

* `file` - (Required) Specifies log files that you want to push to CloudWatch Logs. File can point to a specific file or multiple files (by using wild card characters such as /var/log/system.log*).
* `log_group_name` - (Required) Specifies the destination log group. A log group is created automatically if it doesn't already exist.
* `batch_count` - (Optional) Specifies the max number of log events in a batch, up to `10000`. The default value is `1000`.
* `batch_size` - (Optional) Specifies the maximum size of log events in a batch, in bytes, up to `1048576` bytes. The default value is `32768` bytes.
* `buffer_duration` - (Optional) Specifies the time duration for the batching of log events. The minimum value is `5000` and default value is `5000`.
* `datetime_format` - (Optional) Specifies how the timestamp is extracted from logs. For more information, see the CloudWatch Logs Agent Reference (https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/AgentReference.html).
* `encoding` - (Optional) Specifies the encoding of the log file so that the file can be read correctly. The default is `utf_8`.
* `file_fingerprint_lines` - (Optional) Specifies the range of lines for identifying a file. The valid values are one number, or two dash-delimited numbers, such as `1`, `2-5`. The default value is `1`.
* `initial_position` - (Optional) Specifies where to start to read data (`start_of_file` or `end_of_file`). The default is `start_of_file`.
* `multiline_start_pattern` - (Optional) Specifies the pattern for identifying the start of a log message.
* `time_zone` - (Optional) Specifies the time zone of log event time stamps.

### Load Based Autoscaling

* `downscaling` - (Optional) The downscaling settings, as defined below, used for load-based autoscaling
* `enable` - (Optional) Whether load-based auto scaling is enabled for the layer.
* `upscaling` - (Optional) The upscaling settings, as defined below, used for load-based autoscaling

The `downscaling` and `upscaling` blocks supports the following arguments:

Though the three thresholds are optional, at least one threshold must be set when using load-based autoscaling.

* `alarms` - (Optional) Custom Cloudwatch auto scaling alarms, to be used as thresholds. This parameter takes a list of up to five alarm names, which are case sensitive and must be in the same region as the stack.
* `cpu_threshold` - (Optional) The CPU utilization threshold, as a percent of the available CPU. A value of -1 disables the threshold.
* `ignore_metrics_time` - (Optional) The amount of time (in minutes) after a scaling event occurs that AWS OpsWorks Stacks should ignore metrics and suppress additional scaling events.
* `instance_count` - (Optional) The number of instances to add or remove when the load exceeds a threshold.
* `load_threshold` - (Optional) The load threshold. A value of -1 disables the threshold.
* `memory_threshold` - (Optional) The memory utilization threshold, as a percent of the available memory. A value of -1 disables the threshold.
* `thresholds_wait_time` - (Optional) The amount of time, in minutes, that the load must exceed a threshold before more instances are added or removed.

## Attribute Reference

This resource exports the following attributes in addition to the arguments above:

* `id` - The id of the layer.
* `arn` - The Amazon Resource Name(ARN) of the layer.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider [`default_tags` configuration block](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block).

## Import

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import OpsWorks Custom Layers using the `id`. For example:

```terraform
import {
  to = aws_opsworks_custom_layer.bar
  id = "00000000-0000-0000-0000-000000000000"
}
```

Using `terraform import`, import OpsWorks Custom Layers using the `id`. For example:

```console
% terraform import aws_opsworks_custom_layer.bar 00000000-0000-0000-0000-000000000000
```

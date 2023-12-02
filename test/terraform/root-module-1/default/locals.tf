locals {
  aws_region_to_deploy = var.aws_region

  /*
    * Feature flags
  */
  is_enabled = !var.is_enabled ? false : var.module_config == null ? false : length(var.module_config) > 0
  /*
    * SSM parameter store normalisation process.
  */
  parameters_cfg_normalized = !local.is_enabled ? [] : [
    for cfg in var.module_config : {
      id = trimspace(lower(cfg.name))
    }
  ]

  parameters_cfg_create = !local.is_enabled ? {} : {
    for cfg in local.parameters_cfg_normalized : cfg["id"] => cfg
  }
}

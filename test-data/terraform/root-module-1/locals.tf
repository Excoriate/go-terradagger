locals {
  /*
    * Feature flags
  */
  is_enabled = !var.is_enabled ? false : var.module_config == null ? false : length(var.module_config) > 0
  /*
    * SSM parameter store normalisation process.
  */
}

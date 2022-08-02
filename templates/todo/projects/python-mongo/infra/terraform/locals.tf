locals {
  local_tags = {
    azd-env-name = var.name
  }

  tags = merge(var.tags, local.local_tags)
}
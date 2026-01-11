# Terraform Provider for Firefly III

This is a [Terraform](https://www.terraform.io) provider for [Firefly III](https://www.firefly-iii.org), a free and open source personal finance manager.

## Documentation

For complete documentation on how to use this provider, including all available resources and data sources, please visit the [Terraform Registry documentation](https://registry.terraform.io/providers/renescheepers/firefly3/latest/docs).

## Using the Provider

For usage examples and complete documentation, please see the [provider documentation on the Terraform Registry](https://registry.terraform.io/providers/renescheepers/firefly3/latest/docs).

### Quick Example

```hcl
terraform {
  required_providers {
    firefly3 = {
      source  = "renescheepers/firefly3"
      version = "~> 1.0"
    }
  }
}

provider "firefly3" {
  url   = "https://your-firefly-instance.com"
  token = var.firefly_token
}

resource "firefly3_category" "groceries" {
  name  = "Groceries"
  notes = "Food and household items"
}
```

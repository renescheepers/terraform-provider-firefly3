provider "firefly3" {
  endpoint = "https://firefly3.example.com"
  api_key  = "TOKEN"
}

resource "firefly3_category" "category" {
  name  = "Category 1"
  notes = ""
}

resource "firefly3_rule_group" "group" {
  title       = "Group 1"
  description = "All your rules in this group."
}

resource "firefly3_rule" "rule" {
  rule_group_id = firefly3_rule_group.group.id
  title         = "Example"
  description   = "Description"
  trigger       = "store-journal"

  triggers = [{
    type  = "description_contains"
    value = "description"
    }
  ]
  actions = [
    {
      type  = "set_category"
      value = firefly3_category.category.name
    }
  ]
}

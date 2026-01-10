provider "firefly3" {
  endpoint = "ENDPOINT"
}

terraform {
    required_providers {
      firefly3 = {
        source = "renescheepers/firefly3"
      }
    }
  }

resource "firefly3_rule" "test" {
  rule_group_id = 1
    title = "testing_rule"
    description = "test"
    
    trigger = "update-journal"
    
    triggers = [{
      type = "description_contains"
      value = "some value"
    }
    ]
    actions = [
      {
        type = "set_category"
        value = "Boodschappen"
      }
    ]
}
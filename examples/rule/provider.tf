provider "firefly3" {
  endpoint = "https://ff3.renescheepers.nl"
  api_key = "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJhdWQiOiIxIiwianRpIjoiNWNmNzNmYWE4ZGNlZTg4NWY1NTJmNDQ0MDdhODZkMTc4ZWNhZDlhY2U0NmQxZDQ1NThhMjI0NzdmYzdjMjJkM2RiNGFiYjk0YzVjMmU4ODEiLCJpYXQiOjE3NjgwMzE2MTcuMTQwNzkzLCJuYmYiOjE3NjgwMzE2MTcuMTQwNzk3LCJleHAiOjE3OTk1Njc2MTYuODU2NTQxLCJzdWIiOiIxIiwic2NvcGVzIjpbXX0.IUOc1OnWZQBIUxgl3iaXHCXGbh0iRrsfK4VU8wVfB3jYXyTTQLgP0W9P-XC69t0dte4vBSXGdASDw5WwZ7jTLiV8KAmy1ntIeeexLH8vu4sUXLUm4IgbEz-lkO2uZCTcnDiwYjW854vItwLr0OJi6p8QZ_RPy3FG1pXKyHpDzWZOD3IETz84C1pIapfKriykfGrwC5_ybzvThOdfZg1YI5gwdS0q7JuF0AgKGFxtWTp_Hv160ak9Cxpetqs4L2e8_iYlh1s3eX5ebzhO4mybN5746BsHPLn70cP825Y8IA9LLhIEf1FeXn1tOCZ8SFETkqbUf1XKiv3ii8AkscuU3OmlAAkIqyshkZ39B6blS3cFlQjnF2LWbma8wa2-fpx3ccu9-OQAo3EsxCoxyJBxmi8zL2m-I9PERYSeRawacchWvegne2js7pB6N2aQ4UmRS2FfGzqXY9OGyj4_YklckmBFkdqjUZ37YlckxFvmFWbffT5zmKbt01fRt7oGrYkq1pVnJrBWkLQquBXPKgpTztFL6-0liH5c2iTPcZiLwlziM-iXAQ2zvVwP3cXvYJFvnUCDypM8n-HR5E9x4ald5AN69gIWY4eQKhVfp2-59izy8dBg594CvZDUb1VMLnREPOWjHzcJ5RcA9l-8dhkletRQIjx58xeZ2TKvppcickM"
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
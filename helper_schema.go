package main

import "github.com/hashicorp/terraform/helper/schema"

func resourceSet(v interface{}) int {
	resource := v.(map[string]interface{})
	return schema.HashString(resource["href"].(string))
}

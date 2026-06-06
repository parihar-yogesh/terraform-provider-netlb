package main

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"address": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("NETLB_ADDRESS", nil),
				Description: "Address of the netlb API server (e.g. http://localhost:8080)",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"netlb_pool":           resourcePool(),
			"netlb_virtual_server": resourceVirtualServer(),
			"netlb_monitor":        resourceMonitor(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	address := d.Get("address").(string)
	return NewClient(address), nil
}

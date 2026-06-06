package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// VirtualServer represents the entry point for incoming traffic, routing it to a pool.
type VirtualServer struct {
	Name        string `json:"name"`
	Destination string `json:"destination"`
	Port        int    `json:"port"`
	Pool        string `json:"pool"`
	Monitor     string `json:"monitor"`
}

func resourceVirtualServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVirtualServerCreate,
		ReadContext:   resourceVirtualServerRead,
		UpdateContext: resourceVirtualServerUpdate,
		DeleteContext: resourceVirtualServerDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true, // name is immutable — changing it requires destroying and recreating the resource
			},
			"destination": {
				Type:     schema.TypeString,
				Required: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"pool": {
				Type:     schema.TypeString,
				Required: true,
			},
			"monitor": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "http",
			},
		},
	}
}

func resourceVirtualServerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	vs := VirtualServer{
		Name:        d.Get("name").(string),
		Destination: d.Get("destination").(string),
		Port:        d.Get("port").(int),
		Pool:        d.Get("pool").(string),
		Monitor:     d.Get("monitor").(string),
	}

	_, err := client.doRequest("POST", "/api/virtualservers/"+vs.Name, vs)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating virtual server: %w", err))
	}

	d.SetId(vs.Name)
	return resourceVirtualServerRead(ctx, d, m)
}

func resourceVirtualServerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	resp, err := client.doRequest("GET", "/api/virtualservers/"+d.Id(), nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading virtual server: %w", err))
	}

	// nil response means the resource no longer exists — remove from state
	if resp == nil {
		d.SetId("")
		return nil
	}

	var vs VirtualServer
	if err := json.Unmarshal(resp, &vs); err != nil {
		return diag.FromErr(fmt.Errorf("error parsing virtual server response: %w", err))
	}

	d.Set("name", vs.Name)
	d.Set("destination", vs.Destination)
	d.Set("port", vs.Port)
	d.Set("pool", vs.Pool)
	d.Set("monitor", vs.Monitor)

	return nil
}

func resourceVirtualServerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	vs := VirtualServer{
		Name:        d.Get("name").(string),
		Destination: d.Get("destination").(string),
		Port:        d.Get("port").(int),
		Pool:        d.Get("pool").(string),
		Monitor:     d.Get("monitor").(string),
	}

	_, err := client.doRequest("PUT", "/api/virtualservers/"+d.Id(), vs)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating virtual server: %w", err))
	}

	return resourceVirtualServerRead(ctx, d, m)
}

func resourceVirtualServerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	_, err := client.doRequest("DELETE", "/api/virtualservers/"+d.Id(), nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting virtual server: %w", err))
	}

	d.SetId("")
	return nil
}

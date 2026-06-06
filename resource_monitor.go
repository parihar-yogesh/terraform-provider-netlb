package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Monitor represents a health check that verifies pool members are alive.
type Monitor struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Interval int    `json:"interval"`
	Timeout  int    `json:"timeout"`
}

func resourceMonitor() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMonitorCreate,
		ReadContext:   resourceMonitorRead,
		UpdateContext: resourceMonitorUpdate,
		DeleteContext: resourceMonitorDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true, // name is immutable — changing it requires destroying and recreating the resource
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"interval": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  5,
			},
			"timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  16,
			},
		},
	}
}

func resourceMonitorCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	monitor := Monitor{
		Name:     d.Get("name").(string),
		Type:     d.Get("type").(string),
		Interval: d.Get("interval").(int),
		Timeout:  d.Get("timeout").(int),
	}

	_, err := client.doRequest("POST", "/api/monitors", monitor)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating monitor: %w", err))
	}

	d.SetId(monitor.Name)
	return resourceMonitorRead(ctx, d, m)
}

func resourceMonitorRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	resp, err := client.doRequest("GET", "/api/monitors/"+d.Id(), nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading monitor: %w", err))
	}

	// nil response means the resource no longer exists — remove from state
	if resp == nil {
		d.SetId("")
		return nil
	}

	var monitor Monitor
	if err := json.Unmarshal(resp, &monitor); err != nil {
		return diag.FromErr(fmt.Errorf("error parsing monitor response: %w", err))
	}

	d.Set("name", monitor.Name)
	d.Set("type", monitor.Type)
	d.Set("interval", monitor.Interval)
	d.Set("timeout", monitor.Timeout)

	return nil
}

func resourceMonitorUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	monitor := Monitor{
		Name:     d.Get("name").(string),
		Type:     d.Get("type").(string),
		Interval: d.Get("interval").(int),
		Timeout:  d.Get("timeout").(int),
	}

	_, err := client.doRequest("PUT", "/api/monitors/"+d.Id(), monitor)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating monitor: %w", err))
	}

	return resourceMonitorRead(ctx, d, m)
}

func resourceMonitorDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	_, err := client.doRequest("DELETE", "/api/monitors/"+d.Id(), nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting monitor: %w", err))
	}

	d.SetId("")
	return nil
}

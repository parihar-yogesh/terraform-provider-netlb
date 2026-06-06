package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Pool represents a load balancer pool and its members.
type Pool struct {
	Name     string   `json:"name"`
	Monitor  string   `json:"monitor"`
	LBMethod string   `json:"lb_method"`
	Members  []string `json:"members"`
}

func resourcePool() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePoolCreate,
		ReadContext:   resourcePoolRead,
		UpdateContext: resourcePoolUpdate,
		DeleteContext: resourcePoolDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true, // name is immutable — changing it requires destroying and recreating the resource
			},
			"monitor": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "http",
			},
			"lb_method": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "round-robin",
			},
			"members": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourcePoolCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	members := expandStringList(d.Get("members").([]interface{}))

	pool := Pool{
		Name:     d.Get("name").(string),
		Monitor:  d.Get("monitor").(string),
		LBMethod: d.Get("lb_method").(string),
		Members:  members,
	}

	_, err := client.doRequest("POST", "/api/pools/"+pool.Name, pool)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating pool: %w", err))
	}

	d.SetId(pool.Name)
	return resourcePoolRead(ctx, d, m)
}

func resourcePoolRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	resp, err := client.doRequest("GET", "/api/pools/"+d.Id(), nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading pool: %w", err))
	}

	// nil response means the resource no longer exists — remove from state
	if resp == nil {
		d.SetId("")
		return nil
	}

	var pool Pool
	if err := json.Unmarshal(resp, &pool); err != nil {
		return diag.FromErr(fmt.Errorf("error parsing pool response: %w", err))
	}

	d.Set("name", pool.Name)
	d.Set("monitor", pool.Monitor)
	d.Set("lb_method", pool.LBMethod)
	d.Set("members", pool.Members)

	return nil
}

func resourcePoolUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	members := expandStringList(d.Get("members").([]interface{}))

	pool := Pool{
		Name:     d.Get("name").(string),
		Monitor:  d.Get("monitor").(string),
		LBMethod: d.Get("lb_method").(string),
		Members:  members,
	}

	_, err := client.doRequest("PUT", "/api/pools/"+d.Id(), pool)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating pool: %w", err))
	}

	return resourcePoolRead(ctx, d, m)
}

func resourcePoolDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	_, err := client.doRequest("DELETE", "/api/pools/"+d.Id(), nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting pool: %w", err))
	}

	d.SetId("")
	return nil
}

// expandStringList converts Terraform's []interface{} to []string.
func expandStringList(input []interface{}) []string {
	result := make([]string, len(input))
	for i, v := range input {
		result[i] = v.(string)
	}
	return result
}

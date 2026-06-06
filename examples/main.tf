terraform {
  required_providers {
    netlb = {
      source = "parihar-yogesh/netlb"
    }
  }
}

provider "netlb" {
  address = "http://localhost:8080"
}

resource "netlb_monitor" "http_check" {
  name     = "http-monitor"
  type     = "http"
  interval = 5
  timeout  = 16
}

resource "netlb_pool" "web_pool" {
  name      = "web-pool"
  monitor   = netlb_monitor.http_check.name
  lb_method = "round-robin"
  members   = ["10.0.0.1:80", "10.0.0.2:80"]
}

resource "netlb_virtual_server" "web_vs" {
  name        = "web-vs"
  destination = "192.168.1.100"
  port        = 80
  pool        = netlb_pool.web_pool.name
  monitor     = netlb_monitor.http_check.name
}
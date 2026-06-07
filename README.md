# terraform-provider-netlb

A Terraform provider built in Go that manages network load balancer resources: pools, monitors, and virtual servers.

## Resources

- `netlb_monitor` — health check that verifies backend servers are alive
- `netlb_pool` — group of backend servers with a configurable load balancing method
- `netlb_virtual_server` — entry point for incoming traffic, routes connections to a pool

## How it works

The provider communicates with a local mock REST API server over HTTP/JSON. The mock server stores resources in memory and handles full CRUD for all three resource types.

## Running locally

Start the mock server in one terminal:

```bash
go run ./mockserver/server.go
```

Build the provider in another terminal:

```bash
go build -o terraform-provider-netlb .
```

Add this to `~/.terraformrc`:

```hcl
provider_installation {
  dev_overrides {
    "parihar-yogesh/netlb" = "/Users/your-username/terraform-provider-netlb"
  }
  direct {}
}
```

Apply the example:

```bash
cd examples
terraform apply
```

## Tests

```bash
go test ./... -v
```
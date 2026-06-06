# terraform-provider-netlb

A Terraform provider built in Go that manages load balancer resources — pools, monitors, and virtual servers.

Written to understand how Terraform providers work internally, using the same SDK patterns as [F5Networks/terraform-provider-bigip](https://github.com/F5Networks/terraform-provider-bigip).

## Resources

- `netlb_monitor` — health check for pool members
- `netlb_pool` — group of backend servers with a load balancing method
- `netlb_virtual_server` — frontend that routes incoming traffic to a pool

## How it works

The provider talks to a local mock server over HTTP/JSON, the same way the F5 provider talks to BIG-IP's iControl REST API. The mock server stores state in memory and handles full CRUD for all three resources.

## Running locally

Start the mock server in one terminal:

```bash
go run ./mockserver/server.go
```

Build the provider in another terminal:

```bash
go build -o terraform-provider-netlb .
```

Tell Terraform to use the local binary. Add this to `~/.terraformrc`:

```hcl
provider_installation {
  dev_overrides {
    "parihar-yogesh/netlb" = "/path/to/terraform-provider-netlb"
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
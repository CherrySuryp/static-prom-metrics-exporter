# Static Prometheus Metrics Exporter
Define and export static prometheus metrics via simple YAML configuration

## Flags
```yaml
--help
--config <Path To Config File> # default ./config.yml
--port 1234 # default 9002
--tls-crt /path/to/tls/crt # optional
--tls-key /path/to/tls/key # optional
```
## Configuration file
```yaml
server:
  basic-auth:
    user: "bcrypt-hashed-password" # example: $2a$10$fRXmD.HuavUaUCq4Lp8UK.YmcgzfIxrfH1uZ2l3whKMcy7uthThli
static_metrics:
  - name: bandwidth_limit_bytes
    help: "Network bandwidth limit in bytes" # optional
    value: 123123123123123
  - name: custom_metric
    value: 89
  - name: useless_metric
    help: "This metric's absolutely useless"
    value: 1
```
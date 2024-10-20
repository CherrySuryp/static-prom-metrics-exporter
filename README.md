# Static Prometheus Metrics Exporter
Define and export static prometheus metrics via simple YAML configuration

## Flags
```yaml
--help
--config /path/to/config.yml # string default ./config.yml
--port 1234 # int default 9002
--tls-crt /path/to/tls/crt # string optional
--tls-key /path/to/tls/key # string optional
```
## Configuration file
```yaml
server:
  basic-auth:
    user: "bcrypt-hashed-password" # example: $2a$10$fRXmD.HuavUaUCq4Lp8UK.YmcgzfIxrfH1uZ2l3whKMcy7uthThli
static_metrics:
  - name: bandwidth_limit_bytes
    help: "Network bandwidth limit in bytes" # optional
    value: 123123123123123 # integer or float
  - name: custom_metric
    value: 89
  - name: useless_metric
    help: "This metric's absolutely useless"
    value: 1
```
## Flags
```yaml
--config <Path To Config File>
--help
```
## Configuration file
```yaml
server:
  port: 9090
  tls_crt: "path/to/tls/crt" # optional
  tls_key: "path/to/tls/key" # optional
  basic-auth:
    user: password
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
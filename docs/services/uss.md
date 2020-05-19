# uss

[UPYUN Storage Service](https://www.upyun.com/products/file-storage)

`uss://hmac:<access_key>:<secret_key>/?name=<bucket_name>&work_dir=<prefix>`

## Config

### Storager

| Name | Required | Comments |
| ---- | -------- | -------- |
| `credential` | Y | only support `hmac` protocol |
| `name` | Y | bucket name |
| `work_dir` | N | work dir |
| `endpoint` | Y | specific domain to access this storager |

## Example

Init storager

```yaml
credential: hmac:<operator_name>:<operator_key>
name: <bucket_name>
work_dir: /<work_dir>
endpoint: https:<domain>
```

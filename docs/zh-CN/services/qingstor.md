# qingstor

[QingStor Object Storage](https://www.qingcloud.com/products/qingstor/)

## Config

### Servicer

| Name         | Required | Comments                     |
| ------------ | -------- | ---------------------------- |
| `credential` | Y        | only support `hmac` protocol |
| `endpoint`   | Y        |                              |

### Storager

| Name       | Required | Comments    |
| ---------- | -------- | ----------- |
| `name`     | Y        | bucket name |
| `work_dir` | N        | work dir    |

## Example

Init servicer

```yaml
credential: hmac:<account_name>:<account_key>
endpoint: https:<account_name>.<endpoint_suffix>
```

Init storager

```yaml
credential: hmac:<account_name>:<account_key>
endpoint: https:<account_name>.<endpoint_suffix>
name: <container_name>
work_dir: /<work_dir>
```

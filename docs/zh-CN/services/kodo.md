# kodo

[qiniu kodo](https://www.qiniu.com/products/kodo)

## Config

### Servicer

| Name         | Required | Comments                     |
| ------------ | -------- | ---------------------------- |
| `credential` | Y        | only support `hmac` protocol |

### Storager

| Name       | Required | Comments                                |
| ---------- | -------- | --------------------------------------- |
| `name`     | Y        | bucket name                             |
| `work_dir` | N        | work dir                                |
| `endpoint` | Y        | specific domain to access this storager |

## Example

Init servicer

```yaml
credential: hmac:<access_key>:<secret_key>
```

Init storager

```yaml
credential: hmac:<access_key>:<secret_key>
name: <bucket_name>
work_dir: /<work_dir>
endpoint: http:<domain>
```


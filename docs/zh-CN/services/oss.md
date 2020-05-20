# oss

[Aliyun Object Storage](https://www.aliyun.com/product/oss)

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
credential: hmac:<access_key>:<secret_key>
endpoint: https:<location>.aliyuncs.com
```

Init storager

```yaml
credential: hmac:<access_key>:<secret_key>
endpoint: https:<location>.aliyuncs.com
name: <bucket_name>
work_dir: /<work_dir>
location: <bucket_location>
```

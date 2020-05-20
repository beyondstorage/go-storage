# cos

[Tencent Cloud Object Storage](https://cloud.tencent.com/product/cos)

## Config

### Servicer

| Name         | Required | Comments                     |
| ------------ | -------- | ---------------------------- |
| `credential` | Y        | only support `hmac` protocol |


### Storager

| Name       | Required | Comments    |
| ---------- | -------- | ----------- |
| `name`     | Y        | bucket name |
| `work_dir` | N        | work dir    |
| `location` | Y        | location    |

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
location: <bucket_location>
```

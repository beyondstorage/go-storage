# uss

[UPYUN 存储服务](https://www.upyun.com/products/file-storage)

`uss://hmac:<access_key>:<secret_key>/?name=<bucket_name>&work_dir=<prefix>`

## 配置

### 存储器

| 名称     | 必填 | 评论            |
| ------ | -- | ------------- |
| `凭据`   | Y  | 仅支持 `hmac` 协议 |
| `名称`   | Y  | 桶名称           |
| `工作目录` | N  | 工作目录          |

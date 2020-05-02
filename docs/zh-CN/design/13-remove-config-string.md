---
author: Xuanwo <github@xuanwo.io>
status: 完成
updated_at: 2020-03-06
deprecates:
  - design/3-support-service-init-via-config-string.md
---

# 建议：删除配置字符串

## 二. 背景

[存储](https://github.com/Xuanwo/storage) 在提议 [3-support-service-init-via-config-string](https://github.com/Xuanwo/storage/blob/master/docs/design/3-support-service-init-via-config-string.md)中添加了配置字符串支持，并在提议 [9-remove-storager-init](https://github.com/Xuanwo/storage/blob/master/docs/design/9-remove-storager-init.md) 中更新。 这些建议允许用户进入如下服务：

```go
srv, store, err := coreutils.Open("fs:///?work_dir=/path/to/dir")
if err != nil vol
    log。Fatalf("service init failed: %v", err)

```

随着时间的推移，对 [存储](https://github.com/Xuanwo/storage)的配置有更深刻的理解，我找到 `配置字符串` 不是一个好的解决方案。 `配置字符串` 确实有一些好处：简单字符串，易于构造，易于理解(?)。

然而, 在演示项目的一些经验 [bard](https://github.com/Xuanwo/bard), 这是基于 [存储设备的粘贴bin 服务](https://github.com/Xuanwo/storage), 我找到 `配置字符串` 深深影响最终用户侧配置。 [bard](https://github.com/Xuanwo/bard)的配置看起来像以下：

```yaml
public_url: http://127.0.0.1:8080
listen: 127.0.0.1:8080

key: xxxxxxx
max_file_size: 104857600

database:
  type: sqlite3
  connection: "/tmp/bard/db"

storage: "fs:///?work_dir=/tmp/bard/data"
```

每个应用程序都建立在 [存储](https://github.com/Xuanwo/storage) 上，要么直接暴露配置字符串以结束用户，要么写一个格式配置函数以将他们自己的配置转换为 [存储](https://github.com/Xuanwo/storage) 配置字符串。 这是预料不到的。

不仅如此，配置字符串还难以构建安全配对。 我们需要从字符串解析它们，无法有效地使用现有的配置格式。

## 建议

因此，我提议作如下修改：

- 删除配置字符串的概念
- 将 `Open(fg 字符串)` 重设为 `Open(t 字符串，选择 []*类型。配对)`

添加配置类型 `配置` 以帮助开发者解析对：

```go
类型 config structt v.
    Type 字符串
    选项[string]string
}

func (c*Config) Parse(t 字符串, []*类型.Pair, 错误) {}
```

## 理由

无

## 兼容性

以下软件包将受到影响：

- `coreutils`
- `pkg/config`

## 二． 执行情况

大多数工作将由本提案的作者完成。
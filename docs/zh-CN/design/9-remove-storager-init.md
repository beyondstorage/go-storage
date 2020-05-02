---
author: Xuanwo <github@xuanwo.io>
status: 完成
updated_at: 2020-01-16
updates:
  - design/3-support-service-init-via-config-string.md
---

# 建议：删除存储器Init

## 二. 背景

提交 [存储器：所有 API 现在使用相对路径而不是](https://github.com/Xuanwo/storage/commit/1cb485ec1f64d59cff19414005f9f602b3721cef)，我们首先添加 `Init` 函数在 `存储器` API。

我们添加此函数来解决 `Storager` init 问题：我们需要配置Storager。 例如：我们需要为所有存储器设置 `WorkDir` (过去是 `Basel`)。

因此，我们的接口转向：

```go
// Init 将嵌入存储器本身。
//
// 呼叫者:
// - Init MUST 创建后被调用。
Init(配对...*类型。配对) (错误)
// InitWidContext 将会嵌入存储器本身。
InitWidext(ctx 上下文)。上下文，对...*类型。配对) (错误)
```

我们内部逻辑包括三个部分：

```go
// 1. 创建一个新服务
srv, err = azblob。新建(opt...)
如果是err != nil {
    return
}
名称, 前缀 := 命名空间.ParseObjectStorage(ns)
// 2. 获取一个 Storager 表单服务
商店，err = srv。Get(name)
如果err != nil {
    return
}
// 3。 配对的Init Storager。
err = 商店。Init(pairs.WewWorkDir(前缀))
如果err != nil {
    return
}
```

如果这里没有服务人员，那么init 逻辑将会变成两个部分：

```go
存储 = fs。新()
路径 := 命名空间.ParseLocalFS(ns)
err = store。Init(pairs.WewWorkDir(路径))
如果err != nil {
    return
}
```

看起来我们解决存储器的 init 问题，但不是真的。  有以下问题。

- `Init` 仅用于 `core utils`

如果一个 API 仅用于内部包，我们为什么要导出它？

- 只支持 `WorkDir` ，很难添加更多对。

现在，我们已硬码 `商店。Init(pairs.WewWorkDir(路径))` in `core utils`. 首先，如何添加更多的配对？ 然后，如果只有 `WorkDir` 需要的话，在 `Init` 中做这件事是昂贵的。

- 所有 `Storager` 以同样方式实现相同的接口。

如果实现是相同的，我们不需要将其导出为一个接口。

- `Init` 可以调用时间并可能导致并发问题

用户可以在 `邮件列表`期间更改 `WorkDir` ，不被允许。

## 建议

因此，我提议作如下修改：

- 将 `Init` 和 `新存储` 合并到 `Storager.init(配对...*类型)。配对)`
- 将 `新设` 重设为 `新设(配对...*类型)。配对) (srv *Service, store *Storage, err error)`
- 添加 `Init` 关联到 `Servicer`的 `获取` 和 `新的`
- 重整配置字符串选项句柄。

## 理由

### 重整配置字符串选项句柄。

配置字符串类似：

`qingstor://hmac:<access_key>:<secret_key>@<protocol>:<host>:<port>/<bucket_name>/<prefix>`

最初的设计预计会在用户输入和配置字符串解析之间达到平衡。 让我们考虑以下两种方式：

- `fs://?work_dir=/path/to/dir` and `azblob://hmac:<access_key>:<secret_key>?name=<bucket_name>&work_dir=<prefix>`
- `fs:///path/to/dir` and `azblob://hmac:<access_key>:<secret_key>/<bucket_name>/<prefix>`

很明显，第二种风格的用户输入较少。提议 [3-support-service-init-via-config-string](./3-support-service-init-via-config-string.md) 具有相同的想法。 然而，这种风格在不同类型的服务之间造成了问题。

- `fs:///path/to/dir`
- `azblob://hmac:<access_key>:<secret_key>`
- `azblob://hmac:<access_key>:<secret_key>/<bucket_name>`
- `azblob://hmac:<access_key>:<secret_key>/<bucket_name>/<prefix>`

我们需要服务类型来区分 `/path/to/dir` and `/<bucket_name>/<prefix>`这使配置字符串解析器难以实现。 我们必须添加诸如 `命名空间` 之类的概念，以便我们可以在 `核心工具`中推迟这项工作。

回到第一种风格，如何将此工作委托给最终用户？

- `fs://?work_dir=/path/to/dir`
- `zblob://hmac:<access_key>:<secret_key>?name=<bucket_name>&work_dir=<prefix>`

现在， `name` and `work_dir` 都可以在没有服务类型的情况下解析。

## 兼容性

呼叫服务 `新` 或 `Init` 的用户将直接面临破损的更改。

## 二． 执行情况

大多数工作将由本提案的作者完成。
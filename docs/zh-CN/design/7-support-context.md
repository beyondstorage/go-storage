---
author: Xuanwo <github@xuanwo.io>
status: 完成
updated_at: 2020-01-08
---

# 提议：支助情况

## 二. 背景

Golang广泛使用了上下文，越来越多的项目依靠上下文来完成最后期限。 取消或其他价值，包括gRPC、任选、GCS SDK等。

为 [存储](https://github.com/Xuanwo/storage) 添加上下文支持将使它成为一个为实体世界准备好的存储层：

- 允许设置操作截止日期
- 取消客服请求
- 支持跟踪
- ...

## 建议

因此，我提议作如下修改：

- 为每个公共API添加环境配对支持
- 为每个公共的 API 添加 `ReadWContext` 样式方法
- `读取在Context` 调用 `读取` 含有 `上下文` 配对
- `通过配对 读取` 使用 `Context`

以下是更详细的变化：

### 为每个公共API添加环境配对支持

我们将 `上下文` 作为预定义的 API 对，请确认 `上下文` 是在生成的代码中提供的：

```go
v, ok = values[pairs.如果没关系，则

    结果。上下文=v.(上下文)。上下文有
其他
    结果。context = contextBackground()
}
```

如果对应的 `上下文` ，我们将使用 `上下文`，否则我们将创建一个新的上下文.

本节将更改 `内部/cmd/meta`。

### 为每个公共的 API 添加 `ReadWContext` 样式方法

为每种方法添加 `XxxWContext` API , 例如:

```go
输入 Mover 界面
    Movement(src, dst string, pairs ...*types.配对) (错误)
}
```

将转到：

```go
输入 Mover 界面
    Movement(src, dst string, pairs ...*types.配对) (err 错误)
    移动内容(ctx context上下文，src，dst字符串，对...*类型。配对) (错误)
}
```

这项行动将由手工执行。 我们在这里没有太多的接口，所以没有别的东西来写一个工具。

### 为服务生成 XxxWContext API

让我们生成代码来存档。

首先，我们需要添加 `...*类型。配对` 让每一个 API 携带上下文，这将影响两个API：

- `元数据`
- `统计`

他们不需要，但也有意义：他们都可以叫作API。

然后，我们需要生成 `XxxWiContext` API，这样实现者就不需要关心它。

在生成的代码中，我们将做以下事情：

```go
func (s *存储) ReadWeltext(ctx context上下文，路径字符串，对...*类型。Pair) (r io.读者更近，错误) 电子邮件：
    span, ctx := opentration。StartSpanFromContext(ctx, "github.com/Xuanwo/storage/services/qingstor.存储器。阅读")
    延迟间隔。Finish()

    配对 = append(pairs, ps.Welcome (ctx))
    返回 s.读取(路径，对等...)
}
```

### 服务应处理上下文内容

服务应使用解析对应的环境：

```go
it := s.bucket对象(context.TODO(), &gs.Query{
    Prefix: rp,
})
```

应该更新到：

```go
it := s.bucket对象(选定)上下文， &gs.Query{
    Prefix: rp,
})
```

## 理由

其它执行方式可以如下：

- 在 API 中添加 ctx 并添加 `XXXWOUTContext` API 支持
- 在 API 中添加 ctx ，让不关心上下文的用户使用 `上下文TODO()`
- 添加环境配对支持且不触摸公共接口

### 在 API 中添加 ctx 并添加 `XXXWOUTContext` API 支持

接口将如下所示：

```go
输入 Mover 界面
    移动 (ctx 上下文)。上下文，src，dst字符串，对...*类型。配对(错误)
    移动退出(src, dst 字符串, 配对...*类型。配对) (错误)
}
```

首先，此更改是一个间歇性更改，每次API调用都需要重新设定。

然后，很明显， `移动退出上下文` 比移动上下文</code> 长了 `移动上下文`

最重要的事情是API设计背后的思考： **公平**。

这里有两种开发者：其中一些需要上下文支持，另一些则不关心它。

建议中的设计对两者都是友好的：

- 需要上下文支持的人应该使用 `XxxWiContext` 或在 `Xxx` 中添加上下文配对，他们知道他们在做什么。
- 不需要上下文支持的人可以在没有任何关于上下文的想法的情况下愉快地编写代码。

然而，这种设计对于不需要上下文支持的人来说是不公平的。 他们需要使用 API，如 `MoveWeWoutContext` ，尽管他们不关心上下文了。

### 在 API 中添加 ctx ，让不关心上下文的用户使用 `上下文TODO()`

前一章所述的类似理由。

### 添加环境配对支持且不触摸公共接口

看起来没问题，但有点不便。 这个设计使得很难添加追踪支持。 人们需要把代码换成：

```go
func ReadWeltext(ctx context)。上下文，s *存储，路径字符串，对 ...*类型。Pair) (r io.读者更近，错误) 电子邮件：
    span, ctx := opentration。StartSpanFromContext(ctx, "Read")
    defer span.Finish()

    配对 = append(pairs, ps.Welcome (ctx))
    返回 s.读取(路径，对等...)
}
```

为什么不让我们自己这样做？

## 兼容性

无间断变化

## 二． 执行情况

大多数工作将由本提案的作者完成。
---
author: Xuanwo <github@xuanwo.io>
status: 完成
updated_at: 2019-11-15
---

# 建议：统一存储器行为

## 二. 背景

We provide a `Capable` function for developers to check whether underlying storager support action/pair or not. 事实上，没有人使用这些武器。 作为 `统一的存储层`, 它奇怪的编码器仍然需要处理存储器服务的差异。

然而，确实存在着差异。 要么我们放弃某些存储器服务的能力，要么为开发者提供更好的方法来处理复杂性问题。 我们需要消除基本行动不一致的行为，并为用户使用高层行动提供方便的方式。

## 建议

因此，我提议作如下修改：

### 拆分基础存储

为当前 `存储器` 接口拆分基础存储器，并确保每个存储器都有相同操作和对等的操作。 如果存储不支持特定对，存储需要忽略它们。

基础 `存储器` 可以是：

```go
输入Storager 接口 然后
    String() 字符串

    Init(配对...*类型。配对) (err 错误)
    能力 (动作字符串，密钥...字符串) 布尔
    Metadata(类型)。元数据，错误)

    ListDir(路径字符串，对...*类型。配对程序迭代器。ObjectIterator
    读取(路径字符串，对...*类型。Pair) (r io.读者更近，错误)
    写入(路径字符串，r io)。读者, 配对 ...*类型。配对) (错误)
    状态(路径字符串，对...*类型。配对) (o *类型对象，错误)
    删除(路径字符串，对...*类型。配对) (错误)
}
```

### 拆分高层功能

高级别函数如 `复制`, `移动` 应该分割。 调用者在使用之前应先输入断言：

```go
键入Copier 界面 然后
    Copy(src, dst string, pirs ...*typesPair) (err error)
}

if x, ok := store.(Copier); ok @un.org
    err := x.Copy(老路径，新路径)
    如果是err != nil {
        return err
    }
}
```

将有一个 `段` 接口，用于所有相关的段操作。

```go
键入Segmenter interprete v.
    ListSegments(路径字符串, 对...*类型)配对程序迭代器。片段迭代器
    InitSegment(路径字符串，对...*类型。配对) (id 字符串，错误)
    写入部分(id 字符串，偏移，大小，int64，r io)。读者, 配对 ...*类型。配对) (err 错误)
    CompleteSegment(id 字符串，对...*类型。配对) (err 错误)
    中止 (id 字符串, 对 ...*类型。配对) (错误)
}
```

## 理由

### 为什么复制接口

为了知道存储器是否支持指定的 API 调用，我们需要通过某种方式传送此信息。 据我所知，目前有以下方式将其存档：

- 拷贝接口
- 可复制有趣的来电
- 复制能力
- 视场/不支持错误

我们将逐一讨论这些问题，并为此做一个基准。

对于Copier 接口

```go
键入Copier 界面 然后
    Copy(src, dst string, pirs ...*typesPair) (err error)
}

if x, ok := store.(Copier); ok @un.org
    err := x.Copy(老路径，新路径)
    如果是err != nil {
        return err
    }
}
```

我们为不同的能力创建不同的接口，呼叫者需要为他们打开。

可复制的有趣通话的

```go
输入存储接口 *    
    复制(src, dst字符串, 对...*类型Pair) (err error)
    Copyable() bool


如果存储。Copyable() {
    err := x.Copy(老路径，新路径)
    如果是err != nil {
        return err
    }
}
```

添加不同的 `XXXable` 真空调用，返回不同能力的布尔值，呼叫者需要先检查返回值后才能使用。

复制能力

```go
如果存储，请键入存储界面    
    Capability() Capability
}

Capability() & 类型CapabilityCopy == 1 电子邮件：
    err := x。Copy(老路径，新路径)
    如果是err != nil {
        return err
    }
}
```

`Storager` 支持返回一个 uint64 或表示能力的东西，呼叫者需要在使用前检查此值。

视场/不支持错误

```go
func(store *Storager) Copy() {
    panic("not supported")
}
```

`Storager` 会因不支持的函数出现恐慌或返回错误，呼叫者需要 `恢复` 或使用后检查错误。

基准文件可以在这里找到 [](./1/main_test.go)，结果看起来很好：

```go
goos: linux
目标: amd64
pkg: github. om/Xuanwo/storage/docs/design/1
BenchmarkCopierInterface-8 141224794 7。 9 ns/op
BenchmarkCopyableFuncall8 405428164 2。 4 ns/op
BenchmarkCopyCapability-8 512795942 2。 1 ns/op
基准错误8 48293546 ns/op
BenchmarkPanic-8 16575105 67。 n/op
PASS
```

很明显， `能力` 是最快的，恢复的恐慌是最慢的。

一方面，我们还需要考虑API是易于使用还是易于实施。

Func 调用很容易使用，但开发人员应该在他们的结构中添加另外两个功能。 能力很容易实现，但是增加了另一个概念，让呼叫者能够理解。 接口要清楚得多，但在实际使用中看起来有点微妙。

接口使用的导入越多，我们可以在真正使用前强制使用来检查存储器支持的呼叫器。

因此，我们选择接口，让它具有中心性能和强制性能力。

## 兼容性

### 复制、移动和范围

复制、移动和达标将从 `Storager` 接口中被删除，呼叫者在使用前应先键入电量。

### 与部门有关的业务。

所有与节点相关的操作都将从 `存储器` 接口中删除，呼叫者在使用之前应该先键入 `段` 的申明类型。

## 二． 执行情况

大多数工作将由本提案的作者完成。
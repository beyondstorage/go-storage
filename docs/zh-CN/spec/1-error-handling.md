---
author: Xuanwo <github@xuanwo.io>
status: 草稿
updated_at: 2020-02-18
added_by:
  - 设计/11错误处理.md
---

# 旁观：错误处理

这个速度将描述如何处理 [存储](https://github.com/Xuanwo/storage) 中的错误。

## 定 义

- `错误`: 程序运行不正确
- `包`: 所有有效的插件在 [存储](https://github.com/Xuanwo/storage) 中
- `实现者`: 执行 `包`
- `通话者`: 使用/调用 `包`

## 错误

从 [存储](https://github.com/Xuanwo/storage)的错误可以归类为以下类型：

- 预期错误
- 意外错误

预期的错误是实现者预期的错误。 对这些错误的定义应有足够的评论，对它们的任何改动都应记录在案。

意外的错误是实现者无法预料的错误。 这些错误在依赖性提升且没有更新的情况下会被改变或消失。

错误应该始终作为一个结构来表示，它会带有上下文错误信息。 依赖于软件包实现。软件包可能有一个以上的错误。

```go
类型错误结构如下：
    Op 字符串
    Err 错误

    ContextA 字符串
    ContextB structextB
...
}
```

- `Op` 是指触发此错误的操作。
- `错误` 带有潜在错误
  - 对于预期的错误，相关错误将直接使用
  - 对于意外的错误，应按照或改写传递错误
- `ContextX` 带有上下文错误信息，每个上下文都应该实现 `String() 字符串`

每个错误构建SHOULD 实现以下方法：

- `Error() 字符串`
- `Unwrap() error`

返回了 `Error()` SHOULD 格式相同：

`{Op}: {ContextA}, {ContextB}: {Err}`

`打开` SHOULD 总是返回没有任何操作的潜在错误。

## 实现

本节将描述软件包实现者侧的处理错误。

- 预计错误属于已声明的软件包，只有这个软件包CAN返回了这个错误
- 实现者在确保此操作无法继续或此操作会影响或销毁数据且无法恢复时恐慌。

## 呼叫者

本节将描述处理包调用方块时出现的错误。

- 呼叫者SHOULD 只检查包的预期错误，不检查通过包导入的标签返回的错误
- [存储](https://github.com/Xuanwo/storage)的软件包 CAN 恐慌，操作无法继续，召唤SHOULD 自行恢复

## 示例

预期错误

```go
var (
    // 错误不支持的协议将返回协议。
    ErrUnsupportedProtocol = errors.新("不支持的协议")
    // 错误值表示值无效。
    ErrinvalidValue = 错误。新("无效值")
)
```

结构错误

```go
// 错误表示与端点相关的错误。
类型错误结构是否为
    Op 字符串
    Err 错误

    协议字符串
    值 []字符串
}

func (ae *Error) Error() 字符串否
    if e.values == nil 。
        return fmt。Sprintf("%s: %s: %s", e.Op, e.Protocol, e.ErrError())
    }
    返回 fmt。Sprintf("%s: %s, %s: %s", e.Op, e.Protocol, e.数值，e.ErrError())
}

// 卸载实现x错误。包装器
真空(e *Error) Unwrawrawraw() 错误。
    返回 e。错误
}
```

预计发生错误

```go
err = &Error{"parse", s[0], nil, ErrUnsupportedProtocol}
```

发生意外错误

```go
端口，err := strconv.ParseInt(s[2], 10, 64)
if err != nil 。
    return nil, &Error{"parse", ProtocolHTTP, s[1:], err}
}
```

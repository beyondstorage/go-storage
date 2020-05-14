---
author: Xuanwo <github@xuanwo.io>
status: 完成
updated_at: 2020-02-18
adds:
  - spec/1错误处理.md
---

# 建议：处理错误

## 二. 背景

[存储](https://github.com/Xuanwo/storage) 旨在成为一个生产现成的存储层，错误处理是最重要的部分之一。

当面临错误时，用户应该能够关注以下事：

- 知道发生了什么错误
- 决定如何处理问题
- 挖掘为什么发生错误

为了提供这些能力，我们应该以背景信息返回错误。

## 建议

因此，我提议作如下修改：

- 添加 spec [旁观：处理](../spec/1-error-handling.md) 时出错，以便正常化整个lib 中的错误处理
- 处理相关代码的所有错误都应该被重新设置

Take package `segment` as example, as described in [spec: error handling](../spec/1-error-handling.md), we will have an `Error` struct to carry contextual information:

```go
输入SegError struct v.
    Op 字符串
    Err 错误

    Seg *Segment
}

func (ee *SegError) Error() string Windows
    return fmt.Sprintf("%s: %v: %s", e.Op, e.Seg, e.Err)
}

func (e *SegError) Unwrap() error {
    return e.错误
}
```

所以 `片段` 可能会返回如下错误：

```go
func (s *Segment) ValidateParts() (err error) {
    ...

    // 零件不被允许，导致无法完成。
    if len(s.Parts) == 0 v.
        return &Segerror{"validing parts", ErrSegmentPartsEmpty, s}


...
}
```

然后调用者可以检查这些错误：

```go
err := s.ValidateParts()
```

如果我们不关心由 `片段` 返回的错误：

```go
如果是err != nil {
    return err
}
```

如果我们想要处理一些状态：

```go
如果err != nil && 错误Is(err, partErrSegmentPartsEmpty) {
    log.打印("片段为空")
}
```

如果我们想要获取更多详细信息错误信息：

```go
var e 格错误
如果err != nil && 错误。As(err, &e) vol
    log.Print(e.片段)
}
```


## 理由

- [https://blog.golang.org/error-handling-andgo](https://blog.golang.org/error-handling-and-go)
- <http://joeduffyblog.com/2016/02/07/the-error-model/>
- [https://blog.golang.org/go1.13错误](https://blog.golang.org/go1.13-errors)

## 兼容性

[返回的错误](https://github.com/Xuanwo/storage) 可以更改。

## 二． 执行情况

大多数工作将由本提案的作者完成。

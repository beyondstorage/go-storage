---
author: Xuanwo <github@xuanwo.io>
status: 完成
updated_at: 2020-02-09
---

# 建议：回拨阅读器

## 二. 背景

我们确实需要进度报告能力：

- 为了用户互动，我们需要进度条通知他们当前状态
- 为了服务互动，我们需要推/拉进度状态

在这两种情况下，我们都需要有能力编写进度报告。

## 建议

因此，我提议作如下修改：

- 在 `pkg/iowrap 中添加回调读者/读者`
- 添加 `ReadCallbackFunc` 到 `类型/配对`
- 添加 `ReadCallbackFunc` 对拥有IO 操作的任何方法
- 把 `ReadCallbackFunc` (如果 `在每个服务中都有`)

在完成所有这些工作之后，我们可以在取得进展的情况下很好地开展工作：

```go
导入 (
    "io/ioutil"
    "log"

    "github.com/schollz/progressbar/v2"

    "github.com/Xuanwo/storage/coreutils"
    ps "github.com/Xuanwo/storage/types/pairs"
)


func main()
    // Init a service.
    _, store, err := coreutils.Open("fs:///?work_dir=/tmp")
    如果是err != nil @un.org
        log。Fatalf("service init failed: %v", err)
    }

    bar := progressbar.新建(16* 1024 * 1024 * 1024
    延迟酒吧.Finish()

    r, err := store。已读("test_file", ps.WithReadCallbackFunc(func(b []byte) {
        bar.Add(len(b))
    })
    如果是err != nil 。
        log。Fatalf("服务读取失败： %v", 错误)
    }

    _, err = ioutil。ReadAll(r)
    如果是err != nil 。
        log。Fatalf("ioutil 读取失败: %v", 错误)

}
```

## 理由

### 内置进度条与回调

我们是否应该提供一个内置的进度条？ 不，我认为没有。

进度条远远超过了打印行尾的 `/r`。 我们应该注意窗口宽度，字符宽，屏蔽主题/样式/外观，线程安全，多行支持，多平台支持等等。 作为一个存储lib，我们需要关注存储级别，不要触及这些工作。

相反，我们需要提供一种机制，使每个进展都能很好地发挥作用。 据我所知，Read()中的回调可能是一个好的选择。

### 回调函数(int) vs function([]byte)

读取有两个回调选项： `func(int)` and `function ([]byte)`。

它们的基准没有很多不同：

```go
// Fist 运行
goos: linux
goarch: amd64
pkg: github. om/Xuanwo/storage/docs/design/10
BenchmarkPlainReader-8 269910 4149 ns/op 987。 7 MB/s
BenchmarkIntCallbackReader-8 292818 4188 ns/op 977。 8 MB/s
BenchmarkBytesCallbackReader-8 265878 4176 ns/op 980。 0 MB/s
PASS

// 第二次运行
goos: linux
goarch: amd64
pkg: github. om/Xuanwo/storage/docs/design/10
BenchmarkPlainReader-8 244312 4456 ns/op 919。 2 MB/s
BenchmarkIntCallbackReader-8 262990 4202 ns/op 974。 4 MB/s
BenchmarkBytesCallbackReader-8 240216 4290 ns/op 954。 4 MB/s
PASS
```

结果不稳定, 比较阅读操作, 额外的有趣通话不会消耗太多的 CPU 时间。

在此基础上，返回 `n` 缺少可扩展性，所以我选择返回 `[]byte` 代替。

## 兼容性

没有中断的更改。

## 二． 执行情况

大多数工作将由本提案的作者完成。
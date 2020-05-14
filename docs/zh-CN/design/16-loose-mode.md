---
author: Xuanwo <github@xuanwo.io>
status: 完成
updated_at: 2020-03-12
---

# 建议：松散模式

## 二. 背景

当前 [存储](https://github.com/Xuanwo/storage)的配对处理行为是不一致的。

在所有 `parseStoragePairXXX` 函数中，我们将只通过选择支持的函数忽略不支持对：

```go
v, ok = values[ps.DirFunc]
如果没关系，
    结果。HasDirFunc = true
    结果。DirFunc = v.(类型)。对象函数)
}
```

但在其他对相关的逻辑中，例如 `storage_class` 支持，我们也返回了错误：

```go
func parseStorageClass(in storageclassType) (string, error) {
    switch in {
    case storageclass.热:
        返回 storageClassStandard, nil
    cases storageclass.警告：
        返回 storageClassStandardIA, nil
    默认值：
        返回 ""， &服务。PairError{
            Op:    "parse storage class",
            Err:   services.错误储存类支持,
            配对: []*类型。配对 {Key: ps存储分类，值：in}，
        }

}
```

所以用户可能会混淆我们如何处理与兼容性相关的问题。

## 建议

所以我提议所有服务都要有 `个松散的` 模式。 `loose` mode will be `off` as default, and services will return error when they reach incompatible place. 当 `溢出` 开启时，所有不兼容的错误都将被忽略。

例如：

我们有一个 Storager 不支持 `大小` 对 `读取`。

`丢失` on : 此错误将被忽略. `松散` off: Storager 返回一个兼容性相关错误。

目前，我们在 `配对错误` 中混合了兼容性错误和其他配对相关错误。 我们将添加两个不同的错误： `错误能力不足` and `错误限制不满意`。

`错误能力不足` 意味着此服务不具备此功能， `错误限制不满意` 表示此操作不会限制服务。 `错误能力不足` 可能会被安全忽略，如果您对服务行为一致性不关心很多，并且会在松散模式下被忽略。

基于这些错误，我们将会有新的错误结构，例如 `配对错误` 来传递错误环境：

```go
// NewPairRequired Error 将创建一个新的配对错误.
func NewPairRequired Error(key...字符串) *PairRequired Error(
    )
 return &pairRequired Error{
        Err:  ErrRestrictionDissatisfied,
        Keys: keys,
    }
}

// PairRequired Error 表示此操作需要配对但缺失。
类型配对错误结构变化。
    Err 错误

    Keys []string
}

func (ae *PairRequireedError) Error() string Windows
    return fmt.Sprintf("需要配对" %v: %s", e。Keys, e.ErrError())
}

// 卸载实现x错误。包装器
真空(e *PairRequired Error) Unwrawrawrapper () 错误
    return e.错误
}
```

## 理由

无。

## 兼容性

- 更多 `错误能力不足` 可以返回，因为 `松散了` 模式将被默认打开
- 有些错误可以作为其他错误结构而不是 `配对错误`
- `配对错误` 将被删除

## 二． 执行情况

大多数工作将由本提案的作者完成。
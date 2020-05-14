---
author: Xuanwo <github@xuanwo.io>
status: 完成
updated_at: 2019-11-18
updated_by:
  - 设计/12-support-bot-directory-and prefix-based list.md
---

# 建议：在列表操作中使用回调

当前的 API 设计导致了一些有线结果。 其中一个案件看起来很好：

```go
它 := 商店。ListDir(路径，类型)。Web Recursive(true))

为您的
    o, err := 它下一步()
    如果err != nil && 错误。Is(err, 迭代器)错误) {
        break
    }
    如果err != nil &
        t.触发器(类型)。新错误未处理(错误))
        返回
    }
    商店。Delete(o.姓名)
}
```

为了在路径下移除所有对象，我们显然需要递归列出它们，并逐一删除。 这个算法对对象存储、a.k.a很有效。 基于前缀的存储。 它将在 POSIX 文件系统上失败，因为我们的 `List` 不返回文件夹。

## 建议

因此，我提议作如下修改：

### 为所有列表操作添加回调func

列表操作将不返回 `迭代器。ObjectIterator` 再者，它将允许输入 `WandirFunc` 和 `WiFileFunc` 通过 `Pair`

```go
dirFunc := function(对象 *类型)。Object) 让您能够
    打印出("dir %s", object.名称)
}
fileFunc := func(对象 *类型)Object) 让您能够
    printf("file %s", object.名称)
}

err := 商店。列表("前缀", 类型)。WidDirFunc(dirFunction)类型.Wide FileFunc(fileFunc))
如果err != nil {
    return err
}
```

### 删除列表中的递归项

基于邮件列表</code>中的 `回调工作, 我们可以移除列表中所有递归支持。</p>

<p spaces-before="0">基于目录的存储将只列出一个目录，而基于前缀的存储将只列出一个没有分隔符的前缀。</p>

<h2 spaces-before="0">理由</h2>

<h3 spaces-before="0">为什么不递归删除？</h3>

<p spaces-before="0">如果我们在删除中支持 <code>递归` ，删除目录或删除前缀可能更简单：

```go
商店。删除(路径，类型)。Werrecsive(true))
```

在一个基于前缀的存储服务中，若要支持删除递归，我们需要在 `删除`中列出对象， 这使我们的 API 不是正本，很难实现，并且很难进行单元测试。

为了保持所有API的同一样式， 我们还需要在 `复制`、 `移动` 等中添加递归支持。

除上述原因外，还包括下列原因： 在 `删除` 中实现递归化使我们的 API 无法在并行任务框架中使用：它们不能再被分成子任务。

### 为什么不在列表中返回 Dir

另一个想法可以同时返回 `Dir` and `File` 在 `List` 而不是 `File` 这个方法可能导入两个问题。

首先，虽然我们可以从 `List`获取 `Dir` ，但这并不能解决问题。 调用者需要确保此 `Dir` 下的每一个 `对象` 已被删除。 他们需要像这样代码：

```go
它 := 商店。ListDir(路径，类型)。撤销(true))
dirers := make([]*类型。Object, 0)

for Group
    o, err := it.下一步()
    如果err != nil && 错误。Is(err, 迭代器)错误) {
        break
    }
    如果err != nil &
        t.触发器(类型)。新错误未处理(错误))
        返回
    }
    如果o。类型 == 类型。ObjectTypeDir {
        dirs = append(dir, o)
    } else {
        store.Delete(o.姓名(名称)
    }    
}

for i:=len(dirs)-1;i>=0;i-- *
    store.删除(目录[i])。姓名)
}
```

这么丑陋。

第二，它也使得存储器更难实现这些目标并对它们进行测试，而这并不是我们的目标。

## 兼容性

在列表操作中不再循环，相反，呼叫者需要输入一个有趣的配对。

## 二． 执行情况

大多数工作将由本提案的作者完成。

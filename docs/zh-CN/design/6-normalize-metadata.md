---
author: Xuanwo <github@xuanwo.io>
status: 完成
updated_at: 2020-01-03
---

# 建议：使元数据正常化

## 二. 背景

元数据正在发生混乱。

目前，我们有以下元数据：

```go
const (
    校验和 = "checksum"
    Class = "class"
    Count = "count"
    到期时间 = "expire"
    Host = "host"
    Location = "location"
    size = "size"
    Type = "type"
    WorkDir = "work_dir"
)
```

为了统一不同服务的值名称，我们有一个虚拟名称映射：

- 校验和 -> 内容MD5
- 类型 -> 内容类型
- 类 -> x-qs-storage-class，x-amz-storage-class...
- 大小 -> 桶大小
- 计数 -> 桶计数
- 位置 -> 桶的位置或区域

在更多值提示构造值后，这些名称映射可能会混淆终端用户：

- `对象。大小` 或 `对象。GetSize()` ?
- `对象。输入` 或 `对象。GetType()` ?
- 不同的服务可以使用不同的散列算法，如何区分它们？

## 建议

因此，我提议作如下修改：

### 拆分存储元和对象元数据

- 使用 `映射[string]interface{}` 替换 `元数据`并直接在 `存储` 生成函数。
- Rename `metadata.存储` 到 `元数据。存储元`
- 添加 `元数据。ObjectMeta`
- 添加 `元数据。存储统计`
- 将 `metadata.json` 拆分到 `object_meta.json`, `storage_meta.json` 和 `storage_statistic.json`
- 重建 `元数据` 代码生成器
- 使 `设置XXX` 返回本身，所以我们可以将他们称为一个链：

    ```go
    m := 元数据。NewObjectMeta().
        SetContentType("application/json").
        SetStorageClass("冷")。
        SetETag("xxxxx")
    ```

### 规范化元数据的名称

我们应该通过以下方式实现元数据名称正常化：

**如果此元数据输出为 `对象` 结构值，忽略它们。**

**如果在 [消息头](https://www.iana.org/assignments/message-headers/message-headers.xhtml)中定义了这个元数据，请使用规范风格。**

例如， `content-md5`。

*因为HTTP-2已在 [rfc7540](https://tools.ietf.org/html/rfc7540)中获得批准，我们将使用页眉的小写风格。*

**如果此元数据是私密元，使用最常见的用法。**

例如， `x-qs-storage-class` and `x-amz-storage-class` should be meta `storage-class`

### 规范化元数据的值

这些元值也应正常化：每个服务都应将自己的元转换为储存共享元值。

例如，Amazon S3有以下存储类别： `STANDARD`, `REDUCED_REDUNDANCY`, `STANDARD_IA`, `ONEZONE_IA`, `GLACIER`, `DEEP_ARCHIVE` 等等。

在 Azure 存储中，他们有不同的 `AccessTier`: `AccessTierHot`, `AccessTierCool`, `AccessTierArchive`

我们需要处理所有这些问题：

- 保存它们为元: `存储类`
- 从服务中获取时转换为共享元数据
- 更新时转换为专用元数据

## 理由

无。

## 兼容性

- 元数据的名称和值可以更改。
- 存储将返回不同的元数据类型。

## 二． 执行情况

大多数工作将由本提案的作者完成。
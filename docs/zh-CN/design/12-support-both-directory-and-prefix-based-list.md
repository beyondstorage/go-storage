---
author: Xuanwo <github@xuanwo.io>
status: 完成
updated_at: 2020-03-02
updates:
  - design/2-use-callback-list-operations.md
---

# 建议：支持基于目录和前缀的列表

## 二. 背景

我们删除了在 [2-use-callback-list-operations](./2-use-callback-in-list-operations.md)中的递归支持，我们在同一建议中概述了：

> 基于目录的存储将只列出一个目录，而基于前缀的存储将只列出一个没有分隔符的前缀。

它很容易实现，但是不能满足需求。 主要有两种假设情况：

- 与基于前缀的存储服务目录合作，例如同步文件
- 没有分隔符的列表，这样我们就可以更快地列出目录了

我们可以看到， 前一项提案修正了清单，但没有划界器情景，但消除了在基于前缀的存储服务上与目录共用的可能性。

## 建议

因此，我提议作如下修改：

- 为基于列表的前缀添加 `ObjectFunc`
- 将 `FileFunc` 和 `DirFunc` 视为基于目录的列表

存储器的行为将遵守规则：

- `ObjectFunc` 不能传递到 `FileFunc`
- `ObjectFunc` 无法通过 `DirFunc`
- `需要基于目录列表的支持`
- `基于列表的前缀支持` 是可选的

我们将在 `内部/cmd/service`中进行验证，所以实现者不需要检查它们。

## 理由

无

## 兼容性

没有中断的更改。

## 二． 执行情况

大多数工作将由本提案的作者完成。
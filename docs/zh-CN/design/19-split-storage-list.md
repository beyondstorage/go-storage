---
author: Xuanwo <github@xuanwo.io>
status: 候选项
updated_at: 2020-04-09
updates:
  - design/2-use-callback-list-operations.md
deprecates:
  - 设计/12-support-bot-directory-and prefix-based list.md
---

# 建议：分割存储列表

## 二. 背景

提议 [支持基于目录和前缀列表](./12-support-both-directory-and-prefix-based-list.md) 已被实践证明是一个失败。 在这个建议中，我们对基于列表的前缀支持使用 `ObjectFunc` 并对 `FileFunc`, `DirFunc` 和 `ObjectFunc` 添加许多限制。 问题是用户不知道该存储服务是基于前缀还是基于目录。 所以它们总是回退到基于目录列表的方法，而这种方法对于对象存储服务来说不是不可吸附的。

## 建议

因此，我提议作如下修改：

- 将 `邮件列表` 分成 `ListDir` 和 `邮件列表前缀`
- 从 `Storager 删除 <code>列表`</code>
- 添加接口 `DirLister` for `ListDir`
- 为 `List前缀` 添加接口 `前缀列表`

因此用户需要断言接口 `DirLister` 才能使用 `ListDir`。

同时，我们应当：

- 将 `列表部分` 重命名为 `列表前缀段` 以匹配前缀更改
- 从 `段删除 <code>列表段`</code>
- Add interface `PrefixSegmentsLister` for `ListSegments`

## 理由

无。

## 兼容性

所有 `List` 的 API 调用将被打破。

## 二． 执行情况

大多数工作将由本提案的作者完成。
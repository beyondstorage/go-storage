---
author: Xuanwo <github@xuanwo.io>
status: 完成
updated_at: 2020-03-09
---

# 建议：正常化内容散列检查

## 二. 背景

进行内容散列检查非常常见，在不同服务之间同步文件时尤其如此。 然而，不同的服务并不使用相同的内容散列算法。 例如：

- 大多数对象存储服务使用 `Content-MD5` 头来携带内容 md5 散列项
- 某些对象存储服务使用其唯一的算法，例如 [`kodo` etag散列](https://developer.qiniu.com/kodo/manual/1231/appendix#qiniu-etag)
- 对用户 SaaS 云存储服务总是有自己的散列算法，例如 [`dropbox` Content Hash](https://www.dropbox.com/developers/reference/content-hash)

因此，我们需要使内容散列检查行为正常化，以便我们能够安全和正确地比较不同服务之间的内容散列情况。

## 建议

因此，我提议作如下修改：

- 标准 `content-md5` SHOULD 被填入对象元数据 `content-md5` 被填入 `base64 of 128 位 MD5 digest 中，按照RFC 1864`
- 非标准 `content-md5` 头部被当作服务自定义内容哈希
- 服务自定义的内容散列将被填入对象元数据 `etag` 而不进行任何修改
- 如果服务只返回 `content-md5` ，它应该填入对象元数据 `etag`
- 对象元数据 `content-md5` CAN 安全地跨越服务
- 对象元数据 `etag` CAN 仅用于同一服务

## 理由

HTTP 相关标准

- [超文本转移协议 (HTTP/1.1): 语法和内容](https://www.rfc-editor.org/rfc/rfc7231)
- [超文本传输协议 (HTTP/1.1): 条件请求](https://www.rfc-editor.org/rfc/rfc7232)
- [永久消息头字段名称](https://www.iana.org/assignments/message-headers/message-headers.xml#perm-headers)

存储服务参考文档

- [`kodo` etag hash](https://developer.qiniu.com/kodo/manual/1231/appendix#qiniu-etag)
- [`dropbox` 内容哈希](https://www.dropbox.com/developers/reference/content-hash)

## 兼容性

无间断变化

## 二． 执行情况

大多数工作将由本提案的作者完成。
---
author: Xuanwo <github@xuanwo.io>
status: draft
updated_at: 2020-03-09
---

# Proposal: Normalize content hash check

## Background

It's very common to do content hash check, especially in syncing files between different services. However, different services don't use the same content hash algorithm. For example:

- Most object storage services use `Content-MD5` header to carry content md5 hash
- Some object storage services use their only algorithm, like [`kodo` etag hash](https://developer.qiniu.com/kodo/manual/1231/appendix#qiniu-etag)
- To consumers SaaS cloud storage services always have their own hash algorithm, like [`dropbox` Content Hash](https://www.dropbox.com/developers/reference/content-hash)

So we need to normalize the content hash check behavior so that we can compare content hash between different services safely and correctly.

## Proposal

So I propose following changes:

- Object metadata `content-md5` SHOULD filled with content md5 with normalization or keep empty
  - In string format without any `"` or `'`
- Object metadata `etag` SHOULD filled with services self defined content hash without any modification or keep empty
  - In string format and keep all `"` and `'`
- `content-md5` CAN be used safely across services
- `etag` CAN only be used in same service

## Rationale

HTTP Related Standards

- [Hypertext Transfer Protocol (HTTP/1.1): Semantics and Content]
- [Hypertext Transfer Protocol (HTTP/1.1): Conditional Requests]
- [Permanent Message Header Field Names](https://www.iana.org/assignments/message-headers/message-headers.xml#perm-headers)

Storage Service Reference Document

- [`kodo` etag hash](https://developer.qiniu.com/kodo/manual/1231/appendix#qiniu-etag)
- [`dropbox` Content Hash](https://www.dropbox.com/developers/reference/content-hash)

## Compatibility

No break changes

## Implementation

Most of the work would be done by the author of this proposal.

[RFC7232]: https://www.rfc-editor.org/rfc/rfc7232
[RFC7231]: https://www.rfc-editor.org/rfc/rfc7231
[Hypertext Transfer Protocol (HTTP/1.1): Semantics and Content]: https://www.rfc-editor.org/rfc/rfc7231
[Hypertext Transfer Protocol (HTTP/1.1): Conditional Requests]: https://www.rfc-editor.org/rfc/rfc7232
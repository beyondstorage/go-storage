---
author: Xuanwo <github@xuanwo.io>
status: 完成
updated_at: 2020-03-23
---

# 建议：返回段接口

## 二. 背景

[从第一次发布中存储](https://github.com/Xuanwo/storage) 个已添加的片段支持。 自那时以来，相关的 API 部分出现在以下方面：

```go
键入Segmenter interprete v.
    ListSegments(路径字符串, 对...*类型)配对) (错误)
    InitSegment(路径字符串，对...*类型。配对) (id 字符串，错误)
    写入部分(id 字符串，偏移，大小，int64，r io)。读者, 配对 ...*类型。配对) (err 错误)
    CompleteSegment(id 字符串，对...*类型。配对) (err 错误)
    中止 (id 字符串, 对 ...*类型。配对) (错误)
}
```

数据流看起来很好：

```text
                                                        +------------------------+
                                                        | |
                                                 +------>    完成 |
                                                 | | | |
                                                 | +-------------------------+
                                                 |
       +-------------------------+ | 
 + --------------+ |
       | | | | |
       | Init +-------->+ Write +----+
       | | | | | | | | | |
       +------------+ +-------------+ |
                                                 |
                                                 | +-------------+
                                                 | | |
                                                 +------>    中止|
                                                        |
                                                        +-------------+
```

在这个执行中，我们使用一个被调用的方法： `偏移部分`。 用户需要通过抵销将部分分割成不同的部分，并通过偏移和大小检查完整性。 但大多数对象存储服务使用 `基于索引的部分` 通过 `PartNumber`, 并且很难从 `偏移量转换为` 到 `基于`的索引， 我们需要确保每个部分都有同样的规模。

因此我们需要找到一种方式来支持基于 `偏移的` 和 `索引基于` 部分。

## 建议

因此，我提议返回部分接口而不是段号。

我们应该返回一个 ID 区块，而不是返回 `区块。段` 接口：

```go
请输入Segment interface volume
    String() string

    ID() string
    Path() string
}
```

开发者可以自行实现他们自己的区段逻辑，或使用基于内置的区段。

`段` 可以更改为:

```go
键入Segmenter interprete v.
    ListSegments(路径字符串, 对...*类型)配对) (错误)
    InitSegment(路径字符串，对...*类型。Pair) (seg segment.片段，err 错误)
    WrteSegment(seg 片段)。Segment, r io.读者, 配对 ...*类型。Pair) (err error)
    CompleteSegment(seg segment.片段，对...*类型。配对) (错误)
    中止 segment(seg区段).片段，对...*类型。配对) (错误)
}
```

`段段` 将由服务生成，不能在外部更改。

## 理由

### 基于分片段的 VS 偏移索引值

Index based segment needs `index` and `size` to mark a part, and offset based segment needs `offset` and `size`. 他们不能在没有相同部分大小的情况下转换。 但很难让用户在所有写部分通话中使用相同的部件大小。

为了实现它们，我们需要处理 `大小`， `偏移` 或 `索引` 为 `配对`, 服务可以决定使用哪一个。

这里还有其他实现方式：

1. `偏移BasedSegmenter` vs `IndexBasedSegmenter`
2. `WriteSegmentViaIndex` vs `WriteSegmentViaOffset`

不.1 执行的问题是，这些问题使得服务部门的工作量过多，不足以支助这两种服务。 未解决的问题 实现是不同的段逻辑需要不同的 `init` 和 `中止` 支持 仅支持不同的 `写入` 是不够的。 这也增加了太多的工作来实现 `段`。

### 段接口而不是段 id

首先，返回部分ID是一个BUG。 为了区分多部分上传，我们需要 `upload_id` 和 `object_key`。 我们可能对不同的 `object_key` 有相同的 `upload_id`。

第二，只有返回的部分ID才能迫使我们处理当地与偏远地区之间的不一致问题。

它适用于 `init->写入->完成`, 但如果我们试图中止所有部分会发生什么？ 我们需要 `列表 ->中止`。 `中止` 真菌是 `中止片段( id 字符串，对...*类型)配对) (错误)`。 我们必须从局部地图获取此段ID的相应路径。 但它们是空的，因为它们不是由 [存储](https://github.com/Xuanwo/storage) 发起的。 为了解决这个问题，我们必须更新本地段地图，而 `邮件列表`：

```go
为 _, v = 范围输出。上传者:
    seg := 部分。NepSegment(*v.注如果选项，上传 ID，0)

HasSegmentFunc {
        opt.SegmentFunc(seg)
    }

    s.segmentLock.Lock()
    // 更新客户端部分。
    s.sects[seg.ID] = seg
    s.segmentlock.Unlock()
}
```

这个问题可以通过返回一个分段接口来解决。 更改后，用户可以自己识别一个区段，局部地图已不再需要。

## 兼容性

所有关联的 API 部分都已更改。

## 二． 执行情况

大多数工作将由本提案的作者完成。
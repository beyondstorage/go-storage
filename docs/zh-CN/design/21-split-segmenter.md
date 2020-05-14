---
author: Xuanwo <github@xuanwo.io>
status: 候选项
updated_at: 2020-04-27
---

# 建议：分割段

## 二. 背景

在提议 [18-return segment-interface——改为]((./18-return-segment-interface-instead.md))，我们介绍了 `部分。Segmenter` 接口，但也引入了一个新问题：用户不知道哪种类型的服务支持。

In the `Rationale` section of [18-return-segment-interface-instead]((./18-return-segment-interface-instead.md)), we have already discussed two implementations:

1. `偏移BasedSegmenter` vs `IndexBasedSegmenter`
2. `WriteSegmentViaIndex` vs `WriteSegmentViaOffset`

我们说：

> 不.1 执行的问题是，这些问题使得服务部门的工作量过多，不足以支助这两种服务。

但也许这就是这样的办法，我们可以在这个想法上建立新的隔阂。

## 建议

基于 `偏移BasedSegmenter` and `IndexBasedSegmenter`的想法，我提议做以下修改：

添加新接口 `IndexSegmenter` ，实现以下函数：

- `InitIndexSegment(路径字符串，对...*类型)。Pair) (seg segment.片段，错误)`
- `写入 IndexSegment(seg部分)。Segment, r io.读者, 索引, int64, 大小, 对 ...*类型。配对) (错误)`
- `CompleteSegment(seg区段)片段，对...*类型。配对) (错误)`
- `中止片段 (seg区段)片段，对...*类型。配对) (错误)`

提取 `补全部分` and `中止部分` 到一个未导出的接口 `片段`。

嵌入 `片段` 到 `IndexSegmenter`, `DirSegmentsLister` 和 `前缀片段`

删除 `个片段`

### 添加新接口 `IndexSegmenter`

很容易发现我们不需要实现整个新的片段界面。 介绍了 `个部分之后.段` 接口，我们可以重新使用已经实现的 `补全段` 和 `中止分段`。

在这一变化之后，服务实施者的工作量也减少了。 如果需要实现新的片段方法，我们只需要添加 `偏移片段`。

### 提取未导出的接口并嵌入他人中

我们已经添加 `DirSegmentsLister` and `PrefixSegmentsLister`, 但只列出界面是没有用的，我们总是需要转换为 `部分`。 解压缩 `补全部分` and `中止部分` 到一个未导出的接口后，我们可以让 `DirSegmentsLister` 和 `前缀部分` 更有用。

### 删除 `个片段`

是的，我知道我在这里引入了一个 API 破坏性变化。

很抱歉，但我现在没有时间和精力来维持v1和v2。 既然 [存储](https://github.com/Xuanwo/storage) 几乎没有人只使用过我，我认为这里引入这种更改而不打破别人是安全的。

我意识到这个风险，我将对此负责。

## 理由

无。

## 兼容性

所有 API 调用相关的 `片段` 将会被生成。

## 二． 执行情况

大多数工作将由本提案的作者完成。
---
author: Xuanwo <github@xuanwo.io>
status: 候选项
updated_at: 2020-03-18
added_by:
  - 本记录包括中文发言的译文。
---

# 观感：建议

## 格式

建议将使用 `markdown` 格式。

## 标题

提案的文件名将如下所示：

`<proposal-number>-proposal-name.md`, 例如 `16-loose-mode.md`

他们应该通过 `16个松散模式` 或 `design/16-loose-mode.md`

## 元数据

建议中的元数据将作为预设元数据来传输更多信息。

举 [3-support-service-init-via-config-string](../design/3-support-service-init-via-config-string.md) 作为示例：

```yaml
---
作者：Xuanwo <github@xuanwo.io>
状态：完成
updated_at： 2019-12-23
updated_by
  - design/4-credit al-refacture. d
  - design/9-remove-storager-init.md
已弃用：
  - design/13-remove-config-string.md
-
```

`auther`, `状态` 和 `更新 _at` 是必需的。

`自动` 格式应该为： `名称 <Email>`

如果此建议影响到其他建议， `updated_by`, `updated`, `废弃的`, `废弃的` 等也应该添加。

## 状态

提案具有以下状态： `草稿`, `候选`, `已完成`。

- 提案刚刚创建，但尚未实现： `草稿`
- 提案已执行，但没有最后确认： `候选项`
- 提议已包含在一个次要版本中： `已完成`

因此，当我们发送了一个评论时，我们应该将此建议设置为 `草稿`。

在我们执行了一项建议之后，我们应该将此建议设置为 `个候选`在现阶段，我们可以根据实际反馈更新我们的执行情况。

当我们决定稍微释放时，我们应该考虑所有 `候选` 提案。 不满意的建议将被删除，接受的建议将被标记为 `完成`。

在标记为 `的提议完成了`后，我们不能再更改它的内容(如果更新，元数据也可以更改)。 我们需要提出一项新的建议，以改变某些提案的行为。

## 代码

相关的代码应该是带提案号的子目录。
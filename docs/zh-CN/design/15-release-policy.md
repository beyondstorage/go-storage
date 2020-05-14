---
author: Xuanwo <github@xuanwo.io>
status: 完成
updated_at: 2020-03-10
---

# 建议：发布政策

## 二. 背景

Golang团队确保 [Go 1的兼容性和未来的Go 程序](https://golang.org/doc/go1compat):

> 打算使写给1级规格的程序在该规格的整个生命周期内继续正确地编纂和运行。

[存储](https://github.com/Xuanwo/storage) 应该对其发布策略有类似的说明。

## 建议

因此，我提议 [存储](https://github.com/Xuanwo/storage) 应遵守以下政策：

### 语义版本

**[存储](https://github.com/Xuanwo/storage) 按跟随 [语义版本](https://semver.org/)**

- `1.0.0` 发布后， [存储](https://github.com/Xuanwo/storage) 将严格遵循语义版本
- All exported items in [storage](https://github.com/Xuanwo/storage) except `internal` and `tests` will be included in semantic versioning

[存储](https://github.com/Xuanwo/storage) 使用 [依赖bot](https://dependabot.com/) 自动升级其依赖关系，升级过程将如下所示：

- [依赖机器人](https://dependabot.com/) 将创建新的 PRs 到分支 [依赖性](https://github.com/Xuanwo/storage/tree/dependence) 并在构建成功后合并
- [存储](https://github.com/Xuanwo/storage) 将合并分支 [依赖关系到每个版本的 PR](https://github.com/Xuanwo/storage/tree/dependence)
- 发行版PR 合并后， [分支依赖](https://github.com/Xuanwo/storage/tree/dependence) 应重置为分支大师

### Target Golang Versions

**[存储](https://github.com/Xuanwo/storage) SHOULD 与最后两个黄金主要版本兼容**

假设当前版本为1.14

- 开发人员开发去1.13 **或** 去1.14
- CI 应该传递到 **BOTH** 去1.13 **and** 去1.14
- 任何错误/错误报告到 1.12MAY 标记为 `想修复`
- 不包含新功能 1.14

## 理由

无

## 兼容性

无

## 二． 执行情况

没有与代码相关的更改。
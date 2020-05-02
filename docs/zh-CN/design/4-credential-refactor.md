---
author: Xuanwo <github@xuanwo.io>
status: 完成
updated_at: 2019-12-26
updates:
  - design/3-support-service-init-via-config-string.md
---

# 提议：重新考虑全权证书

## 二. 背景

凭据是连接存储服务的大部分导入部分。

当前实现使用 `凭据。提供器` 输出值的接口：

```go
类型提供器接口 然后
    Value() 值
}

类型值结构变化。
    AccessKey 字符串
    SecretKey 字符串
}
```

我们只为提供商执行 `静态` 协议。

这个实现灵敏度来自爪子：不同的信任度是检索访问密钥和秘密密钥的不同方式。 然而，这种想法并不适用于每个服务。 对于谷歌：获取oauth2令牌的方式不同的可信度。

作为一个统一的储存层，我们需要设法消除这些不一致之处。

## 建议

重新因素 `凭据。提供` 到以下：

```go
输入提供商结构
    
 协议字符串
    args []string
}

func (p *Provider) Protocol() string {
    return p.protocol
}

func (p *Provider) Value() []string {
    return p.args
}
```

用户在使用证书之前需要检查证书的协议，凭据详细信息将存储在 `p. rgs` 可以通过 `获取。Value()`.

只支持一项凭据协议的服务可以是：

```go
credprotocacy, cred := opt.全权证书。Protocol(), opt.全权证书。Value()
如果信用协议 != 凭据。ProtocolHmac 然后
    return nil, fmt.Errorf(错误消息，s, 凭据)。错误不支持的协议)
}
// 与凭据值相关的服务配置。
cfg, err := config.新 (red[0], cred[1])
```

支持多个协议的服务可以是：

```go
cfg := aws。NewConfig()

credProtocol, cred := opt.全权证书。Protocol(), opt.全权证书。Value()
切换信用协议。
案例凭据。ProtocolHmac:
    cfg = cfg.全权证书审查报告。新 StaticCredentials(red[0], cred[1], ""))
case credential.ProtocolEnv:
    cfg = cfg.全权证书审查报告。NewEnvCredentials())
默认值：
    返回 nil, fmt。Errorf(错误消息，s, 凭据)。ErrUnsupportedProtocol)
}
```

为了执行一项新的议定书，应按照事物进行发展。

- 在格式 `协议<Name>` 中添加组合。 `ProtocolEnv`
- 为 `协议<Name>`添加注释，明确描述如何使用值
- 实现init 函数 `新<Name>(值 ...字符串) (*提供者，错误)` 和 `新建<Name>(值 ...字符串) (*提供者，错误)`
- 添加 `协议<Name>` 到 `凭据。解析` 切换大小写
- 添加单元测试实例

## 理由

AWS 使用提供商获取静态、 文件和环境中的密钥和密钥。

```go
键入凭据结构
    credit value
    force刷新布尔

    m sync。 WMutex

    providers Provider
}

type Value struct 哇，
    // AWS Access key ID
    AccessKeyID string

    // AWS Secret Access Key
    SecretAccessKey string

    // AWS Session Token
    SessionToken string

    // Provider used to get 凭据
    Provider name string
}

type Provider interface
    // Retrieve retrieve 返回 nil 如果它成功检索到了值的话。
    // 如果该值无法获取或为空则返回错误。
    检索() (值, 错误)

    // 离线返回，如果凭据不再有效，则需要
    // 才能检索。
    IsExpired() bool
}
```

谷歌有相同的提供商设计，但对于oauth2令牌：

```go
输入凭据结构
    ProjectID字符串// 可能是空的
    TokenSource oauth2。令牌源

    // JSON 包含一个 JSON 凭据文件的原始字节。
    // 如果认证是由
    // 环境提供的，而不是凭据文件，此字段可能为零。 当代码是
    // 在 Google 云平台上运行。
    JSON []byte
}
```

## 兼容性

只对 `凭据` 包进行更改，凭据的配置字符串将受到影响，公共接口中没有其他更改。

## 二． 执行情况

大多数工作将由本提案的作者完成。

为了避免协议 `使用Hmac` 和 `api 键`之间的误解(它们都是静态键)， 我们将当前的 `静态` 协议重命名为 `hmac` 没有这样不准确，但更加清楚。

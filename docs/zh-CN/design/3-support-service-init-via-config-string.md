---
author: Xuanwo <github@xuanwo.io>
status: 完成
updated_at: 2019-12-23
updated_by:
  - design/4-credital-refactor.md
  - '--remove-storager-init.md'
deprecated_by:
  - design/13-remove-config-string.md
---

# 建议：通过配置字符串进入支持服务

## 二. 背景

此项目打算成为Golang统一的存储层，但不同的存储层配置如此不同，我们无法很好地统一它们。

对于定位：我们只需要指定工作岗位。 对于对象存储：我们需要指定主机、 端口、 协议、 访问密钥ID和其他。

我们曾经通过类型和选项支持他们，就像我们在 [qscamel](https://github.com/qingstor/qscamel) 中那样：

每个服务 (qscamel中的端点) 都应该处理他们自己的选项：

```go
type Client struct {
    BucketName          string `yaml:"bucket_name"`
    Endpoint            string `yaml:"endpoint"`
    Region              string `yaml:"region"`
    AccessKeyID         string `yaml:"access_key_id"`
    SecretAccessKey     string `yaml:"secret_access_key"`
    DisableSSL          bool   `yaml:"disable_ssl"`
    UseAccelerate       bool   `yaml:"use_accelerate"`
    PathStyle           bool   `yaml:"path_style"`
    EnableListObjectsV2 bool   `yaml:"enable_list_objects_v2"`
    EnableSignatureV2   bool   `yaml:"enable_signature_v2"`
    DisableURICleaning  bool   `yaml:"disable_uri_cleaning"`

    Path string

    client *s3.S3
}

func New(ctx context)。Context, et uint8, hc *http.Client) (c *Client, err error) {
    ...
    content, err := yaml.Marshal(e)选项
    如果err != nil {
        return
    }
    err = yaml。Unmarshal(contents, c)
    如果是err != nil {
        return
    }
...
}
```

想要使用此服务的开发者应该处理此类型：

```go
切换 t.Src.Type {
...
大小写常数。端点S3：
    src, err = s3。新建(tx, 常量)。SourceEndpoint, context客户端)
    如果是err != nil {
        return
    }

default:
    logrus.错误("类型 %s 不支持。", t.Src.类型)
    err = 常量错误端点支持
    返回
}
```

用户应该在配置中直接设置它们：

```yaml
source:
  type: s3
  path: "/path/to/source"
  options:
    bucket_name: example_bucket
    endpoint: example_endpoint
    region: example_region
    access_key_id: example_access_key_id
    secret_access_key: example_secret_access_key
    disable_ssl: false
    use_accelerate: false
    path_style: false
    enable_list_objects_v2: false
    enable_signature_v2: false
    disable_uri_cleaning: false
```

它行之有效，但没有达到我们的目标。 为了解决这个问题，我们拆分PR [服务的端点和凭据：将端点和凭据分割成不同的配对](https://github.com/Xuanwo/storage/pull/34)。 在这个PR中，我们可以输入对象服务，例如：

```go
srv := qingstor。新()
err = srv.Init(
    pairs.否决(全权证书)。NewStatic(accessKey, secretKey)),
    配对取出点(端点)。NewStaticFromParsedURL(协议, 主机, 端口)，

如果是err != nil 然后
    log。Printf("service init 失败: %v", 错误)
}
```

这是更好的，但还不够。 我们需要一种通用的方式来提供所有服务，例如：

```go
srv := 存储。有点(有点)
```

## 建议

因此，我提议作如下修改：

### 引入“配置字符串”的概念

`配置字符串` 被广泛用于数据库连接：

mysql: `user:password@/dbname?charset=utf8&parseTime=True&loc=Local` postgres: `host=myhost port=myport user=gorm dbname=gorm password=mypassword` sqlserver: `sqlserver://username:password@localhost:1433?database=dbname`

像我们在 URL 中所做的那样，我们可以在一个格式化的字符串中使用不同的部分来表示不同的含义。

存储中的配置字符串类似于：

```
<type>://<config>
             +
             |
             v
<credential>@<endpoint>/<namespace>?<options>
     + + + +
     | +---------+ +------------------------------+
     v v
<protocol><data>   <protocol>:<data>         <key>:<value>[&<key>:<value>]
```

- 凭据： `<protocol>:<data>`, 数据内容由不同的凭据协议决定，静态凭据可以是 `静态凭据：<access_key>:<secret_key>`。
- 端点： `<protocol>:<data>`, 数据内容由不同的端点协议决定，qingstors的有效端点可以是 `https://:qingstor.com:443`。
- 命名空间：命名空间由物体存储的不同存储类型决定 它可能是 `<bucket_name>/<prefix>`, 对于位置 `<path>`
- 选项：多个 `<key>=<value>` 关联到 `&`

因此一个有效的配置字符串可以是：

- `qingstor://static:<access_key_id>:<secret_access_key>@https://:qingstor.com:443/<bucket_name>/<prefix>?zone=pek3b`
- `fs:///<work_dir>`

### 通过类型和配置字符串实现支持嵌入的函数

通过配置字符串的定义，我们可以实现更多的一般服务初始化功能。

我们将添加以下代码段中的更改：

- 在 `coreutils` 软件包中添加 `Open(配置字符串) (Servicer, Storager, 错误)` 函数。 `OpenServicer` and `OpenStorager` 将被添加以便更方便。
- 在 `pkg` 中添加 `配置` 包以进行配置字符串解析。
- 实现 `<service>。新(配对...*配对) (Servicer, error)` 函数，如果服务不能实现服务，请执行 `<service>新(对...*配对) (Storager, 错误)` 代替。

### 从服务界面移除Init

在配置字符串的全新支持下，我们可以移除服务界面中的Init。

## 兼容性

存储init 逻辑将被完全重新设置。

## 二． 执行情况

大多数工作将由本提案的作者完成。
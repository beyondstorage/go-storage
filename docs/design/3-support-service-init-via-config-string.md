---
author: Xuanwo <github@xuanwo.io>
status: candidate
updated_at: 2019-12-23
---

# Proposal: Support service init via config string

## Background

This project intents to be a unified storage layer for Golang, but different storage layers' configuration are so different we can't unify them well.

For posixfs: we only need to specify the workdir.
For object storage: we need to specify host, port, protocol, access key id and others.

We used to support them by type and options, like we did in [qscamel](https://github.com/qingstor/qscamel):

Every service(endpoint in qscamel) should handle their own options:

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

func New(ctx context.Context, et uint8, hc *http.Client) (c *Client, err error) {
	...
	content, err := yaml.Marshal(e.Options)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(content, c)
	if err != nil {
		return
	}
    ...
}
```

Developer who want to use this service should handle the type:

```go
switch t.Src.Type {
...
case constants.EndpointS3:
    src, err = s3.New(ctx, constants.SourceEndpoint, contexts.Client)
    if err != nil {
        return
    }
...
default:
    logrus.Errorf("Type %s is not supported.", t.Src.Type)
    err = constants.ErrEndpointNotSupported
    return
}
```

User should set them in config directly:

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

It works, but it doesn't meet our goal. To address this problem, we split endpoint and credential in PR [services: Split endpoint and credential into different pair](https://github.com/Xuanwo/storage/pull/34). In this PR, we can init an object service like:

```go
srv := qingstor.New()
err = srv.Init(
    pairs.WithCredential(credential.NewStatic(accessKey, secretKey)),
    pairs.WithEndpoint(endpoint.NewStaticFromParsedURL(protocol, host, port)),
)
if err != nil {
    log.Printf("service init failed: %v", err)
}
```

It's better, but not enough. We need a general way to init all service like:

```go
srv := storage.SomeCall(something)
```

## Proposal

So I propose following changes:

### Introduce the concept of "config string"

`config string` is widely used in db connections:

mysql: `user:password@/dbname?charset=utf8&parseTime=True&loc=Local`
postgres: `host=myhost port=myport user=gorm dbname=gorm password=mypassword`
sqlserver: `sqlserver://username:password@localhost:1433?database=dbname`

Like we did in URL, we can use different part in a formatted string to represent different meaning.

Config string in storage would be like:

```
<type>://<config>
             +
             |
             v
<credential>@<endpoint>/<namespace>?<options>
     +            +                 +
     |            +---------+       +----------------------+
     v                      v                              v
<protocol>:<data>   <protocol>:<data>         <key>:<value>[&<key>:<value>]
```

- credential: `<protocol>:<data>`, data's content decided by different credential protocol,static credential could be `static:<access_key>:<secret_key>`.
- endpoint: `<protocol>:<data>`, data's content decided by different endpoint protocol, qingstor's valid endpoint could be `https:qingstor.com:443`.
- namespace: namespace is decided by different storage type, for object storage, it could be `<bucket_name>/<prefix>`, for posixfs, it could be `<path>`
- options: multiple `<key>=<value>` connected with `&`

So a valid config string could be:

- `qingstor://static:<access_key_id>:<secret_access_key>@https:qingstor.com:443/<bucket_name>/<prefix>?zone=pek3b`
- `fs:///<work_dir>`

### Implement functions to support init via type and Config string

With de definition of Config string, we can implement functions for more general service initiation.

We will add following changes in codebase:

- Add `Open(config string) (Servicer, Storager, error)` function in `coreutils` package. `OpenServicer` and `OpenStorager` will be added for more convenient.
- Add `config` package in `pkg` to do config string parse.
- Implement `<service>.New(pairs ...*Pair) (Servicer, error)` function, if service doesn't implement Servicer, implement `<service>.New(pairs ...*Pair) (Storager, error)` instead.

### Remove Init from Servicer interface

With the brand-new support of config string, we can remove Init in Servicer interface.

## Compatibility

Storage init logic will be totally refactored.

## Implementation

Most of the work would be done by the author of this proposal.
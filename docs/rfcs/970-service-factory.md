- Author: Xuanwo <github@xuanwo.io>
- Start Date: 2021-10-27
- RFC PR: [beyondstorage/go-storage#970](https://github.com/beyondstorage/go-storage/issues/970)
- Tracking Issue: [beyondstorage/go-storage#0](https://github.com/beyondstorage/go-storage/issues/0)

# GSP-970: Service Factory

## Background

The way to init a service or storage is the most important part of the project.

- In [GSP-13: Remove config string](./13-remove-config-string.md), we add pairs support.
- In [GSP-48: Service Registry](./48-service-registry.md), we implement a service registry.

For now, we can init a service in this way:

```go
package main

import (
   "log"

   "go.beyondstorage.io/v5/services"
   "go.beyondstorage.io/v5/types"

   _ "go.beyondstorage.io/services/s3/v3"
)

func main() {
   // Init a Storager from connection string. 
   store, err := services.NewStoragerFromString("s3://bucket_name/path/to/workdir")
   if err != nil {
      log.Fatalf("service init failed: %v", err)
   }
}
```

However, only supporting config strings is not enough.

We need to support other configs. For example, we should support init from a map like `map[string]interface{}`, so developers can marshal and unmarshal them from a config file. Even better, we should support init from a struct, so developers can have strong type support.

Besides, the current config string implementation is hard to maintain. We have to maintain a whole pair map between go-storage and services:

```go
var pairMap = map[string]string{"content_md5": "string", "content_type": "string", "context": "context.Context", "continuation_token": "string", "credential": "string", "default_content_type": "string", "default_io_callback": "func([]byte)", "default_service_pairs": "DefaultServicePairs", "default_storage_class": "string", "default_storage_pairs": "DefaultStoragePairs", "disable_100_continue": "bool", "enable_virtual_dir": "bool", "enable_virtual_link": "bool", "endpoint": "string", "excepted_bucket_owner": "string", "expire": "time.Duration", "force_path_style": "bool", "http_client_options": "*httpclient.Options", "interceptor": "Interceptor", "io_callback": "func([]byte)", "list_mode": "ListMode", "location": "string", "multipart_id": "string", "name": "string", "object_mode": "ObjectMode", "offset": "int64", "server_side_encryption": "string", "server_side_encryption_aws_kms_key_id": "string", "server_side_encryption_bucket_key_enabled": "bool", "server_side_encryption_context": "string", "server_side_encryption_customer_algorithm": "string", "server_side_encryption_customer_key": "[]byte", "service_features": "ServiceFeatures", "size": "int64", "storage_class": "string", "storage_features": "StorageFeatures", "use_accelerate": "bool", "use_arn_region": "bool", "work_dir": "string"}
```

As time goes, the map will become bigger.

How about moving the parsing logic to a factory and implementing them on the service side?

## Proposal

So I propose to implement a factory to init a service.

On `go-storage` side, we will add a new interface called `Factory`:

```go
type Factory interface {
    FromString(conn string) (err error)
    FromMap(m map[string]interface{}) (err error)
    WithPairs(ps ...types.Pair) (err error)
    
    NewServicer() (srv types.Servicer, err error)
    NewStorager() (sto types.Storager, err error)
}
```

And all existing functions will rewrite into `Factory` calls:

```go
func NewStoragerFromString(conn string, ps ...types.Pair) (types.Storager, error) {
    f, err := NewFactoryFromString(conn, ps...)
    if err != nil {
        return nil, err
    }
    
    return f.NewStorager()
}
```

On service side, we will:

- Move `namespace.service.new` to `factory.service`
- Move `namespace.storage.new` to `factory.storage`

And use the while `factory` field to generate the `Factory` struct:

```go
type Factory struct {
    // Service pairs.
    // Service required pairs.
    Credential string
    Endpoint   string
    // Service optional pairs.
    ForcePathStyle bool
    UseAccelerate  bool
    UseArnRegion   bool
    // Storage pairs.
    // Storage required pairs.
    Location string
    Name     string
    // Storage optional pairs.
    WorkDir             string
    DefaultContentType  string
    DefaultIoCallback   func ([]byte)
    DefaultStorageClass string
    EnableVirtualDir    bool
    EnableVirtualLink   bool
}
```

This `Factory` struct will implement the go-storage's `Factory` interface. And all existing functions will rewrite into `Factory` calls.

With this struct, we can unify the initialization logic of service and storage together.

## Rationale

TBD

## Compatibility

No public API breaks, but all services should refactor the init logic.

## Implementation

- Add `Factory` interface in go-storage.
- Add `factory` field support.
- Mark `RegisterSchema`, `RegisterServicer`, `RegisterStorager` as deprecated.
- Migrate all services.
- Remove all deprecated functions.

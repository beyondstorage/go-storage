---
author: Xuanwo <github@xuanwo.io>
status: finished
updated_at: 2020-01-16
updates: [3](./3-support-service-init-via-config-string.md)
---

# Proposal: Remove Storager Init

## Background

In commit [storager: All API now use relative path instead](https://github.com/Xuanwo/storage/commit/1cb485ec1f64d59cff19414005f9f602b3721cef), we first added `Init` functions in `Storager` API.

We add this function to address `Storager` init problem: we need to configure the Storager. For example: we need to set `WorkDir` (used to be `Base`) for all Storager.

So our interface turns into:

```go
// Init will init storager itself.
//
// Caller:
//   - Init MUST be called after created.
Init(pairs ...*types.Pair) (err error)
// InitWithContext will init storager itself.
InitWithContext(ctx context.Context, pairs ...*types.Pair) (err error)
```

Our init logic is include three parts:

```go
// 1. Create a new Servicer
srv, err = azblob.New(opt...)
if err != nil {
    return
}
name, prefix := namespace.ParseObjectStorage(ns)
// 2. Get a Storager form Servicer
store, err = srv.Get(name)
if err != nil {
    return
}
// 3. Init Storager with pairs.
err = store.Init(pairs.WithWorkDir(prefix))
if err != nil {
    return
}
```

If there is no Servicer here, init logic will turn into two parts:

```go
store = fs.New()
path := namespace.ParseLocalFS(ns)
err = store.Init(pairs.WithWorkDir(path))
if err != nil {
    return
}
```

It looks like we solve the init problem for Storager, but not really.  There are following problems.

- `Init` is only used in `coreutils`

If an API only used in internal packages, why should we export it?

- Only `WorkDir` supported, and hard to add more pairs

For now, we hardcoded `store.Init(pairs.WithWorkDir(path))` in `coreutils`. Firstly, how to add more pairs? Then, it's expensive to do this in `Init` if only `WorkDir` needed.

- All `Storager` implement the same interface with the same way.

If implementation is the same, we don't need to export it as an interface. 

- `Init` can be called times and may cause concurrent problems

Users can change `WorkDir` during `List`, which is not allowed.

## Proposal

So I propose following changes:

- Merge `Init` and `newStorage` to `Storager.init(pairs ...*types.Pair)`
- Refactor `New` to `New(pairs ...*types.Pair) (srv *Service, store *Storage, err error)`
- Add `Init` related pairs to `Servicer`'s `Get` and `New`
- Refactor config string option handles

## Rationale

### Refactor config string option handles

Config string used to be like:

`qingstor://hmac:<access_key>:<secret_key>@<protocol>:<host>:<port>/<bucket_name>/<prefix>`

The original design expect to reach a balance between user input and config string parse. Let's take following two styles into consideration:

- `fs://?work_dir=/path/to/dir` and `azblob://hmac:<access_key>:<secret_key>?name=<bucket_name>&work_dir=<prefix>`
- `fs:///path/to/dir` and `azblob://hmac:<access_key>:<secret_key>/<bucket_name>/<prefix>`

It's obvious the 2nd style has fewer user input, proposal [3-support-service-init-via-config-string](./3-support-service-init-via-config-string.md) has the same idea. However, this style introduces problem between different kinds of services.

- `fs:///path/to/dir`
- `azblob://hmac:<access_key>:<secret_key>`
- `azblob://hmac:<access_key>:<secret_key>/<bucket_name>`
- `azblob://hmac:<access_key>:<secret_key>/<bucket_name>/<prefix>`

We need service type in order to distinguish `/path/to/dir` and `/<bucket_name>/<prefix>`, this makes config string parser hard to implement, and we have to add concept like `namespace` so that we can defer this work in `coreutils`.

Back to the first style, how about delegating this work to end users?

- `fs://?work_dir=/path/to/dir`
- `azblob://hmac:<access_key>:<secret_key>?name=<bucket_name>&work_dir=<prefix>`

Now, both `name` and `work_dir` can be parsed without service type.

## Compatibility

Users who call services `New` or `Init` directly will facing breaking changes.

## Implementation

Most of the work would be done by the author of this proposal.
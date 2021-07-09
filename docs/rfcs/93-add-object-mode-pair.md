- Author: Xuanwo <github@xuanwo.io>
- Start Date: 2021-06-08
- RFC PR: [beyondstorage/specs#93](https://github.com/beyondstorage/specs/issues/93)
- Tracking Issue: N/A

# GSP-93: Add ObjectMode Pair

## Background

We have introduced `ObjectMode` in [GSP-25]:

```go
type ObjectMode uint32

const (
ModeIrregular ObjectMode = 0

ModeDir ObjectMode = 1 << iota
ModeRead
ModeLink
)
```

But it's not enough, we need `ObjectMode` pair too. For example, the `Create` operation needs `ObjectMode` to create a new object. And `ObjectMode` COULD be used as input restriction for operations. Let me declare the second case more clearly.

[GSP-49] added create dir operation, and [GSP-87] introduced features gates. With those GSP, our user can use `CreateDir` to create a dir object. In services like s3 which doesn't have native `CreateDir` support, we usually simulated it via creating an object that ends with `/`. But the behavior doesn't work in `Stat` and `Delete`.

- If user call `Stat("test")` after `CreateDir("test")`, he will get `ObjectNotExist` error.
- If user call `Delete("test")` after `CreateDir("test")`, no object will be removed and `test/` will be kept in service.

We need to find a way to figure it out.

## Proposal

So I propose to add the `ObjectMode` pair. This pair COULD be used in the following operations:

- `Create`: set the output object's `ObjectMode`
- `Stat`: ObjectMode hint, returns error if `ObjectMode` not meet.
- `Delete`: ObjectMode hint

For `Stat` and `Delete`

- Service SHOULD use the `ObjectMode` pair as a hint.
- Service could have different implementations for different `ObjectMode`.

Take `s3` as an example, we simulate `CreateDir` via creating object ends with `/`. `CreateDir("test")` will create an object `test/` in s3. And we can

- stat this object via `Stat("test", pairs.WithObjectMode(types.ObjectModeDir))`
- delete this object via `Delete("test", pairs.WithObjectMode(types.ObjectModeDir))`

## Rationale

### Alternative design for Dir Object

Except `ObjectMode`, we have following alternative designs:

- `CreateDir("test")` + `Stat("test/")`/`Delete("test/")`
- `CreateDir("test/")` + `Stat("test/")`/`Delete("test/")`

They have the following drawbacks:

- Not consistent with other APIs
- Hard to use: user need to do extra work to change the object path

## Compatibility

All changes are compatible.

## Implementation

- Add `ObjectMode` pairs
- Add `ObjectMode` into operations (all service must support)

[GSP-25]: ./25-object-mode.md
[GSP-49]: ./49-add-create-dir-operation.md
[GSP-87]: ./87-feature-gates.md

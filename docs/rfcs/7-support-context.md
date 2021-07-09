- Author: Xuanwo <github@xuanwo.io>
- Start Date: 2020-01-08
- RFC PR: N/A
- Tracking Issue: N/A

# Proposal: Support context

## Background

context has been widely used in Golang, more and more project rely on context to carry deadline, cancellation or other values, including gRPC, opentracing, GCS SDK and so on.

Add context support for [storage](https://github.com/Xuanwo/storage) will make it a production ready storage layer for real world:

- Allow set deadline for operation
- Support request cancellation
- Support tracing
- ...

## Proposal

So I propose following changes:

- Add context pair support for every public API
- Add `ReadWithContext` style method for every public API
- `ReadWithContext` call `Read` with `Context` pair
- `Read` use `Context` via parsed pair

More detailed changes described as following:

### Add context pair support for every public API

We treat `context` as pre-defined pairs for API, and make sure `context` is provided in generated code:

```go
v, ok = values[pairs.Context]
if ok {
    result.Context = v.(context.Context)
} else {
    result.Context = context.Background()
}
```

If there is `context` in pairs, we will use the `context`, or we will create a new one.

This section will change `internal/cmd/meta`.

### Add `ReadWithContext` style method for every public API

Add `XxxWithContext` API for every method, for example:

```go
type Mover interface {
	Move(src, dst string, pairs ...*types.Pair) (err error)
}
```

will turn into:

```go
type Mover interface {
	Move(src, dst string, pairs ...*types.Pair) (err error)
	MoveWithContext(ctx context.Context, src, dst string, pairs ...*types.Pair) (err error)
}
```

This operation will be executed by hand. We don't have too many interfaces here, so no bother to write a tool.

### Generate XxxWithContext API for services

Let's generate code to archive this.

First of all, we need to add `...*types.Pair` for every API to carry context, this will affect two APIs:

- `Metadata`
- `Statistical`

Not needed, but also make sense: they both could call API.

Then, we need to generate `XxxWithContext` API, so that implementers don't need to care about that.

In generated code, we will do following things:

```go
func (s *Storage) ReadWithContext(ctx context.Context, path string, pairs ...*types.Pair) (r io.ReadCloser, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "github.com/Xuanwo/storage/services/qingstor.Storager.Read")
	defer span.Finish()

	pairs = append(pairs, ps.WithContext(ctx))
	return s.Read(path, pairs...)
}
```

### Services should handle context

Service should use context from parsed pairs:

```go
it := s.bucket.Objects(context.TODO(), &gs.Query{
    Prefix: rp,
})
```

should be updated to:

```go
it := s.bucket.Objects(opt.Context, &gs.Query{
    Prefix: rp,
})
```

## Rationale

Other implement could be following:

- Add ctx in API, and add `XXXWithoutContext` API support
- Add ctx in API, and let users who don't care context to use `context.TODO()`
- Add context pair support and don't touch public interface

### Add ctx in API, and add `XXXWithoutContext` API support

Interface will be like following:

```go
type Mover interface {
	Move(ctx context.Context, src, dst string, pairs ...*types.Pair) (err error)
	MoveWithoutContext(src, dst string, pairs ...*types.Pair) (err error)
}
```

First of all, this change is a break change, every API call need to be refactored.

Then, it's obvious that `MoveWithoutContext` is longer than `MoveWithContext`.

The most important thing is the thought behind API design: **Fair**.

There are two kinds of developers here: some of them need context support, others don't care about it.

Design in the proposal is friendly for both of them:

- people who need context support should use `XxxWithContext` or add context pair in `Xxx` call, they know what they are doing.
- people who don't need context support can write code happily without any idea about context.

However, this design is not fair for people who don't need context support. They need to use API like `MoveWithoutContext` although they don't care about context.

### Add ctx in API, and let users who don't care context to use `context.TODO()`

Similar reason as described in the previous chapter.

### Add context pair support and don't touch public interface

Looks fine, but a bit inconvenient. This design makes it hard to add tracing support. People need to wrap code like:

```go
func ReadWithContext(ctx context.Context, s *Storage, path string, pairs ...*types.Pair) (r io.ReadCloser, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Read")
	defer span.Finish()

	pairs = append(pairs, ps.WithContext(ctx))
	return s.Read(path, pairs...)
}
```

Why not let us do it ourselves?

## Compatibility

No break changes

## Implementation

Most of the work would be done by the author of this proposal.

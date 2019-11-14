---
author: Xuanwo <github@xuanwo.io>
status: draft
updated_at: 2019-11-13
---

# Proposal: Unify storager behavior

## Background

We provide a `Capable` function for developers to check whether underlying storager support action/pair or not. But in reality, no one use them. As a `unified storage layer`, it's strange that coder still need to take care of the difference of storager services.

But differences indeed exist. Either we give up some power of storager services, or provide better way to developers to handle the complexity. We need to eliminate the inconsistent behavior for base actions, and provide convenient ways for user to use high-level actions.

Except for better consistence, current API design also lead to some wired results. One of the case looks like following:

```go
it := store.ListDir(path, types.WithRecursive(true)

for {
    o, err := it.Next()
    if err != nil && errors.Is(err, iterator.ErrDone) {
        break
    }
    if err != nil {
        t.TriggerFault(types.NewErrUnhandled(err))
        return
    }
    store.Delete(o.Name)
}
```

In order to remove all Object under a path, it's obviously that we need to list them recursively and delete them one by one. This algorithm works well for object storage, a.k.a. prefix based storage. But it will fail on POSIX file systems, because our `List` doesn't return folders.

## Proposal

So I propose following changes:

### Split Base Storage

Split base storager for currently `Storager` interface, and make sure every storage will have same operations for same action and pairs. If storage don't support specific pairs, them need to ignore them.

The base `Storager` could be:

```go
type Storager interface {
	String() string

	Init(pairs ...*types.Pair) (err error)
	Capable(action string, key ...string) bool
	Metadata() (types.Metadata, error)

	ListDir(path string, pairs ...*types.Pair) (err error)
	Read(path string, pairs ...*types.Pair) (r io.ReadCloser, err error)
	Write(path string, r io.Reader, pairs ...*types.Pair) (err error)
	Stat(path string, pairs ...*types.Pair) (o *types.Object, err error)
	Delete(path string, pairs ...*types.Pair) (err error)
}
```

### Split high level funcs

High level funcs like `Copy`, `Move` should be split aside. Caller should de type assert before use them:

```go
type Copier interface {
	Copy(src, dst string, pairs ...*types.Pair) (err error)
}

if x, ok := store.(Copier); ok {
    err := x.Copy(oldpath, newpath)
    if err != nil {
        return err
    }
}
```

There will be a `Segmenter` interface for all segments related operations.

```go
type Segmenter interface {
	ListSegments(path string, pairs ...*types.Pair) (err error)
	InitSegment(path string, pairs ...*types.Pair) (id string, err error)
	WriteSegment(id string, offset, size int64, r io.Reader, pairs ...*types.Pair) (err error)
	CompleteSegment(id string, pairs ...*types.Pair) (err error)
	AbortSegment(id string, pairs ...*types.Pair) (err error)
}
```

### Add callback func for all list operations

List operations will not return `iterator.ObjectIterator` anymore, instead, it will allow input `WithDirFunc` and `WithFileFunc` via `Pair`.

```go
dirFunc := func(object *types.Object) {
    printf("dir %s", object.Name)
}
fileFunc := func(object *types.Object) {
    printf("file %s", object.Name)
}

err := store.List("prefix", types.WithDirFunc(dirFunc), types.WithFileFunc(fileFunc))
if err != nil {
    return err
}
```

### Remove recursive for List

Based on the work of `callback in list`, we can remove all recursive support in list.

Directory based storage will only list one directory, and prefix based storage will only list one prefix without a delimiter.

## Rationale

### Why Copier interface

In order to know whether a Storager support a specific API call or not, we need some way to transmit this information. As far as I know, there are following ways to archive it:

- Copier interface
- Copyable func call
- Copy Capability
- Panic/Not Supported Error

We will discuss them one by one and do a benchmark on them.

For Copier interface

```go
type Copier interface {
	Copy(src, dst string, pairs ...*types.Pair) (err error)
}

if x, ok := store.(Copier); ok {
    err := x.Copy(oldpath, newpath)
    if err != nil {
        return err
    }
}
```

We create different interfaces for different capability, and Caller need to do type assert for them.

For Copyable func call

```go
type Storage interface {    
	Copy(src, dst string, pairs ...*types.Pair) (err error)
    Copyable() bool
}

if store.Copyable() {
    err := x.Copy(oldpath, newpath)
    if err != nil {
        return err
    }
}
```

Add different `XXXable` func call which return bool for different capability, and Caller need to check it's return value before use them.

For Copy Capability

```go
type Storage interface {    
	Capability() Capability
}

if store.Capability() & types.CapabilityCopy == 1 {
    err := x.Copy(oldpath, newpath)
    if err != nil {
        return err
    }
}
```

`Storager` supports return an uint64 or something to represent Capability, and Caller need to check this value before use them.

Panic/Not Supported Error

```go
func(store *Storager) Copy() {
    panic("not supported")
}
```

`Storager` will panic or return error for not supported funcs, Caller need to do `recover` or error check after use them.

Benchmark file could be found [here](./1/storager_interface_test.go), and the results looks like:

```go
goos: linux
goarch: amd64
pkg: github.com/Xuanwo/storage/docs/design/1
BenchmarkCopierInterface-8    	144869151	         8.48 ns/op
BenchmarkCopyableFuncCall-8   	397990764	         3.03 ns/op
BenchmarkCopyCapability-8     	503389296	         2.26 ns/op
BenchmarkError-8              	49897329	        25.1 ns/op
BenchmarkPanic-8              	17014038	        71.4 ns/op
PASS
```

It's obvious that `Capability` is the more fast, and panic with recover is the slowest.

Performance is on the one hand, we also need to considerate whether the API is easy to use or easy to implement.

Func call is easy to use but developers should add another two funcs into their struct. Capability is easy to implement, but add another concept for Caller to understand. Interface is far more clear, but looks a bit trick in real usage.

The more import thing for Interface usage is that we can enforce Caller to check Storager support before really use them.

So we pick Interface for it's centered performance and mandatory ability.

### Why callback

TODO

## Compatibility

### Copy, Move and Reach

Copy, Move and Reach will be removed from `Storager` interface, Caller should do type assert before use them.

### Segments Related Operations.

All segments related operations will be remove from `Storager` interface, Caller should do type assert of `Segmenter` before use them.

### No recursive in list

`ListDir` will not support recursive anymore, handle `recursive` in `DirFunc` instead.

## Implementation

Most of the work would be done by the author of this proposal.
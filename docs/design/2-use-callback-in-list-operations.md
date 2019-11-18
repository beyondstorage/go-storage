---
author: Xuanwo <github@xuanwo.io>
status: finished
updated_at: 2019-11-18
---

# Proposal: Use callback in List operations

Current API design leads to some wired results. One of the case looks like following:

```go
it := store.ListDir(path, types.WithRecursive(true))

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

In order to remove all Object under a path, it's obviously that we need to list them recursively and delete them one by one. This algorithm works well for object storage, a.k.a. prefix based storage. It will fail on POSIX file systems, because our `List` doesn't return folders.

## Proposal

So I propose following changes:

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

### Why not recursive in Delete

If we support `recursive` in Delete, delete a dir or delete a prefix maybe much simpler:

```go
store.Delete(path, types.WithRecursive(true))
```

In a prefix based storage service, to support delete recursive means we need to do list objects in `Delete`, which make our API not orthogonal, hard to implement and much hard to do unit test.

In order to keep same style across all APIs, we also need to add recursive support in `Copy`, `Move` and so on.

In addition to the reasons mentioned above, implement recursive in `Delete` makes our API can't be used in concurrent task frameworks: they can't be split into sub-tasks anymore.

### Why not return Dir in List

The other idea could be return both `Dir` and `File` in `List` instead of `File` only. This  method may import two problems.

Firstly, although we can get `Dir` from `List`, this doesn't solve anything. Caller need to make sure every `Object` under this `Dir` has been deleted. They need to code like these:

```go
it := store.ListDir(path, types.WithRecursive(true))
dirs := make([]*types.Object, 0)

for {
    o, err := it.Next()
    if err != nil && errors.Is(err, iterator.ErrDone) {
        break
    }
    if err != nil {
        t.TriggerFault(types.NewErrUnhandled(err))
        return
    }
    if o.Type == types.ObjectTypeDir {
        dirs = append(dir, o)
    } else {
        store.Delete(o.Name)
    }    
}

for i:=len(dirs)-1;i>=0;i-- {
    store.Delete(dirs[i].Name)
}
```

So ugly.

Secondly, it also make it harder for Storager to implement them and test them which doesn't our goal.

## Compatibility

No recursive in List operations anymore, instead caller need to input a func pair.

## Implementation

Most of the work would be done by the author of this proposal.

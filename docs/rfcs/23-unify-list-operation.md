- Author: Xuanwo <github@xuanwo.io>
- Start Date: 2021-01-11
- RFC PR: N/A
- Tracking Issue: N/A

# Proposal: Unify List Operation

## Background

[go-storage]'s list operation has been changed many times:

[Proposal: Unify storager behavior](1-unify-storager-behavior.md) use style like:

```go
ListDir(path string, pairs ...*types.Pair) iterator.ObjectIterator
```

We use a `Recursive` Pair to represent list a dir or list a prefix. 

This design could lead to wired results as [Proposal: Use callback in List operations](2-use-callback-in-list-operations.md) described. So in that preposal, we introduced a callback style list operations:

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

In this design: "Directory based storage will only list one directory, and prefix based storage will only list one prefix without a delimiter."

However, we do need dir support even under prefix based storage. So [Proposal: Support both directory and prefix based list](12-support-both-directory-and-prefix-based-list.md) add a new func called: `ObjectFunc`, and implement with following rules:

- Add `ObjectFunc` for prefix based list
- Treat `FileFunc` and `DirFunc` as directory based list

Obviously, it's hard to maintain and also hard to use. [Proposal: Split storage list](19-split-storage-list.md) called it *a failure by practice*. In the new proposal, we introduce new interfaces to do list operations:

- Split `List` into `ListDir` and `ListPrefix`
- Remove `List` from `Storager`
- Add interface `DirLister` for `ListDir`
- Add interface `PrefixLister` for `ListPrefix`

And with segment support:

- Rename `ListSegments` to `ListPrefixSegments` to match prefix changes
- Remove `ListSegments` from `Segmenter`
- Add interface `PrefixSegmentsLister` for `ListSegments`

In following implementations, we removed object related funcs and come back to object iterator:

```go
ListPrefix(prefix string, pairs ...Pair) (oi *ObjectIterator, err error)
```

In current design, our implementations look like following:

```go
func (s *Storage) listDir(ctx context.Context, dir string, opt *pairStorageListDir) (oi *typ.ObjectIterator, err error) {
	marker := ""
	delimiter := "/"
	limit := 200

	rp := s.getAbsPath(dir)

	input := &listObjectInput{
		Limit:     &limit,
		Marker:    &marker,
		Prefix:    &rp,
		Delimiter: &delimiter,
	}

	return typ.NewObjectIterator(ctx, s.listNextDir, input), nil
}

func (s *Storage) listNextDir(ctx context.Context, page *typ.ObjectPage) error {
	...
}

func (s *Storage) listPrefix(ctx context.Context, prefix string, opt *pairStorageListPrefix) (oi *typ.ObjectIterator, err error) {
	marker := ""
	limit := 200

	rp := s.getAbsPath(prefix)

	input := &listObjectInput{
		Limit:  &limit,
		Marker: &marker,
		Prefix: &rp,
	}

	return typ.NewObjectIterator(ctx, s.listNextPrefix, input), nil
}

func (s *Storage) listNextPrefix(ctx context.Context, page *typ.ObjectPage) error {
    ...
}
```

Maybe we can unify `listDir` and `listPrefix` in `list` here?

## Proposal

So I propose following changes:

- Merge `ListPrefix` and `ListDir` into `List`
- Remove `Lister` and add `List` into `Storager`
- Add a new pair `ListType` with two valid value: `dir` and `prefix`
  - Storager could have their default list type: filesystem could use `dir` and object storage could use `prefix` 

Segments related API should also be changed:

- Merge `ListPrefixSegments` and `ListDirSegments` into `ListSegments`

## Rationale

### How to address the DeleteAll problem?

DeleteAll problem is introduced in [Proposal: Use callback in List operations](2-use-callback-in-list-operations.md). The unified list operation can't delete a folder with content:

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

In this proposal, we change `Recursive` into `ListType` and require object type check after list because we have link object now. So we can return dirs even while listing with prefix list type.

## Compatibility

List related APIs could be changed.

## Implementation

Most of the work would be done by the author of this proposal.

[go-storage]: https://github.com/beyondstorage/go-storage

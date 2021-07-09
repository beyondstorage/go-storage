- Author: Xuanwo <github@xuanwo.io>
- Start Date: 2021-01-14
- RFC PR: N/A
- Tracking Issue: N/A

# Proposal: Object Mode

## Background

We used to use `ObjectType` to represent an Object:

```go
type ObjectType string

const (
	ObjectTypeFile ObjectType = "file"
	ObjectTypeDir ObjectType = "dir"
	ObjectTypeLink ObjectType = "link"
	ObjectTypeUnknown ObjectType = "unknown"
)
```

But when we try to extend this framework to more object types, we faced following problems.

### ObjectType is not orthogonal

Assuming that we have added `append` support, in order to mark an object do support append, we may introduce a new object type: `ObjectTypeAppend`. However, what will happen if we want to check whether an object supports reading as a normal file? Both 
`ObjectTypeFile` and `ObjectTypeAppend` support read as a normal file. This design make it hard to do operation towards specific objects.

### Check all ObjectType is tired

Take [noah](https://github.com/qingstor/noah) as an example:

```go
it := x.GetObjectIter()
for {
    obj, err := it.Next()
    if err != nil {
        if errors.Is(err, typ.IterateDone) {
            break
        }
        return types.NewErrUnhandled(err)
    }
    switch obj.Type {
    case typ.ObjectTypeFile:
       ...
    case typ.ObjectTypeDir:
        ...
    default:
        return types.NewErrObjectTypeInvalid(nil, obj)
    }
}
```

It's tired to check all object type here.

## Proposal

So I propose following changes

### Add Object Mode to replace Object Type

Add Object Mode like `os.FileMode` does:

```go
type ObjectMode uint32

const (
	ModeIrregular ObjectMode = 0

	ModeDir ObjectMode = 1 << iota
	ModeRead
	ModeLink
)
```

One mode means one or more operations we can do on an object: 

- `ModeDir` means we can do list on this object('s path)
- `ModeRead` means we can read it as a normal file
- `ModeLink` means we can use we can read this object's target

So we can compose them together:

- `ModeRead & ModeLink`: Think about a symlink, we can still read it.
- `ModeDir & ModeLink`: Think about a symlink to dir.
- `ModeDir & ModeRead`: Some broken object storage system could have an object end with `/`
- ...

Adding more mode will not affect current implementations.

### Merge Segment into Object

Now we have a much more flexible Object design, we can remove the segment abstraction:

- IndexSegment means an Object with `ModeMultipart` a.k.a. Multipart Object
- OffsetSegment means an Object with `ModePage` a.k.a. Page Object
- Object which supports append mains Object with `ModeAppend`.

After this change, `Object` is the only abstraction under `Storager`. So we can:

- Use `list` to list all kind of objects.
- Use `delete` to delete all kind of objects.
- ...

In this way, we can reduce operations like `ListIndexSegments` and `AbortSegment`.

So we can add interfaces in object style:

```go
// Append object support
type Appender interface {
	CreateAppend(path string, pairs ...Pair) (o *Object, err error)
	WriteAppend(o *Object, r io.Reader, size int64, pairs ...Pair) (n int64, err error)
}
// Random Write support via page object
type Pager interface {
	CreatePage(path string, pairs ...Pair) (o *Object, err error)
	WritePage(o *Object, r io.Reader, size int64, offset int64, pairs ...Pair) (n int64, err error)
}
// Multipart Upload support via multipart object.
type Multiparter interface {
	CompletePart(o *Object, parts []*Part, pairs ...Pair) (err error)
	CreatePart(path string, pairs ...Pair) (o *Object, err error)
	ListPart(o *Object, pairs ...Pair) (pi *PartIterator, err error)
	WritePart(o *Object, r io.Reader, size int64, index int, pairs ...Pair) (n int64, err error)
}
```

## Rationale

This designed is highly influenced by [Microsoft Azblob API](https://docs.microsoft.com/en-us/rest/api/storageservices/blob-service-rest-api)

## Compatibility

Object related operations could be changed.

## Implementation

Most of the work would be done by the author of this proposal.

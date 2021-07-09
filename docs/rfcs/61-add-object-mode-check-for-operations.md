- Author: Prnyself <lanceren@yunify.com>
- Start Date: 2021-05-11
- RFC PR: [beyondstorage/specs#61](https://github.com/beyondstorage/specs/issues/61)
- Tracking Issue: N/A

# AOS-61: Add object mode check for operations

## Background

In [AOS-25], we added support for object modes by bit map. All available object modes are listed below:

```go
type ObjectMode uint32

// All available object mode
const (
    // ModeDir means this Object represents a dir which can be used to list with dir mode.
    ModeDir ObjectMode = 1 << iota
    // ModeRead means this Object can be used to read content.
    ModeRead
    // ModeLink means this Object is a link which targets to another Object.
    ModeLink
    // ModePart means this Object is a Multipart Object which can be used for multipart operations.
    ModePart
    // ModeBlock means this Object is a Block Object which can be used for block operations.
    ModeBlock
    // ModePage means this Object is a Page Object which can be used for random write with offset.
    ModePage
    // ModeAppend means this Object is a Append Object which can be used for append.
    ModeAppend
)
```

It is intended to check object mode at the start of specific operation. For instance, both `WritePart`
and `WriteAppend` got a pointer to `Object` as an input, we need to ensure this `Object` is available 
for certain operation, so we should add object mode check and return `ObjectModeInvalidError`(introduced in [AOS-47]) asap if `Object` not fit.

### Current Practice

The check is implemented by service in the actual method call. 
For example, in [go-service-qingstor](https://github.com/beyondstorage/go-service-qingstor/blob/master/storage.go#L534):
```go
func (s *Storage) writeAppend(ctx context.Context, o *Object, r io.Reader, size int64, opt pairStorageWriteAppend) (n int64, err error) {
    if !o.Mode.IsAppend() {
        err = services.ObjectModeInvalidError{Expected: ModeAppend, Actual: o.Mode}
        return
    }
    ...
}
```

We check `*Object` mode `IsAppend` at the start of `writeAppend` which is called by `WriteAppendWithContext` which is generated.

```go
func (s *Storage) writeMultipart(ctx context.Context, o *Object, r io.Reader, size int64, index int, opt pairStorageWriteMultipart) (n int64, err error) {
    if o.Mode&ModePart == 0 {
        return 0, services.ObjectModeInvalidError{Expected: ModePart, Actual: o.Mode}
    }
    ...
}	
```

We check `*Object` mode `ModePart` at the start of `writeMultipart` which is called by `WriteMultipartWithContext` which is generated, too.

### Operation with Object Input

All function listed below contains `XXXWithContext` method 

- Appender
  - CommitAppend
  - WriteAppend
- Blocker
  - CombineBlock
  - ListBlock
  - WriteBlock
- Multiparter
  - CompleteMultipart
  - ListMultipart
  - WriteMultipart
- Pager
  - WritePage

## Proposal

So I propose that we should add mode check in specific operation, return `ObjectModeInvalidError` if mode not meet, and the check should be generated.

To generate object mode check in different functions, we will follow these steps (take `Multiparter.WriteMultipart` as an example):

1. Add a field named `object_mode` in `operations.toml`.

```toml
[multiparter.op.write_multipart]
description = "will write content to a multipart."
params = ["o", "r", "size", "index"]
results = ["n", "part"]
object_mode = "part"
```

2. Add `ObjectMode` field parsing in `parse.go` and `spec.go`.

```go
// tomlOperation in parse.go
type tomlOperation struct {
    Description string   `toml:"description"`
    Params      []string `toml:"params"`
    Pairs       []string `toml:"pairs"`
    Results     []string `toml:"results"`
    ObjectMode  string   `toml:"object_mode"` // add ObjectMode
    Local       bool     `toml:"local"`
}

// Operation in spec.go
type Operation struct {
    Name        string
    Description string
    Params      []string
    Pairs       []string
    Results     []string
    ObjectMode  string // add ObjectMode
    Local       bool
}
```

3. Modify relevant struct and parse func in [go-storage]

```go
// cmd/definitions/type.go
type Operation struct {
    Name        string
    Description string
    Pairs       []string
    Params      Fields
    Results     Fields
    ObjectMode  string // add ObjectMode
    Local       bool
}
```

```go
// cmd/definitions/type.go
func NewOperation(v specs.Operation, fields map[string]*Field) *Operation {
    op := &Operation{
        Name:        v.Name,
        Local:       v.Local,
        ObjectMode:  v.ObjectMode, // pass ObjectMode from specs.Operation
        Description: formatDescription("", v.Description),
    }
    ...
}
```

```go
// cmd/definitions/type.go
// ObjectParamName returns Object's param name.
func (f Fields) ObjectParamName() string {
    for _, v := range f {
        if v.ftype == "Object" || v.ftype == "*Object" {
            return v.Name
        }
    }
    return "o"
}
```

4. Add templates for code generating in [go-storage]

```go
// cmd/definitions/tmpl/service.tmpl
// {{ $fnk }}WithContext {{ $fn.Description }}
func (s *{{$pn}}) {{ $fnk }}WithContext(ctx context.Context, {{$fn.Params.String}}) ({{$fn.Results.String}}) {
    defer func (){
        {{- $path := $fn.Params.PathCaller }}
        {{- if and (eq $path "") (eq $pn "Service") }}
            {{ $path = ",\"\"" }}
        {{- end }}
        err = s.formatError("{{$fn.Name}}", err {{ $path }} )
    }()
    
    {{- template "mode_check" makeSlice $fn.ObjectMode $fn.Params.ObjectParamName }}

    pairs = append(pairs, s.defaultPairs.{{ $fnk }}...)
    var opt pair{{ $pn }}{{ $fnk }}
    ...
}
```

```go
// cmd/definitions/tmpl/service.tmpl
{{- define "mode_check" }}
    {{- $mode := index . 0 | toPascal }}
    {{- $o := index . 1 }}

    {{- if ne $mode ""}}
    if !{{ $o }}.Mode.Is{{ $mode }}() {
        err = services.ObjectModeInvalidError{Expected: Mode{{ $mode }}, Actual: o.Mode}
        return
    }
    {{- end }}
{{- end }}
```

5. Generate code in specific service by running `make generate`, the code would be generated in `generated.go`:

```go
// generated.go
// WriteMultipartWithContext will write content to a multipart.
func (s *Storage) WriteMultipartWithContext(ctx context.Context, o *Object, r io.Reader, size int64, index int, pairs ...Pair) (n int64, part Part, err error) {
    defer func () {
        err = s.formatError("write_multipart", err)
    }()
    
    if !o.Mode.IsPart() {
        err = services.ObjectModeInvalidError{Expected: ModePart, Actual: o.Mode}
        return
    }
    ...
}
```

In this way, we can ensure the `*Object` input must meet the operation, and service do not need to care about the check in specific operation.

## Rationale

### Alternative 1: generate mode check by interface name

We can add a field `Interface` to specify which interface an operation belongs to. Then we can use this to generate different mode check func.

For example: all operations in `Multiparter` should check `IsPart()`, and all operations in `Appender` should check `IsAppend`.

However, this solution can only fit check functions written in `Golang` (or [go-storage]), our specs are designed 
to describe model for different languages. So this is not suitable.

## Compatibility

No break changes

## Implementation

Most of the work would be done by the author of this proposal, including:
- Implement mode check
- Upgrade go-service-*
- Remove all mode check in go-sevice-*

[AOS-25]: ./25-object-mode.md
[AOS-47]: ./47-additional-error-specification.md
[go-storage]: https://github.com/beyondstorage/go-storage

- Author: Xuanwo <github@xuanwo.io>
- Start Date: 2021-05-06
- RFC PR: [beyondstorage/specs#48](https://github.com/beyondstorage/specs/issues/48)
- Tracking Issue: N/A

# AOS-48: Service Registry

## Background

For now, every service implement the function like:

```go
func New(pairs ...typ.Pair) (typ.Servicer, typ.Storager, error) {}
func NewServicer(ps ...types.Pair) (types.Servicer, error) {}
func NewStorager(ps ...types.Pair) (types.Storager, error) {}
```

Users need to handle types by themselves.

## Proposal

So I propose to implement a service registry in go-storage:

From `go-storage` side:

```go
type (
    NewServicerFunc func(ps ...types.Pair) (types.Servicer, error)
    NewStoragerFunc func(ps ...types.Pair) (types.Storager, error)
)

func RegisterServicer(ty string, fn NewServicerFunc) {}
func NewServicer(ty string, ps ...types.Pair) (types.Servicer, error) {}
func RegisterStorager(ty string, fn NewStoragerFunc) {}
func NewStorager(ty string, ps ...types.Pair) (types.Storager, error) {}
```

From services side, we can generate following code:

```go
func init() {
	services.RegisterServicer("<type>", NewServicer)
	services.RegisterStorager("<type>", NewStorager)
}
```

From user side, they can use:

```go
srv, err := NewServicer("<type>", ps...)
store, err := NewStorager("<type>", ps...)
```

## Rationale

### Return function instead

Instead of call function directly, we can return a init function directly:

```go
type (
    NewServicerFunc func(ps ...types.Pair) (types.Servicer, error)
    NewStoragerFunc func(ps ...types.Pair) (types.Storager, error)
)

func NewServicer(ty string) NewServicerFunc {}
func NewStorager(ty string) NewStoragerFunc {}
```

User needs to:

```go
srv, err := NewServicer("<type>").(ps...)
store, err := NewStorager("<type>").(ps...)
```

## Compatibility

No breaking changes.

## Implementation

- Implement service registry in go-storage
- Implement service code generate in go-storage definitions
- Make sure all service has been updated

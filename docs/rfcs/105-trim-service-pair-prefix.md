- Author: xxchan <xxchan22f@gmail.com>
- Start Date: 2021-06-15
- RFC PR: [beyondstorage/specs#105](https://github.com/beyondstorage/specs/issues/105)
- Tracking Issue: [beyondstorage/go-storage#598](https://github.com/beyondstorage/go-storage/pull/598)

# GSP-105: Trim Service Pair Prefix

- Updated By:
  - [GSP-117](./117-rename-service-to-system-as-the-opposite-to-global.md): Rename `service pair` to `system pair`

## Background

Currently, service pairs have the prefix `<type>_` internally, e.g.,

```go
// go-service-s3/generated.go

const (
	// StorageClass
	pairStorageClass = "s3_storage_class"
	// ...
)

// WithStorageClass will apply storage_class value to Options.
//
// StorageClass
func WithStorageClass(v string) Pair {
	return Pair{
		Key:   pairStorageClass,
		Value: v,
	}
}
```

This prefix is intended to help prevent a service pair overwriting a global pair with the same name. But actually, it is not used to achieve this goal.

### No Conflict Now

Now, when a service pair conflicts with a global pair, the generation process will fail with message "pair conflict: ...". In `definitions`, this is checked in `mergePairs`, which is called by `FormatService`, which is called by `actionService`.

```go 
func mergePairs(ms ...map[string]*Pair) map[string]*Pair {
	ans := make(map[string]*Pair)
	for _, m := range ms {
		for k, v := range m {
			if _, ok := ans[k]; ok {
				log.Fatalf("pair conflict: %s", k)
			}
			v := v
			ans[k] = v
		}
	}
	return ans
}
```

So the prefix is not needed.

The prefix is added here (`service.tmpl`): 

```go 
// Service available pairs.
const (
{{- range $_, $v := .Pairs }}
    {{- if not $v.Global }}
    {{ $v.Description }}
    {{ $v.FullName }} = "{{ $.Name }}_{{ $v.Name }}"
    {{- end }}
{{- end }}
)
```

We can safely remove it.

### What about allowing conflict?

We may need to rethink the cases of pair conflict.

First, a new service pair conflicts with an existing global pair. This is possibly a mistake and we should warn the developer.

Second, a new global pair conflicts with an existing service pair. This can happen if a service feature is lifted into a global one. However, after the new global pair added, corresponding service pair should be removed and thus users have to migrate immediately. We can instead let the service and global pairs compatible and mark the service pair deprecated.

## Proposal

So I propose:
- Trim the type prefix of the internal names of service pairs.
- If a service pair has the same name as a global pair 
  - If their types are the same, a warning is logged.
  - If their types are different, generation fails.

## Rationale

The biggest benefit of trimming the prefix is that we can allow pair conflict and support lifting a service pair to a global pair. 

Besides, when implementing other features involving handling a general pair, we won't need to consider the internal name any more, thus maintainability improved.

### Possible Problems

- When we lift a service pair to a global pair, what if different service pairs in different services have different names/types?
  
  In this case, we should not lift it.

- What if two service pairs have the same name, but different types?
  
  We are not intended to put all pairs into a global namespace, but to include global pairs in a service namespace. So service pairs won't conflict.

- What if two pairs can have the same name, but do completely irrelevant things?
  
  If the service pair is new, then the developer sees the warning, and should give it a new name.

  If the global pair is new, the developer won't be noticed, since global pairs are ignorant of service pairs. But we can find the problem in the proposal process. In the worst case where the conflict is introduced, the service will need to bump the major version to fix this problem.

## Compatibility

Only internal implementation is changed, and it won't break users.

## Implementation

See https://github.com/beyondstorage/go-storage/pull/598

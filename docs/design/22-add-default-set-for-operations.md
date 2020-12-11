---
author: Lance <github.com/Prnyself>
status: candidate
updated_at: 2020-12-10
---

# Proposal: Add default set for operations

## Background

Now we have two different scenarios where pairs are used, 
one is in `new servicer or storage`, and the other in a specific operation. 
But there are actually some pairs that are always accompanied by a storager after initialization, 
and we should support default pairs, for different operations. 
Related issue: [proposal: Add default set for operations #170]

## Proposal

So I propose that [storage] should support default pairs for different operations:

### Add default pairs

We can set an [embedded struct](#handle-conflict) for every service.
So that in different service, it can implement its custom pairs when init by conducting a struct and pass as an argument.
This struct is maintained by different service, which defines its own operations, and should be
generated from specs as same as definitions of operations.

### Handle conflict

If both `Service` and `Storage` have operations with duplication of name, such as delete,
we should handle conflict.
Since the operations with the same name handle different targets (`Servicer` and `Storager`),
I think we should handle the separately.

```go
type DefaultStoragePairs struct {
    Copy   []Pair
    Write  []Pair
    Read   []Pair
    Delete []Pair
}

type DefaultServicePairs struct {
    Delete []Pair
}
```

### Define struct separately

Since one service may just implement a part of operations in [storage],
there are two ways to define default pairs struct:

- Unify struct in [storage], so that we can maintain it only once.
- Generate struct in each service, so every service would only have operation fields that implemented by itself.

For the first one, it would be weird to set a field whose operation is not implemented.
For example, the [fs-service] does not implement `WriteIndexSegment` operation, we do not need a field like `WriteIndexSegmentPair`.
All the pairs are generated from specs, so we choose the second way.

### Parse pairs in specific operation

When parsing pairs in specific operation, we should combine default pairs and `pairs from args`,
and make sure that `pairs from args` can overwrite default pair, and this should be **generated**.
For example, we parse the pairs sequentially from the slice of pairs,
so we should create a new slice of pairs and append `pairs from args` to the slice.

```go
// generated code in Write operation
func (s *Storage) WriteWithContext(ctx context.Context, path string, r io.Reader, pairs ...Pair) (n int64, err error) {
	defer func() {
		err = s.formatError("write", err, path)
	}()
	pairs = append(pairs, s.defaultPairs.Write...)
	var opt *pairStorageWrite
	opt, err = s.parsePairStorageWrite(pairs)
	if err != nil {
		return
	}

	return s.write(ctx, path, r, opt)
}
```

```go
// check if pair was set before, if set, skip
for _, v := range opts {
	switch v.Key {
	// Required pairs
	case "name":
		if result.HasName {
			continue
		}
		result.HasName = true
		result.Name = v.Value.(string)
	}
}
```

### Check if pair is valid  

- Plan A: Checked when init. We can generate valid check func like `parsePairStorageWrite` when init,
and response with pair policy.
So that it is pre-checked before call specific operation, at that time,
it may be inconvenient to change the default pairs.
- Plan B: Check validate when call specific operation. The inappropriate pair may be overwritten by args, 
so we should not worry about its validation when init. In this way, we should announce that we do not check
validation when init, and return error when specific operation was called.

For now, we choose the `Plan B`, and keep `Plan A` as candidate if it meets other inconvenience.
 
## Rationale

Reuse pairs after set when init, do not need to set the same value every time
we call specific operation.

## Compatibility

None

## Implementation

Most of the work would be done by the author of this proposal.

[fs-service]: https://github.com/aos-dev/go-service-fs
[storage]: https://github.com/aos-dev/go-storage
[proposal: Add default set for operations #170]: https://github.com/aos-dev/go-storage/issues/170
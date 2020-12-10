---
author: [Lance]
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

We can set a struct for every supported operation in [storage].
So that in different service, it can implement its custom pairs when init.

### Parse pairs in specific operation

When parsing pairs in specific operation, we should combine default pairs and `pairs from args`,
and make sure that `pairs from args` can overwrite default pair. 

## Rationale

Reuse pairs after set when init, do not need to set the same value every time
we call specific operation.

## Compatibility

None

## Implementation

```go
type DefaultPairs struct {
    copy []Pair
    write []Pair
    read []Pair
}
```

[Lance]: https://github.com/Prnyself
[storage]: https://github.com/aos-dev/go-storage
[proposal: Add default set for operations #170]: https://github.com/aos-dev/go-storage/issues/170
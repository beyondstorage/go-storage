- Author: JinnyYi <github.com/JinnyYi>
- Start Date: 2021-08-26
- RFC PR: [beyondstorage/go-storage#725](https://github.com/beyondstorage/go-storage/issues/725)
- Tracking Issue: [beyondstorage/go-storage#0](https://github.com/beyondstorage/go-storage/issues/0)

# GSP-725: Add Defaultable Property for Pair

- Updates:
  - [GSP-700]: Deprecate `defaultable` in namespace and split default pairs into global and system.

## Background

In [GSP-700], we introduced `defaultable` field in `namespace` to specify defaultable pair and generate default pairs for them to support default pairs configuration.

After introducing defaultable pairs for namespace:
- We need to maintain a pairs map in namespace described in [Pair map for namespace](https://github.com/beyondstorage/go-storage/issues/711). There will be cross-association between `service`, `namespace` and `pair`.
- All the generated default paris are system pairs. For global pairs, shared values are only valid inside the service. Users need to assign default values separately for each service during initialization instead of using a unified function conveniently.

## Proposal

Add `defaultable` property for pair:

```toml
[io_callback]
type = "func([]byte)"
description = "specify what todo every time we read data from source"
defaultable = true
```

- The default value of `defaultable` is `false`. `true` means value of the pair is defaultable and can be shared.
  - Defaultable global pair means for all the services, pairs of operation with the same name will share the default value.
  - Defaultable system pair means for the current service, pairs of operation with the same name will share the default value.

### Implementation

**Generate default pairs**

Generate default pairs prefixed with `default_` for the defaultable pairs, and the generated pairs for defaultable global pairs are also global. So that:
- Support config default pairs via connection string is still in effect.
- `WithXxx()` for the default paris will be generated separately for global and service.

## Rationale

### Alternative Implementation

Based on [GSP-700], defaultable pairs maintained in namespace additionally. We can split default pair into global and system according the `Global` property.

When generating `WithXxx()` for pairs, we can handle pairs group by `default` and `global`. `pairs.WithXxx()` will be generated for global default paris and `<service>.WithXxx()` will be generated for system default paris finally. 

## Compatibility

`defaultable` in `namespace` will be deprecated.

## Implementation

- Add `defaultable` field in pair.
- Generate default pairs for defaultable pairs.

[GSP-700]: ./700-config-features-and-defaultpairs-via-connection-string.md

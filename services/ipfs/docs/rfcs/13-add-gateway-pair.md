- Author: zu1k <i@lgf.im>
- Start Date: 2021-07-15
- RFC PR: [beyondstorage/go-service-ipfs#13](https://github.com/beyondstorage/go-service-ipfs/pull/13)
- Tracking Issue: [beyondstorage/go-service-ipfs#14](https://github.com/beyondstorage/go-service-ipfs/issues/14)

# RFC-13: Add Gateway Pair

Releated issue: [beyondstorage/go-service-ipfs#5](https://github.com/beyondstorage/go-service-ipfs/issues/5)

## Background

With IPFS, we can get the `CID` of an object, then we can use it with a public or internal IPFS Gateway to splice out the access link.

When using internal gateway like `127.0.0.1:8080`, the access link are inaccessible from external network.

When using a public gateway, e.g. `ipfs.io`, `cf-ipfs.com`, if the service is an intranet IPFS cluster, it is possible that the files are not accessible from these public gateway as well.

### How to get the `CID` of an object?

We can [stat](https://docs.ipfs.io/reference/http/api/#api-v0-files-stat) the object, the `hash` would be its `CID`.

```
$ ipfs files stat /part1
QmSvxBhBHn2gu9kAdjCTYfUgDNnyBcNoXXVojRyKdBjZUA  # CID
Size: 262144
CumulativeSize: 262158
ChildBlocks: 0
Type: file
```

### How to splice out the access link?

- https://{gateway URL}/ipfs/{CID}/{optional path to resource}
- https://{CID}.ipfs.{gatewayURL}/{optional path to resource}

## Proposal

I propose to add a pair to let the user specify the `gateway`.

- The `type` of `gateway` should be `String`
- The `format` of `gateway` should follow [go-endpoint](https://github.com/beyondstorage/go-endpoint/blob/master/README.md)
- The `value` of `gateway` should be parsed into `HTTP` or `HTTPS`
- Now we use `gateway` only in `Reach` operation

## Rationale

### Why not infer the `gateway` from the `endpoint`?

If the `endpoint` used is not the loopback address, it is possible to infer the `gateway` address.

For example, if a user tells IPFS to listen to `0.0.0.0` and use the public IP to access, the `endpoint` may be `<IP>:5001` or `<domain>:5001`, we can infer the `gateway` as `<IP>:8080` or `<domain>:8080`.

The problem with this approach is that it is very unstable and the IPFS api does not come with a forensic mechanism, so it is almost impossible for users to access directly through the public IP, so the inferred gateway is most likely invalid.

### Why not use the `gateway` in the configuration?

We can get the `gateway` of IPFS through the [config-show api](https://docs.ipfs.io/reference/http/api/#api-v0-config-show).

This method works if the user has modified the default gateway. 

However, we cannot guarantee that the user has modified this configuration item, and it may be more common for users to use reverse proxy to make IPFS public.

## Compatibility

No compatibility issues at this time.

## Implementation

First we define a pair of `String` type with the name `gateway` in `service.toml`. Then we `generate` code and implement the `Reacher` interface.

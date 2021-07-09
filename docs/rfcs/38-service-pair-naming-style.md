- Author: Xuanwo <github@xuanwo.io>
- Start Date: 2021-04-21
- RFC PR: [beyondstorage/specs#38](https://github.com/beyondstorage/specs/issues/38)
- Tracking Issue: N/A

# AOS-38: Service Pair Naming Style

- Updated By:
    - [GSP-117](./117-rename-service-to-system-as-the-opposite-to-global.md): Rename `service pair` to `system pair`

## Background

There are 10 services have been implemented and more services are on the way, and every service may have their only pairs. We need to design a naming style that meets following needs:

- No ambiguity: the pair should not be confused with other pairs.
- Cross Language: the pair's name should not contain the language related details.
- Easy to follow and learn: the pair's naming style should be easy to follow.

## Proposal

So I propose following service pair naming style.

### Scope

the style only applies for service's pairs and infos, global pairs' name should be discussed in related RFCs.

If a pair could be used with no service type check, we can consider adding it into global service, otherwise, we need to keep it as service pair.

```go
// size & offset can be used in any service
_, err = store.Read("abc", pairs.WithSize(1024), pairs.Offset(4096))
// But s3 storage class can only be used in s3's store

var pairs []types.Pair
switch tp {
    case s3.Type:
        pairs = append(pairs, s3.WithStorageClass("STANDARD_IA"))
}
_, err = store.Read("abc", pairs...)
```

`new` operation's pairs are also not included. Those pairs should follow the SDK's option name.

### Rule

We should adopt a style called `API Native Style` which means our pairs name should be native in original API.

For example, AWS S3 support three kinds of [Server-side encryption](https://docs.aws.amazon.com/AmazonS3/latest/userguide/serv-side-encryption.html)

- Server-Side Encryption with Amazon S3-Managed Keys (SSE-S3)
- Server-Side Encryption with Customer Master Keys (CMKs) Stored in AWS Key Management Service (SSE-KMS)
- Server-Side Encryption with Customer-Provided Keys (SSE-C)


For `SSE-S3`, S3 supports following HTTP headers:

- `x-amz-server-side-encryption` (should be `AES256`)

For `SSE-KMS`, S3 supports following HTTP headers:

- `x-amz-server-side-encryption` (should be `aws:kms`)
- `x-amz-server-side-encryption-aws-kms-key-id`
- `x-amz-server-side-encryption-context`
- `x-amz-server-side-encryption-bucket-key-enabled`

For `SSE-C`, S3 supports following HTTP headers:

- `x-amz-server-side-encryption-customer-algorithm`
- `x-amz-server-side-encryption-customer-key`
- `x-amz-server-side-encryption-customer-key-MD5`

So we should add following pairs:

- `server-side-encryption`
- `server-side-encryption-aws-kms-key-id`
- `server-side-encryption-context`
- `server-side-encryption-bucket-key-enabled`
- `server-side-encryption-customer-algorithm`
- `server-side-encryption-customer-key`
- `server-side-encryption-customer-key-md5`

Pair type should use their real type, for `server-side-encryption-customer-key` here, it's real type is `[]byte`.

Because HTTP headers only allow string, so we have to encode it to base64 like boolean `True`/`False` should be `"true"`/`"false"`. But it doesn't mean we have to set all header-related pairs type to string, this makes it much less easy to use.

So we use the `[]byte` and `boolean` for pairs type based on their real type.

## Rationale

N/A

## Compatibility

This naming style will affect all service pairs that not released.

## Implementation

This proposal is a rule and will be followed by all contributors.

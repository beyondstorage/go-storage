- Author: Xuanwo <github@xuanwo.io>
- Start Date: 2021-06-08
- RFC PR: [beyondstorage/specs#86](https://github.com/beyondstorage/specs/pull/86)
- Tracking Issue: [beyondstorage/go-storage#611](https://github.com/beyondstorage/go-storage/issues/611)

# GSP-86: Add Create Link Operation

- Updated By:
  - [GSP-837: Support Feature Flag](./837-support-feature-flag.md): Move `CreateLink` operation to `Storager`

## Background

We have `ModeLink` for Object which means this Object is a link which targets to another Object. A link object could be returned in `Stat` or `List`. But there is no way to create a link object.

As discussed in [Link / Symlink support](https://github.com/beyondstorage/specs/issues/85), link related support is very different in different services:

- fs: Native support for hardlink and symlink
- oss: Native [Symlink](https://help.aliyun.com/document_detail/45126.html) API
- s3: No native support, [x-amz-website-redirect-location](https://docs.aws.amazon.com/AmazonS3/latest/userguide/how-to-page-redirect.html) can redirect pages but only works for website, see <https://stackoverflow.com/questions/35042316/amazon-s3-multiple-keys-to-one-object#comment57863171_35043462>

## Proposal

I propose to add a new `CreateLink` operation:

```go
type Linker interface {
	CreateLink(path, target string, pairs ...Pair) (o *Object, err error)
}
```

- `path` is the path of link object.
  - `path` COULD be relative or absolute path
- `target` is the target path of this link object.
  - `target` COULD be relative or absolute path

As described in [GSP-87],  `CreateLink` could also be virtual functions.

- Service without native support COULD implement `virtual_link` feature.
- Users SHOULD enable this feature by themselves.

## Rationale

### Hardlink & Symlink

We will only support `Symlink` and simplify it to `Link` instead.

### Data corruption issue

A service doesn't have native support for Link, so we implement one for it. After few months, the service supports it.

> How to handle old-style objects?
> How to avoid data corruption

The feature gates will be kept until next major release.

- If user enabled `virtual_link`, service SHOULD run compatible mode: create link via native methods, but allow read link from old-style link object.
- If user didn't enable `virtual_link`, service will run in native as other services.

Finally, those compatible code will be removed in next major release, and `virtual_link` feature will be removed.

## Compatibility

This proposal is compatible.

## Implementation

- Add new operations and new features in specs
- Implement integration test
- Implement Linker for services

[GSP-87]: ./87-feature-gates.md

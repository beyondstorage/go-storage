- Author: JinnyYi <github.com/JinnyYi>
- Start Date: 2021-10-11
- RFC PR: [beyondstorage/go-storage#837](https://github.com/beyondstorage/go-storage/issues/837)
- Tracking Issue: [beyondstorage/go-storage#0](https://github.com/beyondstorage/go-storage/issues/0)

# GSP-837: Support Feature Flag

## Background

In [GSP-669: Feature Lifecycle](./669-feature-lifecycle.md), we specified the feature lifecycle to address the repetitive work and ineffective tracking, accelerating the iteration of new features. However, we encountered new problems: 

We may find problems with the new interface or feature only during the service implementation, and have to redesign and release the version for this purpose.

Moreover, we do have different capabilities and limitations in storage services, for application or user, the capabilities to access the service is not available.

## Proposal

So I propose to support feature flag for operational rampups, and also provide interfaces for users or applications to acquire service capabilities.

### Feature flag

Feature flag is designed to push features through a predictable lifecycle where a feature can easily be created, enhanced, tested, and then cleaned up by being elevated to a full-fledged feature or discarded altogether.

The basic status of a pre-implemented feature might look like this:

- Feature in development: It begins after the first design proposal is approved and covers the development, testing, and enhancement phases, and possibly even the user experience phase.
- Feature to be released: When the functionality is stable and user-friendly, it can be transferred to a supported interface.

#### Add `development` property for interface

Specify whether the interface is in development status when adding a new interface in go-storage.

```toml
[storage_http_signer]
development = true
```

The data type of `development` is bool: `true` means the interface is in development, and `false` means in stable status. The default value is `false`. 

And go-storage will generate development struct for services.

```go
type DevelopmentInterfaces struct {
	StorageHTTPSigner bool
}

func WithDevelopmentInterfaces(v DevelopmentInterfaces) Pair {
    return Pair{
        Key:   DevelopmentInterfaces,
        Value: v,
    }
}
```

Service implementer should check all operations in development status manually when required.

### Service capabilities

- Services SHOULD declare what the service can do and what operations or features the service can provide.
- Application or user could obtain the capabilities of the service before requesting an operation or feature.

#### Add `implemented` property for operation

For services that support only part of the operations in an interface, `implemented` property will identify whether the operation is implemented:

```toml
[namespace.storage.op.query_sign_http_write]
implemented = false
```

And go-storage will generate `` to expose the capabilities for services.

```go
// Get the supported operations and features
func(s *Storage) GetCapabilities() []string {}

func(s *Storage) FeatureEnabled(op string) bool {}
```

## Rationale

### Feature Flag

Feature flags (also known as feature toggles or feature flipping) is a technique in development that allows you to enable or disable a feature without deploying any codes. The idea behind feature flags is to build conditional feature branches into code in order to make logic available only to certain groups of users at a time. If the flag is `on`, new code is executed, if the flag is `off`, the code is skipped.

Feature flag can be used in the following scenarios:

- Adding a new feature to an application.
- Enhancing an existing feature in an application.
- Hiding or disabling a feature.
- Extending an interface.

Feature flags solutions:

- [Unleash](https://www.getunleash.io/) is a feature management lets you turn new features on/off in production with no need for redeployment.
- [go-feature-flag](https://github.com/thomaspoignant/go-feature-flag) is a simple and complete feature flag solution, without any complex backend system to install.
- [Etsy's Feature flagging API](https://github.com/etsy/feature) used for operational rampups and A/B testing.

## Compatibility

No breaking changes.

## Implementation

- Add `development` property for interface and `implemented` property for operation.
- Generate code for services.
- Refactor go-integration-test to run the tests cases according to service capabilities.
- Add `implemented` property for operations in services and pass integration test.
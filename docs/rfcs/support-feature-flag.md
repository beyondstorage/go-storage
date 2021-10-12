- Author: JinnyYi <github.com/JinnyYi>
- Start Date: 2021-10-11
- RFC PR: [beyondstorage/go-storage#0](https://github.com/beyondstorage/go-storage/issues/0)
- Tracking Issue: [beyondstorage/go-storage#0](https://github.com/beyondstorage/go-storage/issues/0)

# Proposal: Support Feature Flag

## Background

We use go-integration-test to execute integration tests on services. And all services should add integration tests for operations declared in `service.toml`.

For some features or behavior, not all services are supported. When we make these behaviors mandatory test cases, there are problems with some integration test cases not passing.

## Proposal

So I propose to support feature flag to enable or disable operations or features selectively.

- Services should declare what the service can do and what operations and features the service can provide.
- The client application has to obtain the capabilities of the service before requesting an operation or feature.

As for the supported operations and features in go-service-*:

- Each operation SHOULD be completely enabled (implemented) or completely disabled (unimplemented).
- Interface behavior consistency, like [Idempotent Storager Delete Operation](./46-idempotent-delete.md), [Write Behavior Consistency](./134-write-behavior-consistency.md), etc should be the basic behavior of the operation and does not need to be marked.
- Definitions for specific behaviors, like [Write Empty File Behavior](./751-write-empty-file-behavior.md), should be treated as features and added into `features.toml`.

## Rationale

### Feature Flag

Feature flags (also known as feature toggles or feature flipping) is a technique in development that allows you to enable or disable a feature without deploying any codes. The idea behind feature flags is to build conditional feature branches into code in order to make logic available only to certain groups of users at a time. If the flag is `on`, new code is executed, if the flag is `off`, the code is skipped.

Feature flag can be used in the following scenarios:

- Adding a new feature to an application.
- Enhancing an existing feature in an application.
- Hiding or disabling a feature.
- Extending an interface.

Feature flags projects:

- [Unleash](https://www.getunleash.io/) is a feature management lets you turn new features on/off in production with no need for redeployment.
- [go-feature-flag](https://github.com/thomaspoignant/go-feature-flag) is a simple and complete feature flag solution, without any complex backend system to install.
- ...

## Compatibility

All integration tests for services could be affected. We could migrate as follows:

- Release new versions for go-storage and go-integration-test.
- All services bump to the new version with all integration test calls updated.

## Implementation  

From go-storage side:

For features:
- Add feature flag for each feature in `features.toml`. For now, features that need to be added are:
  - write_empty_file
- Generate functions to obtain feature list according to `features` for services.
  ```go
  func(s *Storage) GetFeatures() []string {}
  ```
  
For operations:
- Generate capability list for services, and enables all operations in the supported interfaces by default.
  ```go
  type Storage struct {
    // ...
    capabilities []string
  }
  ```
- Generate functions to identify whether the operation in the supported interface is implemented (or enabled) for services.
  ```go
  func(s *Storage) IsEnabled(op string) bool {
    for index := range s.capabilities {
        if op == s.capabilities[index] {
            return true
        }
    }
    return false
  }
  ```

From services side:

- Services SHOULD declare the supported features in `features` in `service.toml`.
- Services SHOULD maintain the capability list. Usually removes the name of an unsupported operation from the capability list.

From go-integration-test side, we should split test cases into basic cases and feature cases.

- Basic cases are for basic functional testing of interfaces. They SHOULD always be executed.
  - For operations in interface, we need to check whether the operation is enabled first.
- Feature cases are for specific features. They SHOULD be executed when the feature is `on`.
  - First, we need to get the feature list of the storage. Then execute the corresponding feature cases according to the feature names.
  - `map[string]FeatureTestFunc` could be used to maintain a mapping relationship between feature name and feature case.
    ```go
    type FeatureTestFunc func(t *testing.T, store types.Storager)

    var featureTestFnMap map[string]FeatureTestFunc
    ```
  
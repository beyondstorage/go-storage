- Author: JinnyYi <github.com/JinnyYi>
- Start Date: 2021-10-11
- RFC PR: [beyondstorage/go-storage#837](https://github.com/beyondstorage/go-storage/issues/837)
- Tracking Issue: [beyondstorage/go-storage#0](https://github.com/beyondstorage/go-storage/issues/0)

# GSP-837: Support Feature Flag

## Background

Behavior consistency is the most important thing for go-storage, and we do have different capabilities and limitations in storage services.

For storage service development, configuration in `service.toml` makes it easy to specify the interfaces and features supported by the service. However, for application, the capabilities to access the service is not available:

We use go-integration-test to execute integration tests on services. For features or specific behavior, not all services are supported. When we make these behaviors mandatory test cases, there are problems with some integration test cases not passing.

## Proposal

So I propose to support feature flag to enable or disable operations or features selectively for users or applications.

### Capabilities in services

Services SHOULD declare what the service can do and what operations and features the service can provide, and expose interfaces to obtain the capabilities. 

For the supported operations and features in storage services:

- Each operation SHOULD be completely enabled (implemented) or completely disabled (unimplemented).
- Interface behavior consistency, like [Idempotent Storager Delete Operation](./46-idempotent-delete.md), [Write Behavior Consistency](./134-write-behavior-consistency.md), etc should be the basic behavior of the operation and does not need to be marked.
- Definitions for specific behaviors, like [Write Empty File Behavior](./751-write-empty-file-behavior.md), should be treated as features and added into `features.toml`.

### Feature flag in applications

Application (or user) could obtain the capabilities of the service before requesting an operation or feature, then use feature flags to enable or disable features selectively.

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
- ...

### Alternative Way

Instead of expose `GetCapabilities` liked interface in service, we can add flags configuration in go-integration test. And the configurations fields include at least `flagName` and `enabled`.
When executing integration tests on services, services need to init with configs first.

## Compatibility

All integration tests for services could be affected. We could migrate as follows:

- Release new versions for go-storage and go-integration-test.
- All services bump to the new version with all integration test calls updated.

## Implementation

### Capabilities in services

From go-storage side:

For features:
- Add features in `features.toml`. For now, features that need to be added are:
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
  
  func(s *Storage) GetCapabilities() []string {
    if 0 != len(s.capabilities) {
        return s.capabilities
    }
    s.capabilities = append(s.capabilities, types.OpStoragerRead)
    // ...
    return s.capabilities
  }
  ```
  
- Generate functions to identify whether the operation in the supported interface is implemented (or enabled) for services.
  
  ```go
  func(s *Storage) IsEnabled(op string) bool {
    s.GetCapabilities()
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

### Feature flag in go-integration-test

From go-integration-test side, it should run the test cases according to the features, and cases related to operations in `Storager` should always be executed.

We could split test cases into basic cases and feature cases:
- Basic cases are for basic functional testing of interfaces. They SHOULD always be executed.
  - For operations in interface, we need to check whether the operation is enabled first.
- Feature cases are for specific features. They SHOULD be executed when the feature is `on`.
  - First, we need to get the feature list of the storage. Then execute the corresponding feature cases according to the feature names.
  - `map[string]FeatureTestFunc` could be used to maintain a mapping relationship between feature name and feature case.
    
    ```go
    type FeatureTestFunc func(t *testing.T, store types.Storager)

    var featureTestFnMap map[string]FeatureTestFunc
    ```
  
# Workflows

go-storage depends on GitHub actions to run all unit tests and integration tests.

For now, we have the following workflows running:

- auto-merge
- build-test
- cross-build
- unit-test
- services-test

## Auto Merge

`auto-merge.yml` is designed to simplify the dependence management of go-storage.

We use `dependabot` to update all our go modules every day. However, dependabot doesn't support auto-merge PRs after all check success. And dependabot doesn't support run post job after upgrade (we do need to run `make build` after `go-storage` has been upgraded).

So `auto-merge.yml` comes for help.

In the first job `metadata`, we will fetch the dependabot metadata. In the next job `dependabot`, We will use this output to check whether this PR is a dependabot PR. If yes, we will run `make build` and then commit all changes.

To make our maintainer jobs easier, we will use `gh` to enable auto-merge and approve this PR which all checks are passed.

With this workflow's help, our maintainer only needs to care about failed PRs.

## Build Test

`build-test.yml` is used to make sure all our contributors have the same coding style.

In this PR, we will:

- `golangci-lint` for static check.
- `gofmt` for format check. (exit with 1 if changes found)
- `make build-all` to build all packages.
- `git diff` to make sure no code changes after build (exit with 1 if changes found)

## Cross Build

go-storage intends to support running all platforms that go supports. So we set cross-build to make sure go-storage is compiled on the target platform.

We will build on

- GOOS=js GOARCH=wasm
- GOARCH=386 on ubuntu/windows (darwin doesn't support 386 since go 1.15)
- GOARCH=arm on ubuntu/windows (darwin doesn't support arm since go 1.15)
- GOARCH=arm64 on ubuntu/windows/darwin

Every case will be tested on our supported go versions (for now, they are go 1.16 and go 1.17).

## Unit Test

Unit test is used to run all unit tests in go-storage. In this workflow, we will

- `make build-all`
- `make test-all`

on all our supported go versions and systems.

## Services Test

Services test is used to test whether this service is implemented correctly. The test cases live at [tests](../../tests). Every service must pass the tests before they are marked as `stable`.

To reduce the duplicated tests, we have the following filters:

```yaml
on:
  push:
    paths:
      - 'services/azblob/**'
    tags-ignore:
      - '**'
    branches:
      - '**'
  pull_request:
    paths:
      - 'services/azblob/**'
```

Take `azblob` as an example, we will only run tests that have changed files under `services/azblob`. As all commits will be tested before merging, we will ignore the tests triggered by `push tags` (which don't support the path filter).

Depending on the services, the test could be running in github hosted runners or our self-hosted runners.

For open-source services, like `minio`, `fs`, `ftp`, we will run the service in github hosted runners.

For business services, like `s3`, `azblob`, `gcs`, we will run the services in self-hosted runners. All our credentials will be stored at 1password which is accessible to go-storage committers and maintainers. In the actions, we will use `1password/load-secrets-action@v1` to load them:

```yaml
- name: Load secret
  uses: 1password/load-secrets-action@v1
  env:
    STORAGE_AZBLOB_CREDENTIAL: op://Engineering/Azblob/testing/credential
    STORAGE_AZBLOB_NAME: op://Engineering/Azblob/testing/name
    STORAGE_AZBLOB_ENDPOINT: op://Engineering/Azblob/testing/endpoint
```

## More details

### Concurrency

The `concurrency` field is a github native feature that ensures only a single job or workflow using the same concurrency group will run at a time. Other jobs will be canceled.

In our case, we set the group to `${{ github.workflow }}-${{ github.ref }}-${{ github.event_name }}`. This means that for each workflow, for each event in each branch/PR, only one is running.

For example, if we push to a branch 10 times, only the last push will be in progress. Others will be canceled.

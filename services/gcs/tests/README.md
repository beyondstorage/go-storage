## How run integration tests

### Run tests locally

Copy example files and update corresponding values.

```shell
cp Makefile.env.exmaple Makefile.env
```

Run tests

```shell
make integration_test
```

### Run tests in CI

Set following environment variables:

```shell
export STORAGE_GCS_INTEGRATION_TEST=on
export STORAGE_GCS_CREDENTIAL=base64:base64-content
export STORAGE_GCS_NAME=bucketname
export STORAGE_GCS_PROJECT_ID=project-id
```

Run tests

```shell
make integration_test
```

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
export STORAGE_GDRIVE_INTEGRATION_TEST=on
export STORAGE_GDRIVE_NAME=demo
export STORAGE_GDRIVE_CREDENTIAL=file:<abs_path_of_credential>
```

Run tests

```shell
make integration_test
```

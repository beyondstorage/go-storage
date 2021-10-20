## How run integration tests

### Run tests locally

Copy example files and update corresponding values.

```shell
cp Makefile.env.example Makefile.env
```

Run tests

```shell
make integration_test
```

### Run tests in CI

Set following environment variables:

```shell
export STORAGE_IPFS_INTEGRATION_TEST=on
export STORAGE_IPFS_ENDPOINT=http:127.0.0.1:5001
export STORAGE_IPFS_GATEWAY=http:127.0.0.1:8080
```

Run tests

```shell
make integration_test
```

# Sevices

## Available services

### Tire 1 support services

Tire 1 services should pass all integration tests, any failure will block release.

- [azblob](./azblob/)
- [cos](./cos/)
- [fs](./fs/)
- [gcs](./gcs/)
- [kodo](./kodo/)
- [oss](./oss/)
- [qingstor](./qingstor/)
- [s3](./s3/)

### Tire 2 support services

Tire 2 services allow to fail one or more integration tests.

- [dropbox](./dropbox/)
  - community contributed services
- [uss](./uss/)
  - uss have limitations for concurrent put/delete: delete a file shortly after put will make server returning an error

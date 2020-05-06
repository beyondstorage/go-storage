name = "oss"

service {

  op "list" {
    required = ["storager_func"]
    optional = null
  }
  op "new" {
    required = ["credential", "endpoint"]
    optional = ["http_client_options"]
  }
}

storage {

  op "list" {
    required = ["storager_func"]
    optional = null
  }
  op "new" {
    required = ["credential", "endpoint"]
    optional = ["http_client_options"]
  }
  op "write" {
    required = ["size"]
    optional = ["checksum", "storage_class"]
  }
}

name = "gcs"

service {

  op "list" {
    required = ["storager_func"]
    optional = null
  }
  op "new" {
    required = ["credential", "project"]
    optional = ["http_client_options"]
  }
}

storage {

  op "list" {
    required = ["storager_func"]
    optional = null
  }
  op "new" {
    required = ["credential", "project"]
    optional = ["http_client_options"]
  }
  op "write" {
    required = ["size"]
    optional = ["checksum", "storage_class"]
  }
}

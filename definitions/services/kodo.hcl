name = "kodo"

service {

  op "create" {
    required = ["location"]
  }
  op "list" {
    required = ["storager_func"]
  }
  op "new" {
    required = ["credential"]
    optional = ["http_client_options"]
  }
}

storage {

  op "list_dir" {
    optional = ["dir_func", "file_func"]
  }
  op "list_prefix" {
    required = ["object_func"]
  }
  op "new" {
    required = ["name"]
    optional = ["work_dir"]
  }
  op "write" {
    required = ["size"]
    optional = ["checksum", "storage_class"]
  }
}

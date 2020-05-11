name = "azblob"

service {

  op "list" {
    required = ["storager_func"]
  }
  op "new" {
    required = ["credential", "endpoint"]
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
  op "read" {
    optional = ["offset", "size"]
  }
  op "write" {
    required = ["size"]
    optional = ["checksum", "storage_class"]
  }
}

pairs {

  pair "storage_class" {
    type = "StorageClass"
  }
}

infos {

  info "object" "meta" "storage-class" {
    type       = "StorageClass"
    zero_value = "\"\""
  }
}

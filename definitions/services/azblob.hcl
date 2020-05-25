name      = "azblob"
implement = ["prefix_lister", "dir_lister"]

service {

  new {
    required = ["credential", "endpoint"]
  }

  op "list" {
    required = ["storager_func"]
  }
}

storage {

  new {
    required = ["name"]
    optional = ["work_dir"]
  }

  op "list_dir" {
    optional = ["dir_func", "file_func"]
  }
  op "list_prefix" {
    required = ["object_func"]
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

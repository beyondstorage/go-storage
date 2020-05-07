name = "qingstor"

service {

  op "create" {
    required = ["location"]
  }
  op "delete" {
    optional = ["location"]
  }
  op "get" {
    optional = ["location"]
  }
  op "list" {
    required = ["storager_func"]
    optional = ["location"]
  }
  op "new" {
    required = ["credential"]
    optional = ["endpoint", "http_client_options"]
  }
}

storage {

  op "list_dir" {
    optional = ["dir_func", "file_func"]
  }
  op "list_prefix" {
    required = ["object_func"]
  }
  op "list_prefix_segments" {
    optional = ["segment_func"]
  }
  op "new" {
    required = ["name"]
    optional = ["location", "work_dir"]
  }
  op "reach" {
    required = ["expire"]
  }
  op "read" {
    optional = ["offset", "size"]
  }
  op "write" {
    required = ["size"]
    optional = ["checksum", "storage_class"]
  }
}

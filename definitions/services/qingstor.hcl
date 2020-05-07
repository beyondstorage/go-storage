name = "qingstor"

service {

  op "list" {
    required = ["storager_func"]
    optional = ["location"]
  }
  op "new" {
    required = ["credential"]
    optional = ["endpoint", "http_client_options"]
  }
  op "create" {
    required = ["location"]
    optional = null
  }
  op "delete" {
    required = null
    optional = ["location"]
  }
  op "get" {
    required = null
    optional = ["location"]
  }
}

storage {
  op "list_dir" {
    optional = ["dir_func","file_func"]
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
  op "read" {
    required = null
    optional = ["offset", "size"]
  }
  op "write" {
    required = ["size"]
    optional = ["checksum", "storage_class"]
  }
  op "reach" {
    required = ["expire"]
  }
}

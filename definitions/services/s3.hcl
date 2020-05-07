name = "s3"

service {

  op "delete" {
    required = null
    optional = ["location"]
  }
  op "get" {
    required = null
    optional = ["location"]
  }
  op "list" {
    required = ["storager_func"]
    optional = null
  }
  op "new" {
    required = ["credential"]
    optional = ["endpoint", "http_client_options"]
  }
  op "create" {
    required = ["location"]
    optional = null
  }
}

storage {
  op "list_dir" {
    optional = ["dir_func","file_func"]
  }
  op "list_prefix" {
    required = ["object_func"]
  }
  op "new" {
    required = ["name"]
    optional = ["work_dir"]
  }
  op "list_prefix_segments" {
    required = ["segment_func"]
    optional = null
  }
}

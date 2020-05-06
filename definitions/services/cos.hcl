name = "cos"

service {

  op "get" {
    required = ["location"]
    optional = null
  }
  op "list" {
    required = ["storager_func"]
    optional = null
  }
  op "new" {
    required = ["credential"]
    optional = ["http_client_options"]
  }
  op "create" {
    required = ["location"]
    optional = null
  }
  op "delete" {
    required = ["location"]
    optional = null
  }
}

storage {

  op "get" {
    required = ["location"]
    optional = null
  }
  op "list_dir" {
    optional = ["dir_func", "file_func"]
  }
  op "new" {
    required = ["credential"]
    optional = ["http_client_options"]
  }
  op "create" {
    required = ["location"]
    optional = null
  }
  op "delete" {
    required = ["location"]
    optional = null
  }
  op "new" {
    required = ["name", "location"]
    optional = ["work_dir"]
  }
}

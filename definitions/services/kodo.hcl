name = "kodo"

service {

  op "create" {
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
}

storage {

  op "create" {
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
  op "new" {
    required = ["name"]
    optional = ["work_dir"]
  }
}

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
  op "read" {
    required = null
    optional = ["offset", "size"]
  }
}

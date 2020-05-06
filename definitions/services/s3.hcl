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
  op "list_prefix_segments" {
    required = ["segment_func"]
    optional = null
  }
}

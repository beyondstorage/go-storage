name = "dropbox"

service {
}

storage {
  op "list_dir" {
    optional = ["dir_func","file_func"]
  }
  op "new" {
    required = ["credential"]
    optional = ["work_dir","http_client_options"]
  }
  op "read" {
    required = null
    optional = ["size"]
  }
  op "write" {
    required = null
    optional = ["size"]
  }
}

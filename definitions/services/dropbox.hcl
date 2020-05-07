name = "dropbox"

storage {

  op "list_dir" {
    optional = ["dir_func", "file_func"]
  }
  op "new" {
    required = ["credential"]
    optional = ["http_client_options", "work_dir"]
  }
  op "read" {
    optional = ["size"]
  }
  op "write" {
    optional = ["size"]
  }
}

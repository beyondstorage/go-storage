name      = "dropbox"
implement = ["dir_lister", "storager"]

storage {

  new {
    required = ["credential"]
    optional = ["work_dir"]
  }

  op "list_dir" {
    optional = ["dir_func", "file_func"]
  }
  op "read" {
    optional = ["size"]
  }
  op "write" {
    optional = ["size"]
  }
}

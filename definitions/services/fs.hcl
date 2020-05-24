name      = "fs"
implement = ["dir_lister", "storager"]

storage {

  new {
    optional = ["work_dir"]
  }

  op "list_dir" {
    optional = ["dir_func", "file_func"]
  }
  op "read" {
    optional = ["offset", "size"]
  }
  op "write" {
    optional = ["size"]
  }
}

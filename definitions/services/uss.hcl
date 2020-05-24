name      = "uss"
implement = ["prefix_lister", "dir_lister", "storager"]

storage {

  new {
    required = ["credential", "name"]
    optional = ["work_dir"]
  }

  op "list_dir" {
    optional = ["dir_func", "file_func"]
  }
  op "list_prefix" {
    required = ["object_func"]
  }
  op "write" {
    required = ["size"]
  }
}

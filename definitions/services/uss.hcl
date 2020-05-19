name = "uss"

storage {

  op "list_dir" {
    optional = ["dir_func", "file_func"]
  }
  op "list_prefix" {
    required = ["object_func"]
  }
  op "new" {
    required = ["credential", "name"]
    optional = ["work_dir"]
  }
  op "write" {
    required = ["size"]
  }
}

name = "fs"

namespace "storage" {
  implement = ["dir_lister"]

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

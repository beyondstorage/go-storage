name = "uss"


storage {
  op "list_dir" {
    optional = ["dir_func","file_func"]
  }
  op "list_prefix" {
    required = ["object_func"]
  }
  op "new" {
    required = ["name", "credential"]
    optional = ["http_client_options", "work_dir"]
  }
}

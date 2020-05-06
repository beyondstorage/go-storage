name = "uss"

service {
}

storage {

  op "new" {
    required = ["name", "credential"]
    optional = ["http_client_options", "work_dir"]
  }
}

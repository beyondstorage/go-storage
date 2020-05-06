name = "fs"

service {
}

storage {

  op "write" {
    required = null
    optional = ["size"]
  }
}

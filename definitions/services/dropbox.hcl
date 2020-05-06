name = "dropbox"

service {
}

storage {

  op "write" {
    required = null
    optional = ["size"]
  }
}

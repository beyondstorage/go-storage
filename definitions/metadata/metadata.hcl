object_meta "content-md5" {
  type = "string"
  display_name = "ContentMD5"
}

object_meta "content-type" {
  type = "string"
}

object_meta "etag" {
  type = "string"
  display_name = "ETag"
}
object_meta "storage-class" {
  type = "storageclass.Type"
  zero_value = "\"\""
}

storage_meta "location" {
  type = "string"
}

storage_statistic "count" {
  type = "int64"
}
storage_statistic "size" {
  type = "int64"
}

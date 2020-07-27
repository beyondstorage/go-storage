
info "object" "meta" "content-md5" {
  type         = "string"
  display_name = "ContentMD5"
}
info "object" "meta" "content-type" {
  type = "string"
}
info "object" "meta" "etag" {
  type         = "string"
  display_name = "ETag"
}
info "object" "meta" "size" {
  type = "int64"
}
info "object" "meta" "updated_at" {
  type = "time.Time"
}
info "storage" "meta" "location" {
  type = "string"
}
info "storage" "statistic" "count" {
  type = "int64"
}
info "storage" "statistic" "size" {
  type = "int64"
}

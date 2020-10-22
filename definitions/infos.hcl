
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
info "object" "meta" "id" {
  type    = "string"
  export  = true
  comment = "ID is the unique key in service."
}
info "object" "meta" "name" {
  type    = "string"
  export  = true
  comment = "Name is the relative path towards service's WorkDir."
}
info "object" "meta" "size" {
  type = "int64"
}
info "object" "meta" "type" {
  type    = "ObjectType"
  export  = true
  comment = "Type should be one of `file`, `stream`, `dir` or `invalid`."
}
info "object" "meta" "updated_at" {
  type = "time.Time"
}
info "storage" "meta" "location" {
  type = "string"
}
info "storage" "meta" "name" {
  type   = "string"
  export = true
}
info "storage" "meta" "work-dir" {
  type   = "string"
  export = true
}
info "storage" "statistic" "count" {
  type = "int64"
}
info "storage" "statistic" "size" {
  type = "int64"
}

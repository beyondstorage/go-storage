
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
  comment = "ID is the unique key in storage."
}
info "object" "meta" "name" {
  type    = "string"
  export  = true
  comment = "Name is either the absolute path or the relative path towards storage's WorkDir depends on user's input."
}
info "object" "meta" "size" {
  type = "int64"
}
info "object" "meta" "target" {
  type    = "string"
  comment = "Target is the symlink target for this object, only exist when object type is link."
}
info "object" "meta" "type" {
  type    = "ObjectType"
  export  = true
  comment = "Type could be one of `file`, `dir`, `link` or `unknown`."
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

// Code generated by go-bindata. DO NOT EDIT.
// sources:
// tmpl/open.tmpl (1.128kB)

package main

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("read %q: %w", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("read %q: %w", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes  []byte
	info   os.FileInfo
	digest [sha256.Size]byte
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _openTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x9c\x92\x4f\xcb\xd3\x40\x10\xc6\xef\xfb\x29\x1e\xca\x7b\x68\x25\xee\xde\x5f\xe9\xc9\x2a\x78\xb0\x15\xec\x41\x10\x91\x6d\x32\x8d\x8b\xcd\xee\x32\xd9\xa4\x94\x90\xef\x2e\x69\xfe\x34\xc5\x62\x1a\x0f\x21\xd9\x79\x26\xcf\xcc\xfc\x76\x94\xc2\x7b\x97\x10\x52\xb2\xc4\x3a\x50\x82\xc3\x05\xa9\x1b\xce\x28\x8d\x86\xb1\x81\xd8\xea\x93\x8a\xb3\x44\xc5\x8e\xa9\x08\xe6\x94\xbf\xc3\x66\x87\xed\x6e\x8f\x0f\x9b\x4f\x7b\x29\xbc\x8e\x7f\xeb\x94\x30\xe8\x42\x98\xcc\x3b\x0e\x58\x0a\x00\x58\xa4\x26\xfc\x2a\x0e\x32\x76\x99\xfa\x56\x68\x7b\x76\x2a\x0f\x8e\x75\x4a\x8b\x09\x5d\x85\x8b\xa7\x7c\x21\xaa\xea\x2d\x58\xdb\x94\xf0\xf2\x33\xc2\x4b\x46\x41\xe3\x75\x0d\x89\xba\x9e\x72\xc8\x89\x4b\x13\x53\xae\xaa\xaa\xfd\x51\x6e\x75\x46\xa8\xeb\xd6\x95\x6c\xd2\x98\xac\x84\x68\x4a\xc1\x79\xb2\x1f\x0b\x1b\xe3\x58\xd8\x78\xe9\x7c\x80\x94\xf2\xcd\xb5\x0b\xf9\x45\x1b\x5e\x61\x99\x73\x89\xce\x5c\x7e\x6d\xcd\x39\xba\x46\xe8\x16\x6f\xdf\x1c\x81\x98\x9b\xc7\xf1\x4a\x88\x52\xf3\x50\xe0\xb3\xf6\x58\x23\xd3\xfe\x7b\x1e\xd8\xd8\xf4\x47\x2f\x54\x93\xc3\x36\xba\x39\x42\xdb\xa4\x1b\xa8\xeb\xa2\x3f\xb5\xb5\x6f\xd9\xf7\x63\xcb\xfd\xc5\xd3\xeb\xdf\xe1\x2d\x9d\xa3\xc1\xbe\xa3\x32\xfa\xac\x47\x80\xfa\xa9\xff\x0b\xd4\x23\x20\x63\xc3\x87\x60\xc6\x09\x4f\x03\xba\x87\xd3\x69\x33\x91\x0c\x7d\x3f\x8b\xa6\xbb\xf8\x29\x34\x33\xb6\x65\x6c\xf9\x18\xce\x28\x61\x2e\x9c\xfb\x5d\x99\x0b\xa7\xef\xfb\x9f\x70\xfe\x04\x00\x00\xff\xff\x7e\xcd\x13\xa0\x68\x04\x00\x00")

func openTmplBytes() ([]byte, error) {
	return bindataRead(
		_openTmpl,
		"open.tmpl",
	)
}

func openTmpl() (*asset, error) {
	bytes, err := openTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "open.tmpl", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xf7, 0x3, 0x7f, 0x26, 0x97, 0xce, 0x81, 0x33, 0xac, 0x33, 0xcf, 0x4c, 0x8a, 0x1e, 0x67, 0x2b, 0x11, 0x8f, 0xc6, 0x4, 0x30, 0x8f, 0xda, 0x4f, 0x3f, 0x8c, 0xb0, 0x33, 0xe9, 0xd4, 0xf8, 0x84}}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetString returns the asset contents as a string (instead of a []byte).
func AssetString(name string) (string, error) {
	data, err := Asset(name)
	return string(data), err
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// MustAssetString is like AssetString but panics when Asset would return an
// error. It simplifies safe initialization of global variables.
func MustAssetString(name string) string {
	return string(MustAsset(name))
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetDigest returns the digest of the file with the given name. It returns an
// error if the asset could not be found or the digest could not be loaded.
func AssetDigest(name string) ([sha256.Size]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s can't read by error: %v", name, err)
		}
		return a.digest, nil
	}
	return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s not found", name)
}

// Digests returns a map of all known files and their checksums.
func Digests() (map[string][sha256.Size]byte, error) {
	mp := make(map[string][sha256.Size]byte, len(_bindata))
	for name := range _bindata {
		a, err := _bindata[name]()
		if err != nil {
			return nil, err
		}
		mp[name] = a.digest
	}
	return mp, nil
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"open.tmpl": openTmpl,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"},
// AssetDir("data/img") would return []string{"a.png", "b.png"},
// AssetDir("foo.txt") and AssetDir("notexist") would return an error, and
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		canonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(canonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"open.tmpl": &bintree{openTmpl, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory.
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	return os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
}

// RestoreAssets restores an asset under the given directory recursively.
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(canonicalName, "/")...)...)
}

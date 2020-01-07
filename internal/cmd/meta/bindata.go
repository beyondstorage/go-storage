// Code generated by go-bindata. DO NOT EDIT.
// sources:
// meta.tmpl (1.896kB)

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
		return nil, fmt.Errorf("read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("read %q: %v", name, err)
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

var _metaTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x54\x5d\x6f\x9b\x30\x14\x7d\xf7\xaf\xb8\x8b\xa2\x89\x54\xd4\xbc\x6f\xca\xcb\xda\x4a\x9b\xa6\x7e\x68\xad\xa6\x49\x55\x54\xb9\x70\x43\x2d\xc0\x66\xb6\xa1\x8d\x18\xff\x7d\xf2\x07\x84\x46\x49\x9a\x27\xb0\xef\x39\xbe\xe7\x1e\x1f\x48\x12\xb8\x90\x19\x42\x8e\x02\x15\x33\x98\xc1\xf3\x06\x72\x39\xae\xa1\xe5\x0c\xb8\x30\xa8\x04\x2b\x93\xb4\xca\x92\x0a\x0d\xfb\x0a\x97\xb7\x70\x73\xfb\x00\x57\x97\x3f\x1e\x28\xa9\x59\x5a\xb0\x1c\xa1\xeb\x80\xde\xb0\x0a\xa1\xef\x09\xe1\x55\x2d\x95\x81\x88\x00\x00\xcc\x52\x29\x0c\xbe\x99\x19\xf1\xcb\x9c\x9b\x97\xe6\x99\xa6\xb2\x4a\xfe\x34\x4c\xbc\xca\x44\x1b\xa9\x58\x8e\xb3\x0f\xea\x49\x5d\xe4\x89\xc6\xbc\x42\x61\x4e\xc2\xa2\xc8\x6a\xc9\x4f\x04\xa7\x0a\x33\x14\x86\xb3\xf2\x43\xb8\xd9\xd4\xa8\x4f\x43\x25\x35\xe3\x4a\xcf\xc8\x82\x90\x96\x29\x78\x82\x6d\x1b\x7a\xa7\x64\xcb\x33\x54\xa1\x32\xa8\xdd\xdd\x0f\x13\xd3\x7b\xff\x1c\x76\x7d\x17\x7a\xef\x9f\x8a\x90\x24\x81\x87\x4d\x8d\xc0\x35\x98\x17\x04\xdb\x1e\xd6\x52\xbd\xbb\x99\x54\x0a\x6d\x3c\x6c\x09\xb3\x49\x65\x46\x48\xd7\xc1\xfc\x92\x19\x06\x5f\x96\x40\xdd\x35\x76\xdd\x39\x28\x26\x72\x84\xb9\x60\x15\xc6\x30\xcf\x86\xba\x03\xf6\xbd\x23\xa5\xb6\x68\x77\x1d\x0a\xfe\x41\xca\x2a\x2c\x2f\x98\x46\x8f\x18\x0f\x29\x62\x98\xb7\x0e\x98\x4d\xe9\x85\xdb\x2a\x76\x89\x6e\x00\x6b\xdf\xb6\x47\xdf\x07\x42\xdf\x83\x36\xaa\x49\x0d\x74\xee\x1a\x92\x04\xee\x14\x9e\x67\xb8\xe6\x02\x33\xc7\xd2\xae\x70\xe1\xb3\x07\x21\x83\x34\xac\xc9\xc0\xba\x46\xc3\xf6\xd0\xa6\xa2\x71\x13\xc3\xfc\xc9\x69\x6c\xad\x2e\x5b\xff\xce\xb4\x55\x52\xe0\x66\xaa\xba\xef\xe1\x59\xca\x32\x9c\xb0\xaf\xec\x0b\x5c\x64\xf8\xe6\xcd\xa6\xf6\x2e\xae\x59\xed\xc1\xe1\x74\xdb\x1d\x45\x66\x97\x3d\x21\xeb\x46\xa4\x50\x33\xa5\x71\x6a\xc4\xdd\x60\x8c\x35\x23\x92\xb5\xd1\x40\x29\x3d\x73\xa1\xa3\xb6\xb8\x80\xe8\xec\xa0\x7b\x31\xa0\x52\x52\x2d\x82\x7d\x0a\x75\x53\x1a\x3b\xe2\xe7\x83\x94\xae\xf7\xa6\xb5\xac\x6c\x50\x5b\x6c\xc5\x0a\x8c\x2a\x56\x3f\x6a\xa3\xb8\xc8\x57\xee\x57\xb1\x66\x29\x76\xfd\xc2\x41\x6d\xfa\x9e\x62\x70\x57\xee\xdd\x74\x3a\x7d\xcf\xed\x59\x8f\x2d\xfd\x89\x9b\x15\x2c\xa1\xa5\xbf\xed\x8e\xab\x87\x76\xd6\x0b\xbe\x9e\x58\x6f\xe3\xdf\xc2\xa4\xd7\xb8\x2b\x8b\xa9\xfd\xa3\x85\x63\x40\xac\x85\x50\xef\x8d\x49\x1b\x5b\xf6\x72\x10\xe4\x0a\x43\x56\x56\x0e\xc1\xd7\x16\xb1\x95\xee\x2d\x1b\x30\x4e\x7b\xb4\x13\x32\x6f\x42\x0f\x58\x6a\x3c\xc6\x1c\x68\xdf\x58\x5a\xe4\x4a\x36\x22\x8b\x16\x53\x0b\x46\xed\xd5\x89\x61\x55\xf8\xb7\xe1\x0a\xb3\xf7\x99\xdd\x3b\xe3\xde\x98\xae\xde\x39\x3f\x9e\x16\xce\xe1\x6b\xf8\xb4\x63\x85\x69\x94\x00\xc1\xcb\x18\x7c\xfe\x6e\xf0\xf5\x4a\x29\x9b\xc2\x5f\x81\x1c\x1d\x69\x37\x0c\xbb\x73\x6f\x47\x5c\x3f\xf8\xf9\x2d\xc1\xa8\x90\x9f\x09\xfc\x10\xb6\xa5\xd1\xf1\xaf\xf1\x88\xb2\x30\xb4\xef\x10\xdb\xe1\x89\xff\xd7\x05\xcc\xe4\xf5\x7f\x00\x00\x00\xff\xff\x83\x22\x98\xd8\x68\x07\x00\x00")

func metaTmplBytes() ([]byte, error) {
	return bindataRead(
		_metaTmpl,
		"meta.tmpl",
	)
}

func metaTmpl() (*asset, error) {
	bytes, err := metaTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "meta.tmpl", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x63, 0x2f, 0xd, 0x9e, 0x5f, 0xbc, 0xa0, 0x56, 0xf4, 0x7a, 0xdb, 0x59, 0xb8, 0xe, 0x1c, 0x90, 0x2b, 0xd0, 0x85, 0x8, 0xb4, 0xaa, 0x59, 0x8d, 0x73, 0xe8, 0xaa, 0xa5, 0x3d, 0x8e, 0x46, 0xe7}}
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
	"meta.tmpl": metaTmpl,
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
	"meta.tmpl": &bintree{metaTmpl, map[string]*bintree{}},
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

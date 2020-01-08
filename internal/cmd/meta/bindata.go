// Code generated by go-bindata. DO NOT EDIT.
// sources:
// meta.tmpl (1.864kB)

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

var _metaTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x54\x5d\x6f\x9b\x30\x14\x7d\xf7\xaf\xb8\x8b\xa2\x89\x54\xd4\xbc\x6f\xca\xcb\xda\x4a\x9b\xa6\x7e\x68\xad\xa6\x49\x55\x54\xb9\x70\x43\x2d\xc0\x66\xb6\xa1\x8d\x98\xff\xfb\x64\x1b\x08\x8d\x92\x34\x4f\xe0\x7b\xce\xfd\x3a\x3e\x90\x24\x70\x21\x33\x84\x1c\x05\x2a\x66\x30\x83\xe7\x0d\xe4\x72\x3c\x43\xcb\x19\x70\x61\x50\x09\x56\x26\x69\x95\x25\x15\x1a\xf6\x15\x2e\x6f\xe1\xe6\xf6\x01\xae\x2e\x7f\x3c\x50\x52\xb3\xb4\x60\x39\x42\xd7\x01\xbd\x61\x15\x82\xb5\x84\xf0\xaa\x96\xca\x40\x44\x00\x00\x66\xa9\x14\x06\xdf\xcc\x8c\x84\x63\xce\xcd\x4b\xf3\x4c\x53\x59\x25\x7f\x1a\x26\x5e\x65\xa2\x8d\x54\x2c\xc7\xd9\x07\x78\x52\x17\x79\xa2\x31\xaf\x50\x98\x93\xb8\x28\xb2\x5a\xf2\x13\xc9\xa9\xc2\x0c\x85\xe1\xac\xfc\x90\x6e\x36\x35\xea\xd3\x58\x49\xcd\xb8\xd2\x33\xb2\x20\xa4\x65\x0a\x9e\x60\xdb\x86\xde\x29\xd9\xf2\x0c\x55\x8f\x0c\xd3\xee\xc6\xfb\x8d\xe9\x7d\x78\x0e\xd1\xd0\x85\xde\x87\xa7\x22\x24\x49\xe0\x61\x53\x23\x70\x0d\xe6\x05\xc1\xb5\x87\xb5\x54\xef\x6e\x26\x95\x42\x9b\x40\x5b\xc2\x6c\x82\xcc\x08\xe9\x3a\x98\x5f\x32\xc3\xe0\xcb\x12\xa8\xbf\xc6\xae\x3b\x07\xc5\x44\x8e\x30\x17\xac\xc2\x18\xe6\xd9\x80\x7b\xa2\xb5\x3e\x29\x75\xa0\x8b\x7a\x16\xfc\x83\x94\x55\x58\x5e\x30\x8d\x81\x31\x16\x29\x62\x98\xb7\x9e\x98\x4d\xd3\x0b\x1f\x2a\x76\x13\xfd\x02\x4e\xbe\x6d\x0f\x6b\xfb\x04\x6b\x41\x1b\xd5\xa4\x06\x3a\x7f\x0d\x49\x02\x77\x0a\xcf\x33\x5c\x73\x81\x99\xcf\xd2\x1e\xb8\x08\xde\x83\xde\x83\xb4\x3f\x93\x21\xeb\x1a\x0d\xdb\x93\x36\x1d\x1a\x37\x31\xcc\x9f\xfc\x8c\xad\x9b\xcb\xe1\xdf\x99\x76\x93\x14\xb8\x99\x4e\x6d\x2d\x3c\x4b\x59\xf6\x15\xf6\xc1\x01\xe0\x22\xc3\xb7\x20\x36\x75\x77\x71\xcd\xea\x40\xee\xab\xbb\xee\x28\x32\x77\xb4\x84\xac\x1b\x91\x42\xcd\x94\xc6\xa9\x10\x77\x83\x30\x4e\x8c\x48\xd6\x46\x03\xa5\xf4\xcc\x9b\x8e\x3a\x70\x01\xd1\xd9\x41\xf5\x62\x40\xa5\xa4\x5a\xf4\xf2\x29\xd4\x4d\x69\xdc\x8a\x9f\x0f\xa6\x74\x36\x88\xd6\xb2\xb2\x41\xed\xb8\x15\x2b\x30\xaa\x58\xfd\xa8\x8d\xe2\x22\x5f\xf9\x5f\xc5\x9a\xa5\xd8\xd9\x85\xa7\x3a\xf7\x3d\xc5\xe0\xaf\x3c\xa8\xe9\xe7\x0c\x3d\xb7\xb5\x1e\x5b\xfa\x13\x37\x2b\x58\x42\x4b\x7f\xbb\x88\xc7\xc7\x76\x0a\x5a\x98\x94\x1e\xa3\xb2\x08\x6a\x8f\x06\x70\x12\x41\xbd\xd7\x06\x6d\xec\xe8\xcb\xa1\xa1\x07\x06\x2f\xac\x3c\x83\xaf\x1d\x63\x3b\x5a\x90\x64\xe0\xf8\xd9\xa2\x1d\x13\x85\x25\x2d\x60\xa9\xf1\x58\xe6\x90\xf6\x8d\xa5\x45\xae\x64\x23\xb2\x68\x31\x5d\x71\x9c\xbd\x3a\xd1\x8c\x0a\xff\x36\x5c\x61\xf6\xde\x93\x7b\x77\xdc\x6b\xc3\xd5\x58\x96\xaf\x27\xd5\xfa\x3a\x7c\x0d\x9f\x76\xa4\x30\x8d\x12\x20\x78\x19\x43\xf0\xd7\x0d\xbe\x5e\x29\xe5\x5c\xf6\xab\x4f\x8e\x8e\xb4\x1b\x96\xdd\xb1\xf6\x11\xd5\x0f\x7e\x5e\x4b\x30\xaa\xf7\xc7\x84\x7e\x88\xdb\xd2\xe8\xf8\xd7\x76\x64\xb2\x7e\xe9\xd0\x21\x76\xcb\x93\xf0\x2f\xeb\x39\x93\xd7\xff\x01\x00\x00\xff\xff\x5a\xd1\x97\x2a\x48\x07\x00\x00")

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
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xde, 0x91, 0x3d, 0xba, 0xf4, 0x2a, 0x29, 0xac, 0x21, 0x5a, 0x3b, 0x1a, 0xf, 0x98, 0x15, 0x15, 0xa7, 0xd7, 0x20, 0x98, 0xab, 0xe4, 0xab, 0x38, 0xcc, 0x96, 0x48, 0x9e, 0xf6, 0x9a, 0x40, 0x7b}}
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

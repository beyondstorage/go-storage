// Code generated by go-bindata. DO NOT EDIT.
// sources:
// cmd/definitions/tmpl/function.tmpl (567B)
// cmd/definitions/tmpl/info.tmpl (1.702kB)
// cmd/definitions/tmpl/object.tmpl (1.907kB)
// cmd/definitions/tmpl/operation.tmpl (2.051kB)
// cmd/definitions/tmpl/pair.tmpl (490B)
// cmd/definitions/tmpl/service.tmpl (10.807kB)

// +build tools

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

var _cmdDefinitionsTmplFunctionTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xcc\x91\x41\x4b\xc4\x30\x10\x85\xef\xfd\x15\x8f\x25\x87\xae\xec\xe6\x07\x08\x9e\x8a\x82\xb0\xc8\xa2\xde\x25\x64\xa7\x6b\xb0\x99\x94\x66\x5a\x17\x62\xfe\xbb\xa4\x55\x17\x11\x0f\xde\x3c\x25\xcc\xcc\x7b\xf3\x3e\x26\xa5\x2d\x06\xc3\x47\x82\x7a\xda\x40\x4d\xb8\xbc\x82\xd2\x37\x23\xdb\x88\x9c\xab\xd2\x76\x2d\x38\x08\xd4\xa4\x6f\x7d\xdf\x91\x27\x16\x3a\x7c\x36\x55\xcb\x2f\xb3\x66\xd2\x77\xc6\x13\xde\x20\xa1\x31\x9e\xba\x65\xa0\x88\xd5\xa4\x77\xc1\x9a\xb9\xd2\x8e\x6c\x51\x47\x5c\xa4\x04\x75\x56\xec\x4d\x5c\x06\xd6\x48\xa9\x58\xe6\x5c\xa7\xa4\x26\xbd\x37\x83\xf1\x51\x3f\x0e\xce\xef\x4c\x14\xfd\x20\x83\xe3\xe3\x35\x1f\xe2\xab\x93\xe7\x26\x78\x6f\x72\x46\xe8\x05\xbd\x71\xc3\x2f\xa6\xa5\x5c\x62\x7e\xdf\xb4\x2c\xb8\xa7\x38\x76\x12\x3f\x8c\xe7\x00\x15\x00\xf4\x86\x9d\xad\x57\x05\xdc\x9d\xa9\x57\xeb\x6a\xa6\xa2\x2e\xd2\x1f\x71\xac\x9c\x60\x03\x0b\x9d\x44\x37\xcb\xbb\xc1\x3f\x66\xdc\x82\xf8\xeb\xca\x3f\xbe\xef\x01\x00\x00\xff\xff\x41\xc8\xae\xbf\x37\x02\x00\x00")

func cmdDefinitionsTmplFunctionTmplBytes() ([]byte, error) {
	return bindataRead(
		_cmdDefinitionsTmplFunctionTmpl,
		"cmd/definitions/tmpl/function.tmpl",
	)
}

func cmdDefinitionsTmplFunctionTmpl() (*asset, error) {
	bytes, err := cmdDefinitionsTmplFunctionTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "cmd/definitions/tmpl/function.tmpl", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xc1, 0xfc, 0xb, 0xf4, 0xfb, 0xf, 0x94, 0xca, 0x64, 0x7a, 0xb8, 0x48, 0xf9, 0x69, 0xf1, 0x88, 0x46, 0x59, 0xd3, 0xbc, 0x0, 0x4e, 0xe7, 0x8, 0x1e, 0xdf, 0x27, 0x59, 0x81, 0x75, 0x52, 0xa2}}
	return a, nil
}

var _cmdDefinitionsTmplInfoTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xd4\x94\x5f\x6f\xd3\x30\x14\xc5\xdf\xfd\x29\x0e\x55\x85\x1a\xd4\x35\x9b\x84\x78\x80\xe5\x69\x1b\x68\x42\xdb\x90\x36\xf1\x82\x10\x72\x92\x9b\xca\x34\xb6\x23\xdb\x89\x56\x32\x7f\x77\xe4\x24\x74\x4d\xf7\x07\x8a\x78\xe1\xcd\x37\xf6\xbd\xe7\x77\x8f\xaf\x13\xc7\x38\xd1\x39\x61\x49\x8a\x0c\x77\x94\x23\x5d\x63\xa9\x37\x31\x32\x99\xc7\x39\x15\x42\x09\x27\xb4\xb2\xef\x70\x7a\x85\xcb\xab\x1b\x9c\x9d\x9e\xdf\x2c\x58\xc5\xb3\x15\x5f\x12\xdc\xba\x22\xcb\x98\x90\x95\x36\x0e\x33\x06\x00\x93\x42\xba\x09\x8b\x18\x6b\xdb\x03\x18\xae\x96\x84\xe9\x6a\x8e\xa9\x50\x85\xb6\x78\x9b\x60\x71\x1e\x56\x17\xbc\x82\xf7\xac\x6d\x31\xb5\x64\x1a\x91\xd1\x25\x97\x14\xf6\xa7\x2b\xdc\xc1\xe9\x13\x2e\xa9\x0c\x47\x58\x1c\xe3\xbd\xa0\x32\x87\x50\x39\xdd\x42\x28\xb4\xed\x76\x92\xf7\x48\x85\x63\x99\x56\x36\x40\xec\xe8\x36\x5d\xcd\x5e\xdd\xfb\x0e\x71\x37\xfd\x3c\xd4\x0d\x24\xcd\xa2\x83\x08\xf2\x9f\xb8\xcd\x78\xd0\x47\x82\xa3\xe3\xe3\xb0\xbb\xea\x81\x0f\x40\x2a\x0f\xcb\x88\xb1\x60\x00\x76\x7b\x18\xa7\x5b\x67\xea\xcc\xa1\x1d\x94\x37\x6c\xdf\x9e\x62\x0b\x18\x37\xeb\xaa\xaf\xe5\xfd\xd6\x97\xfb\x33\x1b\x86\x2e\x8e\xe3\x60\x00\x6a\x4b\x39\xb8\x05\x0f\x91\xe4\x15\x0a\x6d\xa0\xd3\xef\x94\x39\x34\xbc\xac\x69\x8e\x43\x48\xe2\xca\x42\x69\x07\x4b\x6e\x8e\xa3\xe1\x83\x25\xd7\x95\xea\xea\x08\xe5\xde\xbc\xee\x42\x09\xc9\xab\x2f\xd6\x19\xa1\x96\x5f\x85\x72\x64\x0a\x9e\x51\xeb\xd9\xa0\xfc\xbc\xd7\x61\x57\x14\x81\xfe\xec\xb6\x9b\x10\xef\x59\x51\xab\x0c\x33\x89\x57\xcf\xba\x16\xe1\x03\xb9\xbe\xf1\x53\x61\xab\x92\xaf\x07\x37\x66\xd1\xd8\x8f\xc1\x57\x43\xae\x36\x0a\x72\xf1\xc0\xbe\x40\xfa\xa7\x9a\xd7\x4f\x68\x36\x63\xcd\xe8\x37\x85\x06\xa6\x47\x60\x90\xa0\x19\xf1\xb2\x61\xa0\x4a\xdb\xb1\xfe\x03\x77\x66\x23\xd4\x39\x52\xad\xcb\x68\x20\x12\x05\xe4\x22\xdc\xf0\xcb\x3d\x9f\xc0\x8b\x04\x87\x43\x8d\xe7\xdd\x9e\xc3\x99\x9a\xba\x83\x7e\xbb\xd1\x2d\xa8\x3b\xfc\x20\xa3\x3f\x87\x79\xec\x12\x0a\x5e\x5a\xda\xe7\x96\x2e\x6a\xeb\xf6\x9b\x8e\xbf\xee\x3b\x19\xf7\x5d\x71\x25\xb2\x59\x21\xdd\xe2\xba\x32\x42\xb9\x62\x36\x79\x8c\xf5\x23\xa5\x3c\xbd\x7f\xb9\xbf\xee\x5e\x6c\x9e\xdd\x24\x8a\x1e\x5a\xf4\xdf\xcc\x6e\x6f\xe5\x5d\xb2\x9f\x97\x8f\x8e\x7d\xff\x0f\xdb\xf9\xa5\xdd\x2f\x7f\x06\x00\x00\xff\xff\x26\xf5\xd0\xb3\xa6\x06\x00\x00")

func cmdDefinitionsTmplInfoTmplBytes() ([]byte, error) {
	return bindataRead(
		_cmdDefinitionsTmplInfoTmpl,
		"cmd/definitions/tmpl/info.tmpl",
	)
}

func cmdDefinitionsTmplInfoTmpl() (*asset, error) {
	bytes, err := cmdDefinitionsTmplInfoTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "cmd/definitions/tmpl/info.tmpl", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xa6, 0x9d, 0xdc, 0x29, 0x17, 0xc8, 0x8a, 0x14, 0x25, 0x13, 0xa5, 0x85, 0x13, 0x6, 0xf, 0xb3, 0x4a, 0x21, 0xd1, 0x20, 0x4e, 0x6f, 0x92, 0x60, 0xe2, 0x6e, 0x29, 0x96, 0x4a, 0xee, 0x99, 0x66}}
	return a, nil
}

var _cmdDefinitionsTmplObjectTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xcc\x55\x51\x6f\xdb\x36\x10\x7e\x36\x7f\xc5\x37\x23\xd8\xa4\xc1\x91\x92\x6d\xd8\x43\x16\x3f\x0c\x71\xb6\x06\x68\xec\x02\x76\xfb\x4e\x53\x27\x9b\x8d\x44\x0a\x24\xe5\xda\x75\xf4\xdf\x0b\x4a\xb2\x23\x27\x4e\xe1\xa0\x2f\x7d\x23\x79\x77\xdf\xdd\xf7\xdd\x91\x8c\x63\xdc\xe8\x84\xb0\x20\x45\x86\x3b\x4a\x30\xdf\x60\xa1\xf7\x7b\x88\x3c\x89\x13\x4a\xa5\x92\x4e\x6a\x65\xff\xc1\x68\x82\xf1\x64\x86\xdb\xd1\xdd\x2c\x62\x05\x17\x0f\x7c\x41\x70\x9b\x82\x2c\x63\x32\x2f\xb4\x71\x08\x18\x00\xf4\xd3\xdc\xf5\x9b\x95\x93\x39\xb5\x4b\xbb\x51\xa2\xcf\x42\xc6\xe2\x18\xff\x49\xca\x12\x48\x95\xd0\x1a\x52\x41\xcf\x3f\x93\x70\x98\x4b\xc7\x84\x56\xd6\xe3\x6c\xb7\xe7\x30\x5c\x2d\x08\x67\x0f\x03\x9c\xad\x70\x35\x44\x34\xa9\xfd\xee\xc9\x71\x54\x55\x8d\xda\x44\xde\x79\xa0\xed\x16\x67\xab\x68\xcc\x73\xc2\x23\x9c\xfe\xc0\xad\xe0\x19\xaa\x0a\xa5\x54\xee\xef\xbf\x30\xc4\xe5\xf5\xb5\x77\x7a\xf0\xc1\x1e\x9f\x54\xe2\x97\x4d\x49\x0d\x36\xa4\x85\x5b\x12\x6c\xce\xb3\x8c\xac\x43\xa9\xa4\xf3\x25\x2e\xf4\xb9\x75\xda\xf0\x05\x45\x2c\x8e\x7d\xc0\x78\x32\xbb\x9d\x5e\xf9\x15\x70\xde\x86\xff\x66\x91\x7a\x6a\x16\xd3\x77\x93\x8f\xef\x47\x50\xda\x61\x4e\x10\x4b\x4f\x25\x81\x2e\x9d\x95\x09\xc1\x92\x59\x49\x41\x36\x3a\x0c\xc7\xcd\xbf\x63\x2f\xb1\x8f\xd0\x85\xa4\xe4\x99\x59\x5a\x08\xad\x44\x69\x0c\x29\x07\xcb\x53\x8a\x98\x6f\xc0\xce\x6e\x9d\x29\x85\xc3\xf6\x54\xf5\xbc\x9b\x4c\xbd\x6c\x23\xb2\xc2\xc8\xc2\x77\xfa\xc9\xf8\xaa\x61\x2f\xdd\x93\xdf\x6c\x53\x50\xad\x7d\x55\x75\x4e\x9e\x29\xcd\x7a\x71\x0c\x91\x49\x5f\x7d\x2b\xf4\x6e\xa7\xf0\x65\x29\xc5\xb2\xc3\x94\x67\x72\x45\x11\xeb\xb5\x1e\xd3\x46\x7d\xc3\xea\xa4\x71\xec\xa7\x05\xa5\xa5\x04\xdc\x82\xfb\x5d\xce\x0b\xa4\xda\xec\xa6\x69\xc5\xb3\x92\x06\xb8\x40\x4e\x5c\xd9\xba\x13\x96\xdc\x00\x97\xed\x81\x25\x57\x43\xd5\x38\xf5\x88\xb0\x5e\xa2\x15\xd5\x9b\x3f\xff\x60\xbd\xdc\x5b\xfd\xd4\x46\xf7\xa5\xa3\x35\xab\x18\x3b\x45\xd8\x27\x51\x6f\xd7\xf5\xa5\xa8\x2a\x96\x96\x4a\x20\xd0\xf8\xbd\x71\x0d\xf1\x3f\xb9\x56\x5f\x69\x8b\x8c\x6f\x5a\xe9\x82\xf0\x50\x3c\x6c\xeb\x12\x0d\xb9\xd2\x28\xe8\xe8\x85\xd6\xbe\xa8\x17\xe0\xd3\x57\xc0\x57\x87\xe0\xe1\x2e\xa2\xcd\x72\x04\x1e\x43\xac\x0e\x2a\x60\x6d\x3f\x33\x5b\x67\x7f\x0b\xb1\xe0\x20\xf9\x00\x73\xad\xb3\x70\x9f\xda\x3a\xee\x82\xb0\x69\xae\x4c\xa1\x23\xdf\x96\x5f\x4f\xba\xdf\xbf\x0c\x71\xd1\xe2\x7c\x5f\xac\x01\x9c\x29\xa9\x76\xac\x58\x97\x56\xa7\xb2\x47\x7c\x25\xa3\x3f\xf9\xd9\xa9\x23\x52\x9e\x59\x3a\xaa\xf2\x7d\x69\xdd\xdb\xda\xf8\x43\x2c\x87\x87\x2c\x0b\xae\xa4\x08\xd2\xdc\x45\xd3\xc2\x48\xe5\xd2\xa0\xdf\xce\x7d\x07\xa1\xaa\xfc\x45\x6a\x27\xbf\x1f\x86\x2d\xf5\x9f\x63\xa4\x1a\xf2\x8f\xc3\x53\xd8\x1f\x9d\xc1\xe6\x4d\xe9\x3e\x2f\x2f\xca\x16\x99\x56\x14\xac\x3b\x27\xdb\xfd\x13\x76\xca\xef\x72\xb4\xfc\xf5\x31\xd1\x9e\x3d\x8c\x1d\x86\x75\x80\xff\xda\xaa\x6f\x01\x00\x00\xff\xff\x94\x03\xb8\x29\x73\x07\x00\x00")

func cmdDefinitionsTmplObjectTmplBytes() ([]byte, error) {
	return bindataRead(
		_cmdDefinitionsTmplObjectTmpl,
		"cmd/definitions/tmpl/object.tmpl",
	)
}

func cmdDefinitionsTmplObjectTmpl() (*asset, error) {
	bytes, err := cmdDefinitionsTmplObjectTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "cmd/definitions/tmpl/object.tmpl", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x0, 0x64, 0x33, 0x31, 0x2b, 0x9e, 0xa5, 0x1c, 0x46, 0x19, 0x57, 0xa8, 0x78, 0xb3, 0xc9, 0x92, 0xd7, 0x11, 0xee, 0xde, 0xf1, 0xa5, 0xf2, 0x26, 0x22, 0x58, 0x32, 0x8d, 0xfb, 0xb5, 0x58, 0xdc}}
	return a, nil
}

var _cmdDefinitionsTmplOperationTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xcc\x55\x51\x6f\xd4\x30\x0c\x7e\xcf\xaf\xb0\xaa\x7b\x68\xa5\xd1\xbe\x23\xed\x01\xc1\x90\x26\xa1\xad\x02\x21\x1e\x51\xae\xf5\xdd\x22\xda\x24\x38\xbe\xdb\xa6\x92\xff\x8e\x92\x76\x77\xbd\xd2\x51\x0e\x10\xe2\x2d\x89\xed\xcf\xf6\x97\x7c\xb1\x95\xd5\x17\xb9\x45\xe0\x47\x8b\x4e\x08\xd5\x5a\x43\x0c\xa9\x00\x00\x48\x2a\xa3\x19\x1f\x38\xe9\x77\xca\x24\x22\x13\xa2\xeb\x5e\x00\x49\xbd\x45\x58\x7d\xbe\x80\x95\x82\x97\x97\x90\x5f\x6b\x46\xda\xc8\x0a\x1d\x78\x2f\xba\x0e\x56\x2a\x7f\x83\xae\x22\x65\x59\x19\x1d\x0e\x43\x06\x18\x2c\xca\xd9\x46\x3e\xde\xc8\x16\xc1\x7b\x50\x4f\xc1\xd0\xc5\x4c\x21\x83\xda\x80\x21\x48\xf1\x6b\xf0\x8f\x8e\x89\x43\xda\xab\x0a\x29\xc9\x26\xe7\x6c\x48\x6e\xc3\xb9\xf7\x31\xfe\x03\x93\xd2\xdb\x34\x03\x17\x17\x07\x4c\xd4\x75\x28\x64\xd8\x8f\x9b\x30\x36\x74\xb1\x52\xf9\xad\x8d\x0d\x04\x8f\xa2\x88\xd5\x1a\xdb\xa7\xf9\x06\x6c\x4a\xe9\x2a\xd9\x84\x92\x07\xcb\xa4\xc5\x01\x78\x3e\x26\x1d\x2c\x6f\x0d\xb5\x92\x4b\x49\xb2\x0d\xb9\x32\x38\x35\xbc\x47\xb7\x6b\xd8\x7d\x52\x7c\x57\xf6\x77\x73\xd2\x65\x12\x42\xc6\x2c\x69\xc3\x31\xfa\x9d\xe9\xd3\x2c\xd6\x1e\x90\x5f\xf7\x17\xfb\x3b\x6d\x8c\xc2\xd3\x8a\x1f\x60\x78\x23\xf9\x70\x76\xf1\xd7\xbb\x7c\xba\xb6\xd3\x5d\xdc\xb6\x3b\xc7\x57\xed\x1a\xeb\x8f\x5a\xb5\xb6\xc1\x16\x35\x63\x3d\xf7\xc6\xd2\x4c\x78\x21\x8a\x02\x16\x3d\x23\x28\xac\x11\x30\x00\xd7\x58\x03\x1b\xb8\x93\x7b\x84\x8d\xa1\x7b\x49\x35\x54\xa6\xb5\x92\xd5\xba\x41\x38\x60\xc9\xc0\x9d\xcb\xfb\x47\xbe\x9c\xc3\x31\xed\x2a\x86\xce\x0b\xb1\xd9\xe9\x0a\x52\xb7\x1c\x94\x9d\xd5\xee\xb9\xd8\x13\xcd\x0c\x4a\x24\xe4\x1d\x69\x48\x16\x01\x92\x40\xef\xcf\x35\x75\x46\x31\xff\x44\x44\xa7\x9f\xcd\xac\x8c\x90\x08\x2e\xe1\x06\xef\x6f\x2d\x52\xbc\xe2\x1b\xc3\xd7\xc7\xea\xaf\x88\x0c\xa5\xc9\xb8\x5a\xef\x93\x6c\xfa\xe1\x1c\x99\x14\x5e\x3c\x97\xef\x0c\x7a\xfe\x03\x6d\x0e\xe4\xfd\x09\x41\x23\x46\x9e\x88\x1a\xaf\x8e\xdf\x75\xd4\x54\x29\x15\x95\xa6\x51\xd5\xe3\x41\x3c\x11\xe5\x55\xd3\xc0\xda\x98\xe6\x17\xc7\x52\x60\x1e\x67\x46\x4a\x30\x17\x05\x58\xa9\xc8\x05\xa1\x8f\x67\x52\x77\xf0\x3f\x21\x7b\x9a\x71\x6e\x86\x3c\x3f\x40\x66\x8b\xde\x47\x04\x63\xf3\x32\xd6\xb1\x04\x12\x0c\xfb\x09\xee\x01\xf6\x47\x22\xe7\x97\x5e\x7c\x0f\x00\x00\xff\xff\xb4\x30\x5a\x77\x03\x08\x00\x00")

func cmdDefinitionsTmplOperationTmplBytes() ([]byte, error) {
	return bindataRead(
		_cmdDefinitionsTmplOperationTmpl,
		"cmd/definitions/tmpl/operation.tmpl",
	)
}

func cmdDefinitionsTmplOperationTmpl() (*asset, error) {
	bytes, err := cmdDefinitionsTmplOperationTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "cmd/definitions/tmpl/operation.tmpl", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x76, 0x7e, 0x5d, 0x16, 0x44, 0x6e, 0xcf, 0xa, 0x57, 0xf1, 0xe8, 0x1c, 0xfb, 0x8, 0x6e, 0x26, 0x89, 0xc2, 0x13, 0x34, 0xfd, 0xea, 0x2a, 0xe1, 0x4c, 0x12, 0xa1, 0x32, 0x35, 0xc4, 0x9d, 0x75}}
	return a, nil
}

var _cmdDefinitionsTmplPairTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x90\x31\x6b\xeb\x30\x14\x85\x77\xfd\x8a\x83\xf1\x90\x40\x62\x0d\x6f\xcb\xe3\x4d\x2f\x1d\x4a\x21\xc9\x10\xda\xb1\x28\xf2\x8d\x22\x62\x4b\x42\xbe\x76\x6b\x5c\xff\xf7\x22\x3b\x0d\x94\x2e\xf5\xe4\xfb\xdd\x73\x0e\xe7\x4a\x4a\xfc\xf7\x25\xc1\x90\xa3\xa8\x98\x4a\x9c\x7a\x18\x7f\x9f\xa1\xeb\x52\x96\x74\xb6\xce\xb2\xf5\xae\xf9\x8b\xed\x1e\xbb\xfd\x11\x0f\xdb\xc7\x63\x21\x82\xd2\x57\x65\x08\x41\xd9\xd8\x08\x61\xeb\xe0\x23\x63\x21\x00\x20\xd3\xde\x31\xbd\x73\x26\xe6\xd1\x58\xbe\xb4\xa7\x42\xfb\x5a\x2a\xdf\xac\x4b\xea\xa4\xf1\xeb\x86\x7d\x54\x86\x64\xf7\x47\x86\xab\x91\x17\xe6\xa0\x2b\x4b\x8e\xb3\xc9\x55\xfc\xc2\xc7\x7d\xa0\x26\x13\x62\x29\xc4\x30\xac\x11\x95\x33\x84\xfc\x75\x85\xbc\xc3\xe6\x1f\x8a\x43\x2a\x87\x71\x9c\xb6\x79\x70\xaa\xa6\xc4\xf3\xae\xd8\xa5\xdf\x0f\xb0\x3f\xa8\x46\xab\x2a\x69\xa4\xc4\x8b\xe5\xcb\x30\x7c\x29\xc7\x11\x6f\xb6\xaa\xa0\x42\xa8\x7a\x24\x7e\xf3\x8d\x23\x3a\x55\xb5\x04\xf6\xd8\x87\xe9\x75\x0a\x21\xa5\x98\x25\x5b\x6a\x74\xb4\x13\x4e\xb1\xe7\xd6\xe9\x1f\xc1\x8b\xee\x96\x77\xec\x43\x9a\x97\x48\x55\x31\x4c\x97\x47\xe2\x36\xba\x89\xcc\x20\x7d\x4f\xd4\x6f\x90\x7d\x2b\x91\xad\xee\xdb\xe7\x54\x67\x83\x6e\x26\xa3\x98\x2f\x26\x57\xa6\x06\x9f\x01\x00\x00\xff\xff\x5b\x98\x18\x63\xea\x01\x00\x00")

func cmdDefinitionsTmplPairTmplBytes() ([]byte, error) {
	return bindataRead(
		_cmdDefinitionsTmplPairTmpl,
		"cmd/definitions/tmpl/pair.tmpl",
	)
}

func cmdDefinitionsTmplPairTmpl() (*asset, error) {
	bytes, err := cmdDefinitionsTmplPairTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "cmd/definitions/tmpl/pair.tmpl", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xfb, 0xf0, 0x11, 0x20, 0x7, 0x36, 0xf5, 0x1, 0x25, 0x49, 0xa9, 0xe4, 0x49, 0xaf, 0xb0, 0x5f, 0x1d, 0x1e, 0xf1, 0x82, 0x1c, 0x59, 0xf5, 0xcd, 0x76, 0x89, 0x55, 0x8a, 0x8f, 0xbb, 0xe2, 0xbd}}
	return a, nil
}

var _cmdDefinitionsTmplServiceTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x5a\xdd\x6f\xe3\xb8\x11\x7f\xf7\x5f\x31\x67\xf8\x41\x5e\x38\x52\x8b\xbe\xb9\xc8\xc3\x35\xbb\xdd\x06\x77\xdd\x0d\x36\x69\xef\xe1\xee\x10\xd0\xd2\xc8\x66\x43\x91\x3a\x92\x76\x62\x78\xf5\xbf\x17\xfc\xd0\xb7\xe4\xd8\xdb\x60\x77\x0b\xc4\x2f\x91\xc8\xe1\x70\xe6\x37\x9f\xa4\x12\x45\x70\x25\x12\x84\x35\x72\x94\x44\x63\x02\xab\x3d\xac\x45\xf5\x0e\x3b\x4a\x20\xce\x92\x28\xc1\x94\x72\xaa\xa9\xe0\xea\xaf\xf0\xf6\x23\x7c\xf8\x78\x07\xef\xde\x5e\xdf\x85\x93\x9c\xc4\x0f\x64\x8d\x70\x38\x40\xf8\x81\x64\x08\x45\x31\x99\xd0\x2c\x17\x52\x43\x30\x01\x00\x98\xc6\x82\x6b\x7c\xd2\x53\xf7\x46\xc5\x74\xe2\x9e\xd6\x54\x6f\xb6\xab\x30\x16\x59\x44\x84\xba\x48\x70\x17\xad\xc5\x85\xd2\x42\x92\x35\x46\xbb\xbf\x44\xf9\xc3\x3a\x42\x9e\xe4\x82\xf2\x72\xf5\x49\x6b\x62\x89\x09\x72\x4d\x09\x3b\x67\xd5\x46\xeb\x3c\x66\x14\x4f\xdf\x4b\xa1\xdc\xd1\x18\x95\xa3\x0f\x4f\x58\xa1\xf7\xb9\x21\x9f\x4f\x26\x3b\x22\xe1\x1e\x6a\x49\xc3\x1b\x29\x76\x34\x41\xe9\x67\x4a\xbd\xc3\x7f\x13\xb6\x45\x3f\x78\xeb\x38\x95\x34\xe5\xfe\xe1\xad\x7b\x78\x27\xa5\x28\xe7\x6a\x6d\xc2\x8f\xb9\x35\xdc\x64\x12\x45\x70\xb7\xcf\x11\xa8\x02\xbd\x41\x30\xc2\x40\x2a\x64\xcb\x76\xb1\xe0\x4a\x3b\xb2\x4b\x98\x36\x66\xa6\x76\xbd\xdf\x09\xc8\x8e\x50\x46\x56\x0c\x21\x27\x54\xaa\xd0\xaf\x0b\x26\x87\xc3\x05\x48\xc2\xd7\x08\xb3\xfb\x05\xcc\x76\xb0\xbc\x84\xf0\xc6\xd0\x18\xee\x06\x27\x43\x41\x53\xe0\x42\xc3\x6c\x17\xbe\x67\x62\x45\x58\x3d\x67\xc6\xde\xa2\x8a\x25\xb5\x52\xb7\x27\xfe\xbe\x65\xcc\x8b\xe3\xa5\x9b\x95\xe2\xdd\x3b\x8a\x4a\xd8\x72\x27\xe4\x89\xe1\xd1\x78\x9c\x5b\x45\x3e\xae\xfe\x83\xb1\xfe\x27\x6a\x92\x10\x4d\xc0\xd8\x08\x55\x89\x28\x64\xe5\xb8\x81\x47\x58\xd2\x70\x62\xf1\xea\xad\x93\xdb\x58\xc3\x61\x58\xef\x6b\x9e\x8a\x13\xf5\xbe\x80\x59\xce\x8d\xf0\xcb\xcb\x4a\x8f\xcf\xa0\xc5\x0d\x51\x71\x9b\x8e\xa6\x16\x23\xaa\x72\x46\xf6\xa5\xd9\xc0\xff\x1a\x8c\x2e\x47\xc8\x1a\x50\x3c\x0b\xb9\xe3\x54\x14\x9e\xca\xba\x45\x9f\x4d\xe3\xb1\xb0\xe0\xbe\x47\xdd\xc1\xe9\x91\x32\x06\x6b\xd4\x5d\xfc\x52\x29\x32\x3f\x16\x4e\xa2\xc8\x2c\xbe\x80\xbb\x0d\x55\x90\x6e\x79\x6c\xc5\x51\x1b\xb1\x65\x89\xc5\x6d\x85\x10\x13\xc6\x5c\x9e\x2a\x6d\x45\xb3\x9c\x61\x86\x5c\xa3\x0c\xcb\xf5\x08\x12\xf5\x56\x72\xca\xd7\xdd\x1d\xa9\x02\x89\x24\x01\xc1\xd9\x1e\x08\x4f\x3a\xfc\x33\x91\xd0\x94\x62\x12\x4e\x8c\x00\x7d\x4d\x02\x01\x6f\xdc\xc8\xbc\xcb\xf9\x60\x71\x11\xd9\x02\xc4\x83\x31\xa3\x08\xdf\xa3\xf6\x01\x53\x2d\x9f\x5b\x22\x9a\x1a\x9a\x43\x65\x36\x27\x2d\x88\x2c\x0c\xda\x4c\x1d\xb9\x43\xdc\x13\xb5\x09\x0e\x25\xe6\x6a\x18\x73\xd5\xc7\x9c\x72\x2d\x4e\xc3\xdc\x62\x54\x83\x2e\x78\x8c\x0b\xc8\x19\x12\x85\x90\x91\x07\x04\xb5\x95\x08\x84\x31\xb0\x8c\x37\x44\xc1\x0a\x91\xc3\xa3\xa4\x5a\x23\x87\x15\xa6\x42\xa2\x91\xc1\xc3\xd9\x13\xb2\x86\x73\x01\x95\x27\x54\xca\x97\x90\x86\xb7\x7d\x20\x45\x36\x37\x9a\x1f\xcf\x37\x63\x31\x77\x4a\xbc\x45\x11\xfc\x42\xf5\xa6\x15\x06\x16\x52\x92\xe7\x6c\x0f\xad\x7c\x03\x3b\x93\xa5\xc1\xe0\xea\x92\xad\x05\x76\x30\xb4\x2c\x0e\x5d\xc6\xc1\xae\x1d\x61\x73\x30\x2a\xf4\x1d\xc4\x8c\xd6\x83\x3f\xe1\x7e\xd9\xcf\x8c\x8b\x6a\xde\x96\x8e\x25\xec\x16\xde\x89\x5a\x91\xda\x78\x1c\x06\xd1\xb0\x53\x39\x89\xb1\x95\xc1\x34\x66\x39\x33\xad\xc1\xd4\x24\xfe\x7b\x8e\x8f\x53\xeb\x0a\xb7\xcc\xc4\x62\x30\x84\xe5\xdc\x0e\xe2\xe3\x38\x1b\xd5\xe4\x51\x72\xb0\x5a\xf1\x78\x64\xf7\xd2\x55\x9f\x5d\xda\xd5\xd3\xb6\x32\xf5\xbe\xed\x04\x6c\x14\xa7\x3c\xc1\x27\x08\xe1\x4f\x23\xf9\x77\x66\x7d\xbd\x41\xf8\x67\xcb\x3b\x8a\xe0\x2d\xa6\x64\xcb\xb4\x33\x2c\x14\x85\x73\x43\xaa\xcc\x9e\x66\xc2\xd5\x4a\x5b\x53\x54\x8e\x31\x4d\x69\x0c\xc4\x6a\xe1\x6a\xcb\xf0\xfa\xaa\xc4\x94\x02\x74\x2d\xe5\x04\x6a\xa5\xf3\xbe\x3b\xc3\xaf\xbf\x1b\x76\xdd\xd4\x5d\x4c\x4e\x63\x3b\x6c\xb7\x16\xf6\x39\x87\xd9\xee\x68\x71\xe8\xe1\xef\xdc\xa7\x8e\xc7\x36\xfc\xe5\x78\xca\xfb\x68\xfb\x09\x9b\x67\x67\x29\x1f\x8a\x5f\x63\x11\xb3\x49\x05\xa7\x79\x30\x4b\x8a\xa2\xec\x80\x72\x22\x15\x26\x1e\x61\x67\x83\xce\x0a\xb3\xa0\x28\xda\x36\x70\x56\xf4\x78\xda\x91\x28\x82\x4f\xf8\xc7\x96\x4a\x4c\xdc\xec\x10\xa8\x66\xa2\x14\xb7\xa2\xf6\x78\xfd\x83\x28\xbb\x29\xa1\x72\xc8\x76\x00\x2b\x21\x58\x5d\x95\x47\xc9\xea\xe9\x91\x52\xed\xa5\x75\x59\x8a\xb0\xd3\xa4\xad\xa8\x5f\x5a\xda\x93\xe5\x7d\x5f\x1d\x50\x4e\x12\xb8\x26\xff\xda\xf8\x96\x6e\x27\x15\xde\x8c\xf8\x9e\xad\x21\x96\xc2\x25\x79\x65\xc3\xc7\x16\xe5\x37\x23\xfe\xea\x8a\xc6\x31\xae\x81\xc8\x75\xe9\x93\x73\x08\x46\xf8\x2c\x00\xcd\x11\xa1\xac\xab\x12\x95\xc9\x4a\xcb\xcb\xb1\x38\xa9\x4b\x8d\x85\x7d\x09\x66\x97\xb2\x9a\xd8\x3f\x26\x97\xdd\x2f\xc0\x26\x0c\x67\x0c\x2b\x48\xbd\x50\x3d\x52\x1d\x6f\x60\x17\xfe\x84\xfb\xc6\xf0\x70\xc4\x9c\x19\x35\xe6\x17\x9b\x5e\xa4\xb2\x4a\xa3\x0c\x2e\x2b\x12\xdf\x70\x39\x6d\xc3\x67\x5c\xe1\xd0\x5a\x66\x77\x10\x5c\x53\xbe\xc5\xd6\x44\xd1\x7a\x3b\x8d\xf5\x25\x68\xd9\x61\xe3\x17\x1e\x5f\xb5\x73\x67\xc0\x30\xe8\x3a\xdf\xbc\x05\x5a\x23\x62\xc6\xa3\xfc\xcc\x48\x7f\x45\xf8\x38\xc2\x43\x79\xe9\x19\x88\x7b\xb9\xe9\x0c\x8c\x6d\x97\x59\x1e\x2f\x62\xc1\x77\x28\x75\x23\x5a\xc7\x75\x0c\xc6\x35\x9c\x77\xad\xd8\x3a\x98\xbc\x80\x75\x47\x2d\xdc\xb7\xf2\xff\x64\xe9\xd3\xad\x6d\x6f\x54\x5a\x52\x00\x32\x63\x80\x2f\x16\x26\x25\x4c\x8d\xb9\xef\x80\xe7\x14\xcd\x14\x7a\x46\xc2\xa3\x29\xfc\x70\xae\x19\xfc\xc9\x61\x2c\xc3\x17\x8b\xfa\x1a\xc9\x14\x8f\x72\x4b\x7b\x97\x74\x30\x47\x0c\xb5\xfc\xf5\x77\xa5\x25\xe5\xeb\x83\x6f\x93\x5b\x3e\x7a\x51\x14\x50\x54\x1a\x75\x14\x6e\x9e\x5b\x9d\xe0\x0b\xe0\x94\x4d\x8e\xb7\x85\xdf\x77\x4b\x58\xad\x78\x6d\x0a\x5f\x9b\xc2\x6f\xde\x14\x06\x0a\xde\x54\x93\xf3\xd7\x16\xf1\xcb\x5b\xc4\xef\xb4\xc5\xf8\xaa\x4d\xdc\x77\x8a\xc1\xb7\x6c\xb3\x54\xf8\x55\x1b\xad\xef\xaa\xfd\x81\xc3\x81\xa6\xc0\xdd\xe5\xce\xd4\xdf\xc0\x4e\x9b\x88\x36\xd4\x51\xa1\xd9\xf3\x46\x30\x1a\xef\xc3\x1f\x19\x83\xcf\x9f\xdb\x63\x8d\x74\x37\x3e\x75\x7e\x6f\x79\x66\x83\xf3\x2f\xae\xb6\x79\x2e\xa4\x2e\x7b\x1c\x33\xb8\xdc\xf5\x55\xea\x8f\x1c\x0e\x06\xa7\x01\xed\xbf\xbc\x4d\xf4\x6c\x79\xd2\xe1\x7a\xb4\x7d\xf4\xd7\x89\x6d\xa7\x3d\xc9\x54\x67\x9a\x69\xc8\x3d\x5f\x1e\xec\x62\x72\x12\xc8\x83\x67\x88\x2e\x76\xff\x37\x9d\xf5\x89\x8d\xf5\x8b\xf5\xd5\x8d\xfb\xf2\x17\xbb\xf2\x1e\xc0\x78\xfc\xe6\xd8\x7d\xb6\xfc\x59\xc4\xed\xc2\xd4\xbe\x54\x66\x66\xfa\xbe\x94\x75\xf4\x7a\xb9\x42\xc3\xa4\xad\x51\x66\xa7\xb3\x69\x7e\x0d\x7d\xf6\xe2\xba\x66\x7b\x0c\xc8\xe6\xdc\xd0\x69\x05\x7c\xdf\xf8\xcc\x81\x05\x5c\x01\x6c\xc6\xa4\x7d\x1c\xfa\x60\x1b\x45\x25\x79\xfb\x63\x9e\xed\x2d\x63\x89\x06\x14\x02\xfe\xbf\x31\x60\xb5\x2f\xf3\x48\xe8\xda\xb4\xba\x9b\x9c\xe5\xbc\x28\xe6\x8d\x4d\x03\x7b\x11\x1e\xde\x10\x49\x32\x15\xde\x5a\x9f\x35\x14\x7e\xfc\x93\xf5\xc0\xe6\x44\x1d\x24\xb1\x7e\xf2\x97\x15\x66\xd7\xf0\x6f\x24\x7e\x58\x4b\xb1\xe5\x49\x30\xef\x06\x92\x6a\xa4\x9e\x5f\xa8\xde\x5c\xb9\x35\x41\xac\x9f\x16\xd0\x92\xe0\x8a\x30\x86\xb2\xec\x1b\x86\x60\x6a\xac\x3f\x82\xd8\x31\x9d\x3b\x12\x54\x1a\xf8\xb1\x8e\x44\x67\x61\x92\x60\x8a\xd2\x6e\x1e\xcc\x0f\x9d\x74\x66\xd3\x80\xde\x94\x3e\xe1\xd9\xdf\x10\xbd\x71\x4a\x77\x33\xba\x0f\x2d\xc2\x13\x08\xf0\x0f\xbf\x78\x3a\x9d\xfb\xb7\x46\x2d\x98\x0f\xd5\x6d\x97\xe9\xf4\x06\x2e\x61\xba\xf8\x6d\xfa\xdb\xb4\x57\x32\x06\xea\x8f\xf9\xa1\x94\x60\x9a\xa3\x54\xc8\x8c\x68\x9b\xd3\x82\xa9\x53\xdd\xf8\x71\x51\x4c\xed\x49\xa2\xde\xa0\x28\xa0\xb6\x78\x11\xcc\x27\xdd\xec\x60\x6a\x97\xe9\x5e\xdd\x27\x60\x91\x20\x4c\xa7\x03\xb2\xcc\x32\x33\x55\xb6\xba\x35\x71\x3f\x73\x95\x3f\x93\xdb\xbd\x0f\x78\x38\xdd\x32\xfb\xe2\x1b\xc0\xd0\xf0\x08\xaf\x6d\xe6\xb7\x1b\x14\x45\x30\x1f\xa8\x7f\x5e\xed\x32\xaf\xd7\xfb\x5f\xf3\x1d\x61\xd4\x27\xf7\x77\x4f\x39\xc6\x1a\x93\x25\x98\xa9\x06\xcf\x05\xfc\x18\xeb\x2d\x61\xcb\xb6\xf7\x74\xe4\xf1\xe2\x0c\xf5\x19\x26\x50\x9e\x6f\x14\xda\x67\x32\xb8\x04\x92\xe7\xc8\x13\x7b\xe6\x53\x0b\x50\xa1\x8f\x7c\xfb\x81\xb1\x11\x75\x61\x18\xd6\x36\xda\x11\x69\x4e\x6a\x63\xf5\xad\xde\x44\xe4\x7a\x51\xf9\xc3\xd1\x13\xa8\xdd\xbf\xde\x81\xa6\x76\xd9\x0f\x97\xa6\x76\x75\xc0\xee\xa8\xda\xd8\xae\x99\x2d\xda\x79\xf3\x8a\x64\x68\x9b\xf0\x7e\xc6\xb8\x93\x34\xfb\x99\x28\xed\x53\xc7\x3b\x9e\x98\x53\xe7\xe6\x4a\x64\x19\x29\x0a\xa3\x42\x99\x4b\xc6\xd2\x7e\xb7\x34\xbd\x26\xff\xe7\x13\xdd\x57\x71\xbf\x28\x82\xeb\x35\x17\x12\xdd\xcd\x05\x3c\x6e\x28\x43\xd8\x10\x9e\x30\xca\xd7\x60\xed\x66\x14\xf4\xff\x97\xd1\xf2\xda\xfb\x93\x7d\xf6\x1c\xf7\x7b\x01\xcf\xb3\x06\xa1\x9c\xea\x2a\x0b\x9d\xf8\x3f\x1a\xe6\x57\xa5\xa7\x4f\xb8\xa6\x4a\xa3\x1c\xfb\xa7\x00\x19\x98\xc3\xe9\x02\x3e\xe0\xe3\x28\xc9\xbc\x7f\x9b\xf5\xdf\x00\x00\x00\xff\xff\xc4\x11\x06\x35\x37\x2a\x00\x00")

func cmdDefinitionsTmplServiceTmplBytes() ([]byte, error) {
	return bindataRead(
		_cmdDefinitionsTmplServiceTmpl,
		"cmd/definitions/tmpl/service.tmpl",
	)
}

func cmdDefinitionsTmplServiceTmpl() (*asset, error) {
	bytes, err := cmdDefinitionsTmplServiceTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "cmd/definitions/tmpl/service.tmpl", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x74, 0x5, 0x7d, 0x9b, 0x45, 0xb0, 0x1d, 0xe7, 0xe5, 0xe5, 0xce, 0x51, 0x9d, 0xf9, 0xc1, 0x1b, 0x4a, 0xf5, 0x43, 0xd6, 0x4d, 0xd1, 0x12, 0x19, 0x2c, 0xbf, 0x67, 0x57, 0x3c, 0x97, 0x96, 0x39}}
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
	"cmd/definitions/tmpl/function.tmpl":  cmdDefinitionsTmplFunctionTmpl,
	"cmd/definitions/tmpl/info.tmpl":      cmdDefinitionsTmplInfoTmpl,
	"cmd/definitions/tmpl/object.tmpl":    cmdDefinitionsTmplObjectTmpl,
	"cmd/definitions/tmpl/operation.tmpl": cmdDefinitionsTmplOperationTmpl,
	"cmd/definitions/tmpl/pair.tmpl":      cmdDefinitionsTmplPairTmpl,
	"cmd/definitions/tmpl/service.tmpl":   cmdDefinitionsTmplServiceTmpl,
}

// AssetDebug is true if the assets were built with the debug flag enabled.
const AssetDebug = false

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
	"cmd": {nil, map[string]*bintree{
		"definitions": {nil, map[string]*bintree{
			"tmpl": {nil, map[string]*bintree{
				"function.tmpl":  {cmdDefinitionsTmplFunctionTmpl, map[string]*bintree{}},
				"info.tmpl":      {cmdDefinitionsTmplInfoTmpl, map[string]*bintree{}},
				"object.tmpl":    {cmdDefinitionsTmplObjectTmpl, map[string]*bintree{}},
				"operation.tmpl": {cmdDefinitionsTmplOperationTmpl, map[string]*bintree{}},
				"pair.tmpl":      {cmdDefinitionsTmplPairTmpl, map[string]*bintree{}},
				"service.tmpl":   {cmdDefinitionsTmplServiceTmpl, map[string]*bintree{}},
			}},
		}},
	}},
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

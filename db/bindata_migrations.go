// Code generated by go-bindata.
// sources:
// db/migrations/1-initial.sql
// db/migrations/2-userchanges.sql
// db/migrations/3-event.sql
// db/migrations/4-stream-activelisteners.sql
// db/migrations/5-username.sql
// db/migrations/6-user-stream.sql
// DO NOT EDIT!

package db

import (
	"bytes"
	"compress/gzip"
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
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
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

var _dbMigrations1InitialSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb4\x92\xbf\x4e\xf4\x30\x10\xc4\x7b\x3f\xc5\x96\xdf\x27\x38\x89\x3e\x2d\xaf\x40\x6d\x39\xf1\x40\x56\xf8\x1f\xbb\x6b\xb8\x7b\x7b\x94\x0b\x82\x44\x5c\x4a\x3a\x6b\x7f\xb3\xf2\xcc\xd8\xa7\x13\xdd\x65\x7e\x91\x60\xa0\xa7\xe6\x26\xc1\x72\xb2\x30\x26\x90\x9a\x20\x64\xaf\x90\x77\xc8\x3f\x47\x44\xc4\x91\x14\xc2\x21\x51\x13\xce\x41\x2e\xf4\x8a\xcb\xfd\x15\xe9\x5c\xc5\x3c\x47\x32\x9c\x8d\x4a\x35\x2a\x3d\x25\xea\x85\xdf\x3a\x56\x49\x09\x19\x7b\xbc\xce\xe7\xaa\x76\xc4\xb8\x18\xa4\x84\xe4\xb9\xdd\xc2\x38\x1f\x60\xf7\x7f\x70\xfb\x38\x5d\x21\xfa\x97\x31\x22\x74\x12\x6e\xc6\xb5\xdc\xc2\xab\x99\xe8\x83\x91\x71\x86\x5a\xc8\xed\xe7\x82\x88\xe7\xd0\x93\xd1\xd4\x45\x50\xcc\x7f\x4b\xd6\xe5\xd6\xc7\xc4\x3a\x23\x2e\x85\xfc\xde\x7a\xf8\x32\xdf\xc7\xc5\xc2\x78\x28\xbb\xb6\xe2\xb6\x8f\xfe\x58\x3f\x8a\x8b\x52\xdb\xb6\xa5\x61\x3b\xd9\x7d\x83\xe1\x33\x00\x00\xff\xff\x8f\xd2\xca\x99\x31\x02\x00\x00")

func dbMigrations1InitialSqlBytes() ([]byte, error) {
	return bindataRead(
		_dbMigrations1InitialSql,
		"db/migrations/1-initial.sql",
	)
}

func dbMigrations1InitialSql() (*asset, error) {
	bytes, err := dbMigrations1InitialSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "db/migrations/1-initial.sql", size: 561, mode: os.FileMode(420), modTime: time.Unix(1491226672, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _dbMigrations2UserchangesSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\xcf\x41\x0a\xc2\x30\x10\x85\xe1\x7d\x4f\x31\x7b\xe9\x09\xba\xf5\x0a\xae\x65\x9a\x3c\x25\x34\x99\x94\xc9\x54\x3d\xbe\x08\x2a\x11\x8a\xc9\xfe\xe7\x3d\xbe\x71\xa4\x43\x0a\x57\x65\x03\x9d\xd6\x81\xa3\x41\xc9\x78\x8e\xa0\xad\x40\x0b\xb1\xf7\xe4\x72\xdc\x92\x10\x12\x87\x48\x86\x87\x4d\xff\x43\x09\x6e\x11\x4e\xe8\x69\xf9\xc6\xc6\xba\x69\xd7\xf0\x65\x0e\xbe\x6b\xd4\x39\x94\x72\xb6\xbc\x40\xde\xfd\x50\x4b\x8f\xf9\x2e\x3b\x0b\x5e\xf3\xfa\x83\xdd\xbb\xa9\xa3\x0f\xb4\xd5\x7d\x91\xad\xf0\x05\x6c\x8e\x55\xb8\xe9\x19\x00\x00\xff\xff\x01\x3a\x88\x70\xc0\x01\x00\x00")

func dbMigrations2UserchangesSqlBytes() ([]byte, error) {
	return bindataRead(
		_dbMigrations2UserchangesSql,
		"db/migrations/2-userchanges.sql",
	)
}

func dbMigrations2UserchangesSql() (*asset, error) {
	bytes, err := dbMigrations2UserchangesSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "db/migrations/2-userchanges.sql", size: 448, mode: os.FileMode(420), modTime: time.Unix(1490194356, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _dbMigrations3EventSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xc4\x53\xb1\x6e\xeb\x30\x0c\xdc\xfd\x15\x1c\x63\xbc\x17\xe0\xed\x59\xdf\x2f\x74\x16\x64\x8b\x49\x88\x5a\x94\x4a\x51\x49\xf3\xf7\x85\xac\xb8\x89\x1d\xa7\x43\x3b\x74\x33\xee\x48\xf3\x74\x47\x6e\xb7\xf0\xc7\xd3\x41\xac\x22\xbc\xc4\xa6\x17\x2c\x5f\x6a\xbb\x01\x01\x4f\xc8\xba\x69\x00\x00\xc8\x41\x42\x21\x3b\x40\x14\xf2\x56\x2e\xf0\x8a\x97\xbf\x23\x95\x8e\x41\xd4\x90\x03\xc5\x77\x05\x0e\x0a\x9c\x87\x01\x32\xd3\x5b\xc6\x5a\xc2\xd6\xe3\x9c\xae\xb8\xc3\xd4\x0b\x45\xa5\xc0\x6b\x74\x15\xe3\x8c\x55\x50\xf2\x98\xd4\xfa\x78\x1b\xe0\x70\x6f\xf3\xa0\xd0\x67\x11\x64\x35\x9f\x25\x57\x59\x6a\x45\x0b\x06\x0b\x02\xd9\xad\xc1\x92\x99\x89\x0f\x86\xc3\x19\x88\xb5\x82\xde\x6a\x7f\x24\xb7\x00\xb2\x0c\xa3\xdc\xa6\xdd\x35\x73\xc7\x92\x0a\x5a\xff\x8b\x96\x8d\xaf\xfe\xae\x65\xc8\xee\x07\x6e\x6b\x4e\xc5\xa7\xc7\x96\x7f\xd7\x92\xdc\x15\xe9\x1d\x8a\xe9\x43\x66\xfd\xb2\x38\xe6\x6e\xa0\x74\x34\x93\xd5\xcb\x67\x4e\xff\x7a\x5a\x30\x06\x61\x12\xca\x09\xc5\xd4\x04\xd7\xb6\x2b\x3c\x21\x55\x2c\xa7\x58\x52\x7a\x32\x60\xbc\x8d\x87\x08\x2b\xb7\x0f\x82\x74\xe0\x92\xf7\xe6\x36\xa5\x05\xc1\x3d\x0a\x72\x8f\x09\x72\x42\x49\x1b\x72\xed\x63\xc7\x52\xfa\xac\x6f\x46\xae\xf7\x4f\xca\x66\x7d\xf5\x94\xa7\xbd\x6b\xcb\xe6\xde\x5f\xfe\xff\x70\xe6\xc6\x49\x88\xf7\x97\xbf\xbb\x47\xea\xe0\xdd\x47\x00\x00\x00\xff\xff\x48\x52\xe0\x1e\x2f\x04\x00\x00")

func dbMigrations3EventSqlBytes() ([]byte, error) {
	return bindataRead(
		_dbMigrations3EventSql,
		"db/migrations/3-event.sql",
	)
}

func dbMigrations3EventSql() (*asset, error) {
	bytes, err := dbMigrations3EventSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "db/migrations/3-event.sql", size: 1071, mode: os.FileMode(420), modTime: time.Unix(1491226663, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _dbMigrations4StreamActivelistenersSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x74\xcc\xb1\x0d\x02\x31\x0c\x05\xd0\x3e\x53\xfc\x1e\x9d\x44\x7f\x2d\x2b\x50\x23\x73\x31\x28\x92\x63\x47\xce\x0f\xac\x4f\x4b\x01\x0b\xbc\x6d\xc3\xa9\xb7\x67\x0a\x15\xd7\x51\xc4\xa8\x09\xca\xdd\x14\x93\xa9\xd2\x8b\xd4\x8a\x23\x6c\x75\x87\x1c\x6c\x2f\xbd\x59\x9b\x54\xd7\x9c\x68\x4e\x78\x10\xbe\xcc\x50\xf5\x21\xcb\x88\xf3\x5e\xca\x37\x7b\x89\xb7\xff\x82\x6b\xc6\xf8\x27\xef\x9f\x00\x00\x00\xff\xff\xd7\x36\xcd\xf0\x98\x00\x00\x00")

func dbMigrations4StreamActivelistenersSqlBytes() ([]byte, error) {
	return bindataRead(
		_dbMigrations4StreamActivelistenersSql,
		"db/migrations/4-stream-activelisteners.sql",
	)
}

func dbMigrations4StreamActivelistenersSql() (*asset, error) {
	bytes, err := dbMigrations4StreamActivelistenersSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "db/migrations/4-stream-activelisteners.sql", size: 152, mode: os.FileMode(420), modTime: time.Unix(1489843670, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _dbMigrations5UsernameSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x5c\xcc\x31\x8e\xc2\x40\x0c\x05\xd0\xde\xa7\xf8\x5d\x8a\x55\x4e\x90\x76\xaf\xc0\x01\x0c\xf3\x41\x91\x3c\x9e\xe0\xb1\x05\xc7\x47\xa2\x42\x69\x5f\xf1\xd6\x15\x7f\x7d\x7f\x84\x26\x71\x39\x44\x2d\x19\x48\xbd\x1a\x51\x93\x31\x45\x5b\xc3\x6d\x58\x75\xff\x82\x6b\x27\x92\xef\x84\x8f\x84\x97\x19\x1a\xef\x5a\x96\x58\x16\x94\xef\xcf\xe2\x26\xf2\xdb\xfe\x8f\x97\x9f\xe2\xe0\x94\x16\xe3\x38\xcf\xdb\x27\x00\x00\xff\xff\x07\x5a\x5c\xc0\x8f\x00\x00\x00")

func dbMigrations5UsernameSqlBytes() ([]byte, error) {
	return bindataRead(
		_dbMigrations5UsernameSql,
		"db/migrations/5-username.sql",
	)
}

func dbMigrations5UsernameSql() (*asset, error) {
	bytes, err := dbMigrations5UsernameSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "db/migrations/5-username.sql", size: 143, mode: os.FileMode(420), modTime: time.Unix(1491320025, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _dbMigrations6UserStreamSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x8f\x3d\x0e\xc2\x30\x0c\x46\xf7\x9e\xe2\xdb\x3a\xd0\x9e\xa0\x2b\x57\x60\x46\x26\x31\x01\xc9\x49\x2a\xc7\x16\x1c\x9f\x81\xa5\x84\x1f\x89\xd9\xf6\xf3\x7b\xf3\x8c\x5d\xbe\x26\x25\x63\x1c\xd6\x81\xc4\x58\x61\x74\x12\x86\x37\xd6\x36\x00\x00\xc5\x88\x50\xc5\x73\x81\x50\x49\x4e\x89\x61\x7c\x37\x94\x6a\x28\x2e\x82\xc8\x67\x72\x31\x8c\xe3\x32\xbc\x40\x9a\x29\x53\xfe\x9b\x32\xf5\x17\x99\x2c\x5c\x8e\xad\xba\x86\x9f\xbf\xb7\x3d\xfb\x7a\x2b\x5d\x91\xf2\xb3\x28\x6a\x5d\x7b\x99\xef\xe6\x9f\xb6\xa7\xb7\xc9\x56\x71\x79\x04\x00\x00\xff\xff\xb2\x3f\x7f\x26\x58\x01\x00\x00")

func dbMigrations6UserStreamSqlBytes() ([]byte, error) {
	return bindataRead(
		_dbMigrations6UserStreamSql,
		"db/migrations/6-user-stream.sql",
	)
}

func dbMigrations6UserStreamSql() (*asset, error) {
	bytes, err := dbMigrations6UserStreamSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "db/migrations/6-user-stream.sql", size: 344, mode: os.FileMode(420), modTime: time.Unix(1491337051, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
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

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
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
	"db/migrations/1-initial.sql": dbMigrations1InitialSql,
	"db/migrations/2-userchanges.sql": dbMigrations2UserchangesSql,
	"db/migrations/3-event.sql": dbMigrations3EventSql,
	"db/migrations/4-stream-activelisteners.sql": dbMigrations4StreamActivelistenersSql,
	"db/migrations/5-username.sql": dbMigrations5UsernameSql,
	"db/migrations/6-user-stream.sql": dbMigrations6UserStreamSql,
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
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
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
	"db": &bintree{nil, map[string]*bintree{
		"migrations": &bintree{nil, map[string]*bintree{
			"1-initial.sql": &bintree{dbMigrations1InitialSql, map[string]*bintree{}},
			"2-userchanges.sql": &bintree{dbMigrations2UserchangesSql, map[string]*bintree{}},
			"3-event.sql": &bintree{dbMigrations3EventSql, map[string]*bintree{}},
			"4-stream-activelisteners.sql": &bintree{dbMigrations4StreamActivelistenersSql, map[string]*bintree{}},
			"5-username.sql": &bintree{dbMigrations5UsernameSql, map[string]*bintree{}},
			"6-user-stream.sql": &bintree{dbMigrations6UserStreamSql, map[string]*bintree{}},
		}},
	}},
}}

// RestoreAsset restores an asset under the given directory
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
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
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
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}


// Code generated by "esc -o ./specfs/fs.go -pkg specfs api.json"; DO NOT EDIT.

package specfs

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

type _escLocalFS struct{}

var _escLocal _escLocalFS

type _escStaticFS struct{}

var _escStatic _escStaticFS

type _escDirectory struct {
	fs   http.FileSystem
	name string
}

type _escFile struct {
	compressed string
	size       int64
	modtime    int64
	local      string
	isDir      bool

	once sync.Once
	data []byte
	name string
}

func (_escLocalFS) Open(name string) (http.File, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	return os.Open(f.local)
}

func (_escStaticFS) prepare(name string) (*_escFile, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	var err error
	f.once.Do(func() {
		f.name = path.Base(name)
		if f.size == 0 {
			return
		}
		var gr *gzip.Reader
		b64 := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(f.compressed))
		gr, err = gzip.NewReader(b64)
		if err != nil {
			return
		}
		f.data, err = ioutil.ReadAll(gr)
	})
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fs _escStaticFS) Open(name string) (http.File, error) {
	f, err := fs.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.File()
}

func (dir _escDirectory) Open(name string) (http.File, error) {
	return dir.fs.Open(dir.name + name)
}

func (f *_escFile) File() (http.File, error) {
	type httpFile struct {
		*bytes.Reader
		*_escFile
	}
	return &httpFile{
		Reader:   bytes.NewReader(f.data),
		_escFile: f,
	}, nil
}

func (f *_escFile) Close() error {
	return nil
}

func (f *_escFile) Readdir(count int) ([]os.FileInfo, error) {
	if !f.isDir {
		return nil, fmt.Errorf(" escFile.Readdir: '%s' is not directory", f.name)
	}

	fis, ok := _escDirs[f.local]
	if !ok {
		return nil, fmt.Errorf(" escFile.Readdir: '%s' is directory, but we have no info about content of this dir, local=%s", f.name, f.local)
	}
	limit := count
	if count <= 0 || limit > len(fis) {
		limit = len(fis)
	}

	if len(fis) == 0 && count > 0 {
		return nil, io.EOF
	}

	return fis[0:limit], nil
}

func (f *_escFile) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f *_escFile) Name() string {
	return f.name
}

func (f *_escFile) Size() int64 {
	return f.size
}

func (f *_escFile) Mode() os.FileMode {
	return 0
}

func (f *_escFile) ModTime() time.Time {
	return time.Unix(f.modtime, 0)
}

func (f *_escFile) IsDir() bool {
	return f.isDir
}

func (f *_escFile) Sys() interface{} {
	return f
}

// FS returns a http.Filesystem for the embedded assets. If useLocal is true,
// the filesystem's contents are instead used.
func FS(useLocal bool) http.FileSystem {
	if useLocal {
		return _escLocal
	}
	return _escStatic
}

// Dir returns a http.Filesystem for the embedded assets on a given prefix dir.
// If useLocal is true, the filesystem's contents are instead used.
func Dir(useLocal bool, name string) http.FileSystem {
	if useLocal {
		return _escDirectory{fs: _escLocal, name: name}
	}
	return _escDirectory{fs: _escStatic, name: name}
}

// FSByte returns the named file from the embedded assets. If useLocal is
// true, the filesystem's contents are instead used.
func FSByte(useLocal bool, name string) ([]byte, error) {
	if useLocal {
		f, err := _escLocal.Open(name)
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(f)
		_ = f.Close()
		return b, err
	}
	f, err := _escStatic.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.data, nil
}

// FSMustByte is the same as FSByte, but panics if name is not present.
func FSMustByte(useLocal bool, name string) []byte {
	b, err := FSByte(useLocal, name)
	if err != nil {
		panic(err)
	}
	return b
}

// FSString is the string version of FSByte.
func FSString(useLocal bool, name string) (string, error) {
	b, err := FSByte(useLocal, name)
	return string(b), err
}

// FSMustString is the string version of FSMustByte.
func FSMustString(useLocal bool, name string) string {
	return string(FSMustByte(useLocal, name))
}

var _escData = map[string]*_escFile{

	"/api.json": {
		name:    "api.json",
		local:   "api.json",
		size:    12873,
		modtime: 1586569523,
		compressed: `
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
`,
	},
}

var _escDirs = map[string][]os.FileInfo{}

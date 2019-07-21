// Code generated by "esc -private -o ./swagger-assets.go -pkg httpintf ./swagger.json"; DO NOT EDIT.

package httpintf

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

// _escFS returns a http.Filesystem for the embedded assets. If useLocal is true,
// the filesystem's contents are instead used.
func _escFS(useLocal bool) http.FileSystem {
	if useLocal {
		return _escLocal
	}
	return _escStatic
}

// _escDir returns a http.Filesystem for the embedded assets on a given prefix dir.
// If useLocal is true, the filesystem's contents are instead used.
func _escDir(useLocal bool, name string) http.FileSystem {
	if useLocal {
		return _escDirectory{fs: _escLocal, name: name}
	}
	return _escDirectory{fs: _escStatic, name: name}
}

// _escFSByte returns the named file from the embedded assets. If useLocal is
// true, the filesystem's contents are instead used.
func _escFSByte(useLocal bool, name string) ([]byte, error) {
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

// _escFSMustByte is the same as _escFSByte, but panics if name is not present.
func _escFSMustByte(useLocal bool, name string) []byte {
	b, err := _escFSByte(useLocal, name)
	if err != nil {
		panic(err)
	}
	return b
}

// _escFSString is the string version of _escFSByte.
func _escFSString(useLocal bool, name string) (string, error) {
	b, err := _escFSByte(useLocal, name)
	return string(b), err
}

// _escFSMustString is the string version of _escFSMustByte.
func _escFSMustString(useLocal bool, name string) string {
	return string(_escFSMustByte(useLocal, name))
}

var _escData = map[string]*_escFile{

	"/swagger.json": {
		name:    "swagger.json",
		local:   "./swagger.json",
		size:    4824,
		modtime: 1563728133,
		compressed: `
H4sIAAAAAAAC/+xXT3Pctg6/+1Ng9N5MLutdJ5P3Drm95CVpDu14Et/qHCARKzGmSJWEdq1m/N07IKm1
tFrHbjPtKb5YSwEg/vzwA/T1DKConA19S6F4Bb+eAQAU2HVGV8ja2c2X4GxxBvB5JbKdd6qvniYbqoZm
ZhvmrljdP4fxxz5MtPZY1+SLV1C8WF9EiULbrStewdckrShUXndyo0hdNQRd7zsXCNwWuNEBJk6BDsAO
Ou92WhH87/IDuB15+Onq6lJesKtrQx4C+Z2uaHVttYV9o6sGBtdDhRa0ZfJYMew1N8ANjcKgLaCYrj22
LbKuYI/DegyLNRsSF/PlId1u090GB/LRg4aOvTiY2JEPOc7dc0nSXUxIiYEukRs536QcdchNuE/SBju9
2T3feGeM63mjwzlZLA2pdazSKAhQ1MSTn8sEf6TODFBidQP7hrgRrxuCLSH3nmDrPCDUekcWOm0cg1aS
9HwdOA/W8fravh7g/7TF3vAqVWmvjYGSQBGTb7UlNb/AYA10qwOvrq0cdIF65cCjVa4FJQWQ8ICs/AtR
KXqwurbOgxYwkCdAO0CLtkdjBgjEWaElyxAYuQ8xiIP6mP2YjEV35PNTuE9/nyfai375U9rHHZSPJ300
OQmnbTDWCwO5dOeS4bmhGP8DzvRti34QRLxpqLqBjwla8C4D4VNK5Tvn4TKamSi7jnyM9oMSAx9CVnqb
QDIV9RQ6ZwOFGSgBihcXF0dHAMW/PW3F4r82B73NfXk/5rNionQ3C/jlk4167/zj9v7z/fbOjp/S/3zP
qc6ujSsF3H9Xi8tV0sk9w3jVfVtfSStLv4eOKh3f9YG2vYlNVWGgIEY9RT5Vzj5jaHBHB7JYX9tPfdVE
UTFUIleNsGpFIWhbr+Quhf4GDPa2amBrdBe+l1CuGrIwe+3JUHRBiZfCyUA78oM8p2ihHKAjX5FlrOkH
TXwPTQg43kcwCVyewhXvM/R+cMbTOeOwMSjaaqslq5Ml4e0i6NdODVPiOLFrVc6yTM6x2dDscQhpAViB
pxq9MhRCnr+TrgtHpNGQ6dLUVrQjI3WPu1ovTciAARBsbwy48gtVPDN+YmmKs3zSlQUPXVy/kvrkfNzL
lvEfovLEvc8EYu93HYklrww6gKffegpManJp5yUO1kdYLO4BdozRYy7O/TK7c7GxrOcNOYZaOmcI7fzl
7Xntzi22RzEXS9AcoogqHVY3WEetWnPTl+vKtRtU2Jr+96A3Of8bumVtebsRejk8YKeL2dx6OwXvMc4e
L9WxNqS0CQCtgtgZwA0ytISWBUZlxKqwshLiRghuy3v0BGRrbYn8U8omhh+pWPQt9gVqOwKaUZuQPkgo
uXe6YkcRf8uZPGgULU4f7NTAkgzRmbmSMiXfNiVJu6lIx8ZVaHSI3DtzaOawfA3JB9pCYOt8i5xF/vty
KTDD4RuJYyFCt9h2seIvL57PXt7NZYuWQkjgfEoqsnSK+4CcST7cnIdg3zhQ6TAGnCbS+tr+4uYQk1ph
jUxqNEJWnfeBfKQ5gp/z1S0OUDVo6/jZGJmxlwZfCddpzmzoCRWW2mgeVte2SxuJcnH78CR7V9IdWVi+
QO0gn53g8iaDFr70gaMpbWvQvH64moG9tvUjtcohFPOCPDjP5nwTof8PUc1h3C0XguLEhH940J2YDGki
BNA2AV07C1gKU6flMY0C/yyAP2Lwb5B32t1wzjn3o38ytTcPDOuT3wjzReEbYU7l0uw70BjCmMQcZg4x
9VCDCoJrCW60VcItnXeloTaDfuS/cZxqW5lekRqhXzo1/LUsLKbIcuE5u/sjAAD//0A+Hc/YEgAA
`,
	},
}

var _escDirs = map[string][]os.FileInfo{}

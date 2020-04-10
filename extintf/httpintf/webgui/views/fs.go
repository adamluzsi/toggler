// Code generated by "esc -o ./views/fs.go -ignore fs.go -pkg views -prefix views ./views"; DO NOT EDIT.

package views

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

	"/doc/show.html": {
		name:    "show.html",
		local:   "views/doc/show.html",
		size:    60,
		modtime: 1565038982,
		compressed: `
H4sIAAAAAAAC/6quTklNy8xLVVDKTczMU6qt5VJQUFCwScksswOzQKC6Wg8mrg+WqK5OzUuprQUEAAD/
/+o8KgA8AAAA
`,
	},

	"/flag/create.html": {
		name:    "create.html",
		local:   "views/flag/create.html",
		size:    2159,
		modtime: 1572749977,
		compressed: `
H4sIAAAAAAAC/8xVT2vbThC951MMc/n9epDk5CwJQkugEEpIybmMtWNpYbW77J/8wfi7l5XkRBV2apdS
4otnNbPvvXmz0m63gjdSM2BPUuNudwEAUAr5CI0i7ytsjA6sA9ZDZsh2V4tk1jEJkD5rWAd2WG+3kH+j
nmG3K4vuarZ3hmyj46ydAS8LVLY2z5lyLQylcfrLLhd7hn0b43qgJkijKyw2ilqEnkNnRIXW+IC/0A7V
r1FGSraaBYKgQNlahidygnX2RKHppG4rPMQ58kpWwnM4nB5KpLYxgKaeK/wxakIIL5Yr7KQQrBEeSUWe
hNYXx6GW9qUBOKOy1plo8biGYbOiNSvYGFfhhilEx1jfjAGkaZXFUPEblLEbKd5AptaS53kK8V2A6Tf2
H/g5nFRuFTXcGSX4TT0kxoE8z/N3mi8LIR//hauWXXoDqGWs741SJga4e312rr0ztLnDbgTOfXAUuH3J
Z3WnG+9IJ9xe6gpXCD09V3i5Wp2EMB3W1Uew3DMLrL8zCzAbCB2DIy1MD9ZzFGa/knrMjd6dO4mB5NAM
EnriPsN4Hfs1uz848+msT7z7rpKuPM/hf2PTd4/Up48wEsGN9NLoaysf7m+x/jKt4fruKzzc357r/gLv
3XdhX3trWtlcW3nGYKJT508ldNIPn6H/PDTRB9PDXgKopGHf9F+a0TqGYPSk2Md1LxcX21QwizPrZE/u
BevPjikw3Chqy2JMHuEqi+PXWlmkK3Nxa4+6DyyncLtlLXa7nwEAAP//kRpWMG8IAAA=
`,
	},

	"/flag/index.html": {
		name:    "index.html",
		local:   "views/flag/index.html",
		size:    817,
		modtime: 1564181141,
		compressed: `
H4sIAAAAAAAC/4SS3YrjMAyF7/sUwrCXiaHXrpeFZWFhGIaZJ1BjNRG4dnHUQsf43YekGTpN/+6CvhOd
I8k5O9pwIFBb5KBKWQAAGMcHaDz2/Uo1MQgFUXYkI+2WM1h1hA64rxoKQknZ9+h93Av889jCC/didLe0
i3MLwbWn7y67faLqVDl/Vl1M/BmDoP/hffp5sJvX0mVhEtpX3JLR0t2mb5SGyNg+0PxphGPorwVGz00H
zRhtlm0d3fFSmTMkDC1BDdPOn05zAs7mDPUwFZRitLjHwukQ9YckFGqP9XliKCVn2CUOAuqXetruJhjh
3t+Ho8CzNQhdos1K6Y3H9je71RDv/18oRVlyLEajNdrzAx99z+h28Ov75AwU3HzjRs8OZPT4AqczGu34
YBc5U3ClfAUAAP//oPGN5jEDAAA=
`,
	},

	"/flag/show.html": {
		name:    "show.html",
		local:   "views/flag/show.html",
		size:    5544,
		modtime: 1584743764,
		compressed: `
H4sIAAAAAAAC/9RYTW+rOBTd51dcefM+JELbdcKoUvqkSNWoSvXWIwffgPWMjWxTpUL895ENtIQmQJLO
x2MTiK+P7XPP9TGUJcMdlwgko1ySqpoBACwYf4FYUGOWJFbSorQk8i2+Nb3rNQYpUgbcBDFKi5pEZQnz
H4Im8z9phlBVizC96wB04PNCY5B00PsBItiqfSB0Aj60aH6C214f32+ndAY0tlzJJQl3giYEMrSpYkuS
K2PJwbA++u0uoIInEtkR3Bqbo2AG7fFmH8JlXliQNMMl+asel4B9zXFJUs4YSgIvVBToxneMnkbqM+SI
1koEiVZFTk5PwXcWdIsCdkovSY7a5YQmSKKNEkIVFp7e/luEPnQErl4VZwdog12aqybCZWGu67Hnxmpq
MXmdnwlVk6ipTJBAxuWS3BDI6H5Jbm9uJiE0xL8ps2Fj/tzO6J0WqKoBihch4y//Ru4MOjE+IzJQO7Ap
gqaSqQxygwVT7ROXdVu9nnNT6ge5MJluAs9T+9cZlEW2RX1VwjbNqMNZcldZht8XJqfyQ9UHGRpDEwy4
FFwiiRahi4u+h80O+J8mnmHMDVfyPuc/N48kWjXPcP+0hp+bx3Nz3MO7tnRbuEeV8Pg+52ekv9BiUnRZ
At+dKtTVwfhPa6iq0eI+0oeUJUo2kO/OlQsaY6oEQ70kNuUGHDVfDMSFsSqDlhIQDr/NE3xVubMiKr5d
s590XcVnxN2e8JW+8ZLoDFzOxlDXqxHMbWGtkg2KKbYZ73luE9C5D3LNM6pfSfRMXxDcOPCkcYcaZYxm
EdZhJxhahKd9eRG6Wv+djwkMBVr8XyRRyVjw+NeSaLSFlhArueM6+/rlXiO8qgJMofGPL9+OpdttYW4h
g6W2QYHU1AIYKJZWDZ8hhn7pNc/9puHjcFmCP5fA/IkLZQ1U1Ww26AvJsaPfh6giuA3uTh0SawkfavZA
0GHuJvNRyEN5fydvslWhdPaQobRzp68xaR2VrZ/pvEaaeJhwl3O2o8NPBaj1HacY/9qq/eRujTc9NPOF
qvIQyBosZK2zRC52/rC3qCUVfnLDxLa+PjuTu8Gan5KTY/nY2/VqAPZgWefDO40O4v/wATX2MPhVnmPQ
jjjM2MbS3VxmU839H6n2sJAGL67531gfny0PEnkmh83m03XR/68sASVrN43phnSxN4wecoYc4qK3n4NP
FW5SD15VkXdSaCUE69XgK0/3C8U7yGz09aYr5NmoUVjc28Gwg5eEegXrFdiUWpCIzIBVsEVojS7cCuXd
Irp81ziX43e/JNHD2z08W2oLM5XjDshkjqeYe8+Pz6Dlwu3jpDVeZykbTLixqMFr4IS7HNs5ukfUdq3N
UeLvAAAA//8AwZenqBUAAA==
`,
	},

	"/layout.html": {
		name:    "layout.html",
		local:   "views/layout.html",
		size:    2293,
		modtime: 1586569394,
		compressed: `
H4sIAAAAAAAC/7yWUWvrNhTH3/spFN3XKVrZHi4jNoyuZYMWRulbKUOxjp3TypKRjpOFkO8+JDuJkzhQ
unCfcmQd//8/nejI2mw0lGiB8UZVwLfbG8YYm020K2jdAFtQbfLuWQyZUbbKOFi+ewhK5zcpTmNCMpCT
qyoDfia74WC+BlKsWCgfgDLeUim+91KHaatqyLiGUHhsCJ3lrHCWwFLGXzrl39ids+SdYWvXevYAiloP
7BkMqACBrZAWrKMIkwsGS4RV4zwN1FeoaZFpWGIBIg1+YmiRUBkRCmUgu53+zIfrMWg/mAeT8UBrA2EB
QJwtPJQZlyoEoCCb1kMRul9Ro50WIQyZPqMR3y+dJaFWEFwNvcZBZCLEK5bMELC/7tn3t/xrjJVHHYSH
0DgbcAnCGS1wFHvyClZj+SZEfoZRUU8RH1wLZYxBiCOOq9UjAhi1di0FGVCDqMG2fTF+VCFGGT5VgZns
+jLFc6fX+U0/oXHJUGe80+WsMCqEjKetWZ1uKPYEtu3biB3hqx71W0TiSTFGj2g/9pqpYnGdA9y98p+q
nre+As+wcPZIPOWERtl8JtPPwVaqIeFuLR1D7xoLNVeepRW14nYX1Frcil9PUaLEsARJ6jinX+9plogV
Rlvt/rRvPH/ZnXpDyr1Ea841DAYaset3yHk+EtQ831dflkZVEq2Gf/mYeCz9IwZiD0ZVIXLNpMH/61d4
UAQXDe/SdLK8gmODxpEs0eqLhn/HFHb38ngFO+2KIJ/vf//j6X5aX7aMaV92G80/6ioZVqqqwMsWLyEw
Ur6Kn9B/5kZFpP4V0WLaf6P6o7Qz2ZqTrpAal8OuS8ORvlNo93z9N3Ss734567vNhqBuTNwnvcq0v3mc
Gh7F3YWABV8cDsn3IFucvgceT4s0v3+zO/j68zDdYzYbsHq7/S8AAP//xUOCrfUIAAA=
`,
	},

	"/login.html": {
		name:    "login.html",
		local:   "views/login.html",
		size:    360,
		modtime: 1560121981,
		compressed: `
H4sIAAAAAAAC/2SQUY7yMAyE33sKy+80F2h6if+/QGgMWJvYUeKyQlXvvmrLSrC8JfONRx4vS6QLCwHm
wILr2gEADJHvMKXQmsdJxUgMx53s9KI1/+IyVzptAkKYjFU8uqRXFoRMdtPosWh7HT8imFJsZO/yjljK
bCAhk0fTLxIEexTyWEJr31ojQklhopumSNXj/83T9z2O3WfYeTZTeQa0+ZzZ8G31p+HlfSqVc6gPHP/x
VYBlcAf4U8F9dhjcdonjP7jI97FbFpK4rj8BAAD///JUCoxoAQAA
`,
	},

	"/pilot/edit.html": {
		name:    "edit.html",
		local:   "views/pilot/edit.html",
		size:    1731,
		modtime: 1584743764,
		compressed: `
H4sIAAAAAAAC/7RUTW+cPBC+768YWe8pEmu9OVbApUmqXKKo7b2axQNYNTa1h00qi/9eYUgEaXejpq1P
u/LwfIxnnhjlBXy4/gyy18axJKUZLuQ47mJUVGtLIDrUVozjDgAgV/oIlcEQClE5y2RZlOkm3baXLy6z
llCBDllFlsmL8n6igffOsncG7tGSeQcxwj5dXD8yeYvm9grGMZftZbl7Bo8RPNqGYH9DyIOnG4NNgEXY
+uS18x1gxdrZQizOaoONDMSZd8a4gQV0xK1ThehdYPETyHwWN/3gKUuoz78yNLqxpAQoZMwOmh/QK7LZ
A3LVatsU4n+xkr8VqMmoQFyeoAXIte0HBv7eUyFarRRZARY7KsSXWbiAI5qBFv1vQkqd2dMj3149o8UI
//3iMf6EYGr9lmH/kQxhSE/4BL87TbAautT/ah6frPFu6M8oSx8bPJCB2vlCkJ0ev0tD+0LFHXaUZi6V
vwIZyFDFoNUGctvUFdVZsATo+mlWnxr00GomowNP0xUj6Bro27IhnxiZYFsyjrMgUjGSVeP4OuN0rpNE
Uq/Lk7O+3zZyMFh9PW9kU/JGIzfOVwRXNG38EflfOhrsnIon/awK3ujmzjGEnipd67/mJJezlHJ3vuww
MDu7rHIYDp1msdm8uUCUAY+Uy/nfmWSQSh9PZaA8HYK5nBK2XEc/WTVF/W4Fu3T1RwAAAP//XMJYBcMG
AAA=
`,
	},

	"/pilot/find.html": {
		name:    "find.html",
		local:   "views/pilot/find.html",
		size:    460,
		modtime: 1563994919,
		compressed: `
H4sIAAAAAAAC/2yQz27CMAzG7zyF5SNSG+3ecgJNSDtw2Au4jQFLqROlLn9U9d0nCpsY283JL/l++TKO
bgnvm09wSUI0txf1sHTTtBhHz3tRBuxIFKdpAQBQeTlBG6jva2yjGqvhaiYz3cfcAbUmUWt8ikTo2I7R
15hib/idkIbMxXznZyooyEHZI3gyKhqxM2XPWpzJ2qPooca3J+PdKhx8z/Z7e0aiaTCwa+IajS+GoNRx
jfPTSr7Ydo2QArV8jMFzrnF3I7AbmiAtbC7GWSnAdl2WJf4jaAazqA9DPzSdvNR7HHiai5Slo3zF1Uck
D7Owcnf20sz9rVa52zfd15XzclotxpHVT9NXAAAA//+9jaN8zAEAAA==
`,
	},

	"/doc": {
		name:  "doc",
		local: `views/doc`,
		isDir: true,
	},

	"/flag": {
		name:  "flag",
		local: `views/flag`,
		isDir: true,
	},

	"/pilot": {
		name:  "pilot",
		local: `views/pilot`,
		isDir: true,
	},

	"/views": {
		name:  "views",
		local: `./views`,
		isDir: true,
	},
}

var _escDirs = map[string][]os.FileInfo{

	"views/doc": {
		_escData["/doc/show.html"],
	},

	"views/flag": {
		_escData["/flag/create.html"],
		_escData["/flag/index.html"],
		_escData["/flag/show.html"],
	},

	"views/pilot": {
		_escData["/pilot/edit.html"],
		_escData["/pilot/find.html"],
	},

	"./views": {
		_escData["/doc"],
		_escData["/flag"],
		_escData["/layout.html"],
		_escData["/login.html"],
		_escData["/pilot"],
	},
}

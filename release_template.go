package bindata

const tmplImport string = `
import (
	"os"
	"time"
)

`

// used on web servers like, headers should be included:
// "Vary": "Accept-Encoding"
// "Content-Encoding": "gzip"
const tmplImportCompressIrisMode = `
import (
	"fmt"
	"os"
	"strings"
	"time"
)

`

const tmplReleaseHeader = `
type gzipAsset struct {
	bytes []byte
	info  gzipFileInfoEx
}

type gzipFileInfoEx interface {
	os.FileInfo
	MD5Checksum() string
}

type gzipBindataFileInfo struct {
	name        string
	size        int64
	mode        os.FileMode
	modTime     time.Time
	md5checksum string
}

func (fi gzipBindataFileInfo) Name() string {
	return fi.name
}
func (fi gzipBindataFileInfo) Size() int64 {
	return fi.size
}
func (fi gzipBindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi gzipBindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi gzipBindataFileInfo) MD5Checksum() string {
	return fi.md5checksum
}
func (fi gzipBindataFileInfo) IsDir() bool {
	return false
}
func (fi gzipBindataFileInfo) Sys() interface{} {
	return nil
}

`

// keep error there in order to be compatible with the existing API.
const tmplReleaseCommon string = `

func %s() (*gzipAsset, error) {
	bytes := _%s
	info := gzipBindataFileInfo{
		name: %q,
		size: %d,
		md5checksum: %q,
		mode: os.FileMode(%d),
		modTime: time.Unix(%d, 0),
	}

	a := &gzipAsset{bytes: bytes, info: info}

	return a, nil
}

`

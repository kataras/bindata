// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

package bindata

// nolint: gas
import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// writeOneFileRelease writes the release code file for each file (when splited file).
func writeOneFileRelease(w io.Writer, c *Config, a *Asset) (err error) {
	_, err = fmt.Fprint(w, tmplImport)
	if err != nil {
		return
	}

	return writeReleaseAsset(w, c, a)
}

// writeRelease writes the release code file for single file.
func writeRelease(w io.Writer, c *Config, toc []Asset) (err error) {
	err = writeReleaseHeader(w, c)
	if err != nil {
		return
	}

	for i := range toc {
		err = writeReleaseAsset(w, c, &toc[i])
		if err != nil {
			return
		}
	}

	return
}

// writeReleaseHeader writes output file headers.
// This targets release builds.
func writeReleaseHeader(w io.Writer, c *Config) (err error) {
	if _, err = fmt.Fprint(w, tmplImportCompressIrisMode); err != nil {
		return
	}

	_, err = fmt.Fprint(w, tmplReleaseHeader)
	return
}

// writeReleaseAsset write a release entry for the given asset.
// A release entry is a function which embeds and returns
// the file's byte content.
func writeReleaseAsset(w io.Writer, c *Config, asset *Asset) (err error) {
	fd, err := os.Open(asset.Path)
	if err != nil {
		return
	}

	err = compress(w, asset, fd)
	if err != nil {
		_ = fd.Close()
		return
	}

	err = fd.Close()
	if err != nil {
		return
	}

	return assetReleaseCommon(w, c, asset)
}

// sanitize prepares a valid UTF-8 string as a raw string constant.
// Based on https://code.google.com/p/go/source/browse/godoc/static/makestatic.go?repo=tools
func sanitize(b []byte) []byte {
	// Replace ` with `+"`"+`
	b = bytes.Replace(b, []byte("`"), []byte("`+\"`\"+`"), -1)

	// Replace BOM with `+"\xEF\xBB\xBF"+`
	// (A BOM is valid UTF-8 but not permitted in Go source files.
	// I wouldn't bother handling this, but for some insane reason
	// jquery.js has a BOM somewhere in the middle.)
	return bytes.Replace(b, []byte("\xEF\xBB\xBF"), []byte("`+\"\\xEF\\xBB\\xBF\"+`"), -1)
}

func compress(w io.Writer, asset *Asset, r io.Reader) (err error) {
	_, err = fmt.Fprintf(w, "var _%s = []byte(\n\t\"", asset.Func)
	if err != nil {
		return err
	}

	gz := gzip.NewWriter(&StringWriter{Writer: w})
	_, err = io.Copy(gz, r)
	if err != nil {
		_ = gz.Close()
		return err
	}

	err = gz.Close()
	if err != nil {
		return
	}

	_, err = fmt.Fprint(w, `")`)

	return
}

// nolint: gas
func assetReleaseCommon(w io.Writer, c *Config, asset *Asset) (err error) {
	fi, err := os.Stat(asset.Path)
	if err != nil {
		return
	}

	mode := uint(fi.Mode())
	modTime := fi.ModTime().Unix()
	size := fi.Size()
	if c.Mode > 0 {
		mode = uint(os.ModePerm) & c.Mode
	}
	if c.ModTime > 0 {
		modTime = c.ModTime
	}

	var md5checksum string
	if c.MD5Checksum {
		var buf []byte

		buf, err = ioutil.ReadFile(asset.Path)
		if err != nil {
			return
		}

		h := md5.New()
		if _, err = h.Write(buf); err != nil {
			return
		}
		md5checksum = fmt.Sprintf("%x", h.Sum(nil))
	}

	_, err = fmt.Fprintf(w, tmplReleaseCommon, asset.Func, asset.Func,
		asset.Name, size, md5checksum, mode, modTime)

	return
}

// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

package bindata

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

const (
	// DefPackageName define default package name.
	DefPackageName = "main"

	// DefOutputName define default generated file name.
	DefOutputName = "bindata_gzip.go"
)

// List of errors.
var (
	ErrNoInput       = errors.New("No input")
	ErrNoPackageName = errors.New("Missing package name")
	ErrCWD           = errors.New("Unable to determine current working directory")
)

// InputConfig defines options on an asset directory to be convert.
type InputConfig struct {
	// Path defines a directory containing asset files to be included
	// in the generated output.
	Path string

	// Recusive defines whether subdirectories of Path
	// should be recursively included in the conversion.
	Recursive bool
}

// Config defines a set of options for the asset conversion.
type Config struct {
	// cwd contains current working directory.
	cwd string

	// Name of the package to use. Defaults to 'main'.
	Package string

	// Tags specify a set of optional build tags, which should be
	// included in the generated output. The tags are appended to a
	// `// +build` line in the beginning of the output file
	// and must follow the build tags syntax specified by the go tool.
	Tags string

	// Input defines the directory path, containing all asset files as
	// well as whether to recursively process assets in any sub directories.
	Input []InputConfig

	// Output defines the output file for the generated code.
	// If left empty, this defaults to 'bindata_gzip.go' in the current
	// working directory and the current directory in case of having true
	// to `Split` config.
	Output string

	// Prefix defines a regular expression which should used to strip
	// substrings from all file names when generating the keys in the table of
	// contents.  For example, running without the `-prefix` flag, we get:
	//
	// 	$ bindata /path/to/templates
	// 	go_bindata["/path/to/templates/foo.html"] = _path_to_templates_foo_html
	//
	// Running with the `-prefix` flag, we get:
	//
	//	$ bindata -prefix "/a/path/to/some/" /a/path/to/some/templates/
	//	_bindata["templates/foo.html"] = templates_foo_html
	Prefix string

	// Ignores any filenames matching the regex pattern specified, e.g.
	// path/to/file.ext will ignore only that file, or \\.gitignore
	// will match any .gitignore file.
	//
	// This parameter can be provided multiple times.
	Ignore []*regexp.Regexp

	// Include contains list of regex to filter input files.
	Include []*regexp.Regexp

	// When nonzero, use this as mode for all files.
	Mode uint

	// When nonzero, use this as unix timestamp for all files.
	ModTime int64

	// Perform a debug build. This generates an asset file, which
	// loads the asset contents directly from disk at their original
	// location, instead of embedding the contents in the code.
	//
	// This is mostly useful if you anticipate that the assets are
	// going to change during your development cycle. You will always
	// want your code to access the latest version of the asset.
	// Only in release mode, will the assets actually be embedded
	// in the code. The default behaviour is Release mode.
	Debug bool

	// Perform a dev build, which is nearly identical to the debug option. The
	// only difference is that instead of absolute file paths in generated code,
	// it expects a variable, `rootDir`, to be set in the generated code's
	// package (the author needs to do this manually), which it then prepends to
	// an asset's name to construct the file path on disk.
	//
	// This is mainly so you can push the generated code file to a shared
	// repository.
	Dev bool

	// Split the output into several files. Every embedded file is bound into
	// a specific file, and a common file is also generated containing API and
	// other common parts.
	// If true, the output config is a directory and not a file.
	Split bool

	// MD5Checksum is a flag that, when set to true, indicates to calculate
	// MD5 checksums for files.
	MD5Checksum bool
}

// NewConfig returns a default configuration struct.
func NewConfig() *Config {
	c := new(Config)
	c.Package = DefPackageName
	c.Output = DefOutputName
	c.Ignore = make([]*regexp.Regexp, 0)
	c.Include = make([]*regexp.Regexp, 0)
	return c
}

func (c *Config) validateInput() (err error) {
	for _, input := range c.Input {
		_, err = os.Lstat(input.Path)
		if err != nil {
			return fmt.Errorf("Failed to stat input path '%s': %v",
				input.Path, err)
		}
	}
	return
}

// validateOutput will check if output is valid.
//
// (1) If output is empty, set the output directory to,
// (1.1) current working directory if `split` option is used, or
// (1.2) current working directory with default output file output name.
// (2) If output is not empty, check the directory and file write status.
func (c *Config) validateOutput() (err error) {
	// (1)
	if len(c.Output) == 0 {
		if c.Split {
			// (1.1)
			c.Output = c.cwd
		} else {
			// (1.2)
			c.Output = filepath.Join(c.cwd, DefOutputName)
		}

		return
	}

	// (2)
	dir, file := filepath.Split(c.Output)

	if dir != "" {
		err = os.MkdirAll(dir, 0700)
		if err != nil {
			return fmt.Errorf("Create output directory: %v", err)
		}
	}

	if len(file) == 0 {
		if !c.Split {
			c.Output = filepath.Join(dir, DefOutputName)
		}
	}

	if c.Split {
		return
	}

	var fout *os.File

	fout, err = os.Create(c.Output)
	if err != nil {
		return
	}

	err = fout.Close()
	if err != nil {
		return
	}

	return
}

// validate ensures the config has sane values.
// Part of which means checking if certain file/directory paths exist.
func (c *Config) validate() (err error) {
	if len(c.Package) == 0 {
		return ErrNoPackageName
	}

	err = c.validateInput()
	if err != nil {
		return
	}

	err = c.validateOutput()
	if err != nil {
		return
	}

	return
}

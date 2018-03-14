// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

/*
Package bindata converts any file into manageable Go source code. Useful for
embedding binary data into a go program. The file data is gzip
compressed before and after converted to a raw byte slice.

The following paragraphs cover some of the customization options
which can be specified in the Config struct, which must be passed into
the Translate() call.


Debug vs Release builds

When used with the `Debug` option, the generated code does not actually include
the asset data. Instead, it generates function stubs which load the data from
the original file on disk. The asset API remains identical between debug and
release builds, so your code will not have to change.

This is useful during development when you expect the assets to change often.
The host application using these assets uses the same API in both cases and
will not have to care where the actual data comes from.

An example is a Go webserver with some embedded, static web content like
HTML, JS and CSS files. While developing it, you do not want to rebuild the
whole server and restart it every time you make a change to a bit of
javascript. You just want to build and launch the server once. Then just press
refresh in the browser to see those changes. Embedding the assets with the
`debug` flag allows you to do just that. When you are finished developing and
ready for deployment, just re-invoke `bindata` without the `-debug` flag.
It will now embed the latest version of the assets.


Path prefix stripping

The keys used in the `_bindata` map are the same as the input file name passed
to `bindata`. This includes the path. In most cases, this is not desirable,
as it puts potentially sensitive information in your code base.  For this
purpose, the tool supplies another command line flag `-prefix`.  This accepts a
[regular expression](https://github.com/google/re2/wiki/Syntax) string, which
will be used to match a portion of the map keys and function names that should
be stripped out.

For example, running without the `-prefix` flag, we get:

	$ bindata /path/to/templates/

	_bindata["/path/to/templates/foo.html"] = path_to_templates_foo_html

Running with the `-prefix` flag, we get:

	$ bindata -prefix "/.*\/some/" /a/path/to/some/templates/

	_bindata["templates/foo.html"] = templates_foo_html


Build tags

With the optional Tags field, you can specify any go build tags that
must be fulfilled for the output file to be included in a build. This
is useful when including binary data in multiple formats, where the desired
format is specified at build time with the appropriate tags.

The tags are appended to a `// +build` line in the beginning of the output file
and must follow the build tags syntax specified by the go tool.


Splitting generated file

When you want to embed big files or plenty of files, then the generated output
is really big (maybe over 3Mo). Even if the generated file shouldn't be read,
you probably need use analysis tool or an editor which can become slower
with a such file.

Generating big files can be avoided with `-split` command line option.
In that case, the given output is a directory path, the tool will generate
one source file per file to embed, and it will generate a common file
nammed `common.go` which contains commons parts like API.

*/
package bindata

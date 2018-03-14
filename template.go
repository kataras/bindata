package bindata

const tmplTypeBintree string = `
type gzipBintree struct {
	Func     func() (*gzipAsset, error)
	Children map[string]*gzipBintree
}

var _gzipbintree = &gzipBintree`

const tmplBinTreeValues string = `{Func: %s, Children: map[string]*gzipBintree{`

const tmplFuncAssetDir string = `

// GzipAssetDir returns the file names below a certain
// directory embedded in the file by bindata.
// For example if you run bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then GzipAssetDir("data") would return []string{"foo.txt", "img"}
// GzipAssetDir("data/img") would return []string{"a.png", "b.png"}
// GzipAssetDir("foo.txt") and GzipAssetDir("notexist") would return an error
// GzipAssetDir("") will return []string{"data"}.
func GzipAssetDir(name string) ([]string, error) {
	node := _gzipbintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, &os.PathError{
					Op: "open",
					Path: name,
					Err: os.ErrNotExist,
				}
			}
		}
	}
	if node.Func != nil {
		return nil, &os.PathError{
			Op: "open",
			Path: name,
			Err: os.ErrNotExist,
		}
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

`

const tmplFuncAsset string = `

// GzipAsset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func GzipAsset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _gzipbindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("GzipAsset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, &os.PathError{Op: "open", Path: name, Err: os.ErrNotExist}
}

// MustGzipAsset is like GzipAsset but panics when GzipAsset would return an error.
// It simplifies safe initialization of global variables.
// nolint: deadcode
func MustGzipAsset(name string) []byte {
	a, err := GzipAsset(name)
	if err != nil {
		panic("asset: GzipAsset(" + name + "): " + err.Error())
	}

	return a
}

// GzipAssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or could not be loaded.
func GzipAssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _gzipbindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("GzipAssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, &os.PathError{Op: "open", Path: name, Err: os.ErrNotExist}
}

// GzipAssetNames returns the names of the assets.
// nolint: deadcode
func GzipAssetNames() []string {
	names := make([]string, 0, len(_gzipbindata))
	for name := range _gzipbindata {
		names = append(names, name)
	}
	return names
}

//
// _gzipbindata is a table, holding each asset generator, mapped to its name.
//
var _gzipbindata = map[string]func() (*gzipAsset, error){
`

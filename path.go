package filever

import (
	"path/filepath"
	"strings"
)

// Values are empty when not applicable. See README for naming.
type PathParts struct {
	FullPath string //  Relative or absolute. E.g. `e/app~fv=4mIbJJPq.min.js`.
	Path     string `json:",omitempty"` // Just the path, relative or absolute.  E.g. `e/`.
	File     string `json:",omitempty"` // Just the file.  E.g. `app~fv=4mIbJJPq.min.js`.
	Base     string `json:",omitempty"` // No path, extension, or version. E.g. `app`.
	Ext      string `json:",omitempty"` // Extension.  Includes "sub" extensions.  E.g. `min.js`.
	BaseExt  string `json:",omitempty"` // Base Extension.  E.g. `.js`

	// FileVer specific.
	FileVer  string `json:",omitempty"` // E.g. `app~fv=4mIbJJPq.min.js`.
	BarePath string `json:",omitempty"` // E.g. `e/app.min.js`
	BareFile string `json:",omitempty"` // No path or version.  E.g. `app.min.js`.
	Bare     string `json:",omitempty"` // No path, version, or extension.  E.g. `app`.
	DelimVer string `json:",omitempty"` // E.g. `~fv=4mIbJJPq`.
	Version  string `json:",omitempty"` // E.g. `4mIbJJPq`.
}

// Populate populates PathParts from FullPath.  Populates FileVer specific
// fields only if FileVer exists.
func (p *PathParts) Populate() {
	p.Path, p.File = PathCut(p.FullPath)

	// strings.Cut splits on first instance of char but excludes first "." in ext.
	var found bool
	p.Base, p.Ext, found = strings.Cut(p.File, ".")
	if found {
		p.Ext = "." + p.Ext // add back "." 😞
	}

	li := strings.LastIndex(p.File, ".")
	if li > 0 {
		p.BaseExt = string(p.File[li:])
	}

	// FileVer specific
	p.DelimVer = VerAnySizeRegexC.FindString(p.Base)

	_, p.Version, found = strings.Cut(p.DelimVer, Delim)
	if found {
		p.FileVer = p.File
	}
	p.BareFile = VerAnySizeRegexC.ReplaceAllString(p.File, "")
	p.BarePath = p.Path + p.BareFile
	p.Bare = VerAnySizeRegexC.ReplaceAllString(p.Base, "")
}

func Populated(fullPath string) *PathParts {
	p := &PathParts{FullPath: fullPath}
	p.Populate()
	return p
}

// PathCut, from the full path, returns directory path and the file name.  e.g.
// for `..subdir/test_5~fv=wwjNHrIw.js` returns `..subdir/` and
// `test_5~fv=wwjNHrIw.js`
func PathCut(path string) (dir, base string) {
	// Alternatively use `string(os.PathSeparator)`
	base = filepath.Base(path)
	dir = path[:len(path)-len(base)]
	return
}

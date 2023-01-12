package filever

import (
	"path/filepath"
	"strings"
)

// Values are empty when not applicable. See README for naming.
type PathParts struct {
	Full    string //  The full path, relative or absolute. E.g. `e/app~fv=4mIbJJPq.min.js`.
	Dir     string `json:",omitempty"` // Just the directory, relative or absolute, without the file.  E.g. `e/`.
	File    string `json:",omitempty"` // Just the file, no directory.  E.g. `app~fv=4mIbJJPq.min.js`.
	Base    string `json:",omitempty"` // No dir, extension, or version. E.g. `app`.
	Ext     string `json:",omitempty"` // Extension.  Includes "sub" extensions.  E.g. `min.js`.
	BaseExt string `json:",omitempty"` // Base Extension.  E.g. `.js`

	// FileVer specific.
	FileVer  string `json:",omitempty"` // E.g. `app~fv=4mIbJJPq.min.js`.
	DelimVer string `json:",omitempty"` // E.g. `~fv=4mIbJJPq`.
	Version  string `json:",omitempty"` // E.g. `4mIbJJPq`.
	BarePath string `json:",omitempty"` // E.g. `e/app.min.js`
	BareFile string `json:",omitempty"` // No dir or version.  E.g. `app.min.js`.
	Bare     string `json:",omitempty"` // No dir, version, or extension.  E.g. `app`.
}

// Populate populates PathParts from FullPath.  Populates FileVer only if
// FileVer exists, but other FileVer specific fields will populate regardless
// (such as Bare).
func (p *PathParts) Populate() {
	p.Dir, p.File = PathCut(p.Full)

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
	p.BarePath = p.Dir + p.BareFile
	p.Bare = VerAnySizeRegexC.ReplaceAllString(p.Base, "")
}

func Populated(fullPath string) *PathParts {
	p := &PathParts{Full: fullPath}
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

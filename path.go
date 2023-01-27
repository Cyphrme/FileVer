package filever

import (
	"net/url"
	"os"
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

	// Example: https://example.com:8081/bob/joe.txt?name=ferret#nose?name=bob
	Scheme        string `json:",omitempty"` // e.g. `https`
	Authority     string `json:",omitempty"` // e.g. `example.com:8081`
	Host          string `json:",omitempty"` // e.g. `example.com`
	Port          string `json:",omitempty"` // e.g. `:8081`
	URIPath       string `json:",omitempty"` // e.g. `/bob/joe.txt`
	Query         string `json:",omitempty"` // e.g. `name=ferret`
	Fragment      string `json:",omitempty"` // e.g. `nose?name=bob`
	Anchor        string `json:",omitempty"` // e.g. `nose`
	FragmentQuery string `json:",omitempty"` // e.g. `?name=bob`
	Quag          string `json:",omitempty"` // e.g. `?name=ferret#nose?name=bob`
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

	// URL.  Will no populate unless scheme is present.
	u, err := url.Parse(p.Full)
	if err != nil || u.Scheme == "" {
		return
	}
	p.Scheme = u.Scheme
	p.Authority = u.Host
	p.Host, p.Port, found = strings.Cut(u.Host, ":")
	if found {
		p.Port = ":" + p.Port // add back ":" 😞
	}
	// Do not include the path separator if the path is only the path separator.
	if u.Path != "/" && u.Path != string(os.PathSeparator) {
		p.URIPath = u.Path
	}
	p.Query = u.RawQuery
	p.Fragment = u.Fragment // Note: Fragment is not sent by browsers on requests.
	p.Anchor, p.FragmentQuery, found = strings.Cut(p.Fragment, "?")
	if p.Query != "" {
		p.Quag += "?" + p.Query
	}
	if p.Fragment != "" {
		p.Quag += "#" + p.Fragment
	}
}

func Populated(fullPath string) *PathParts {
	p := &PathParts{Full: fullPath}
	p.Populate()
	return p
}

// PathCut, from the full path, returns directory path and the file name.  e.g.
// for `..subdir/test_5~fv=wwjNHrIw.js` returns `..subdir/` and
// `test_5~fv=wwjNHrIw.js`
//
// If there is an ending separator, e.g. "/" in "https://localhost:8081/", it
// will be removed.
func PathCut(path string) (dir, base string) {
	// Remove ending separator if present.
	lastchar := string(path[len(path)-1])
	if lastchar == "/" || lastchar == string(os.PathSeparator) {
		path = path[:len(path)-1]
	}

	scheme, rest, found := strings.Cut(path, "://")
	if found {
		d, b := PathCut(rest)
		return scheme + "://" + d, b
	}

	li := strings.LastIndex(path, "/")
	if li > 0 {
		dir = path[:li+1]
		base = path[li+1:]
		return
	}
	return "", path
}

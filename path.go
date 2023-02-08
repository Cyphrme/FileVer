package filever

import (
	"net/url"
	"os"
	"strings"
)

// Values are empty when not applicable. See README for naming.
type PathParts struct {
	Full     string `json:"full"`                //  The full path, relative or absolute. E.g. `e/app~fv=4mIbJJPq.min.js`.
	Dir      string `json:"dir,omitempty"`       // Just the directory, relative or absolute, without the file.  E.g. `e/`.
	File     string `json:"file,omitempty"`      // Just the file, no directory.  E.g. `app~fv=4mIbJJPq.min.js`.
	FileBase string `json:"file_base,omitempty"` // No dir, extension, or version. E.g. `app`.
	Ext      string `json:"ext,omitempty"`       // Extension.  Includes "sub" extensions.  E.g. `min.js`.
	ExtBase  string `json:"ext_base,omitempty"`  // Base Extension.  E.g. `.js`

	// FileVer specific.
	FileVer  string `json:"filever,omitempty"`   // E.g. `app~fv=4mIbJJPq.min.js`.
	DelimVer string `json:"delim_ver,omitempty"` // E.g. `~fv=4mIbJJPq`.
	Version  string `json:"version,omitempty"`   // E.g. `4mIbJJPq`.
	BarePath string `json:"bare_path,omitempty"` // E.g. `e/app.min.js`
	BareFile string `json:"bare_file,omitempty"` // No dir or version.  E.g. `app.min.js`.
	Bare     string `json:"bare,omitempty"`      // No dir, version, or extension.  E.g. `app`.

	// Example: https://example.com:8081/bob/joe.txt?name=ferret#nose?name=bob
	// Origin // TODO
	// PostOrigin //TODO
	// URIBase   // TODO
	Scheme        string `json:"scheme,omitempty"`         // e.g. `https`
	Authority     string `json:"authority,omitempty"`      // e.g. `example.com:8081`
	Host          string `json:"host,omitempty"`           // e.g. `example.com`
	Port          string `json:"port,omitempty"`           // e.g. `:8081`
	URIPath       string `json:"uri_path,omitempty"`       // e.g. `/bob/joe.txt` // TODO think about putting this up above.
	Query         string `json:"query,omitempty"`          // e.g. `name=ferret`
	Fragment      string `json:"fragment,omitempty"`       // e.g. `nose?name=bob`
	Anchor        string `json:"anchor,omitempty"`         // e.g. `nose`
	FragmentQuery string `json:"fragment_query,omitempty"` // e.g. `?name=bob`
	Quag          string `json:"quag,omitempty"`           // e.g. `?name=ferret#nose?name=bob`
}

// Populate populates PathParts from FullPath.  Populates FileVer only if
// FileVer exists, but other FileVer specific fields will populate regardless
// (such as Bare).
func (p *PathParts) Populate() {
	p.Dir, p.File = PathCut(p.Full)

	// strings.Cut splits on first instance of char but excludes first "." in ext.
	var found bool
	p.FileBase, p.Ext, found = strings.Cut(p.File, ".")
	if found {
		p.Ext = "." + p.Ext // add back "." 😞
	}

	li := strings.LastIndex(p.File, ".")
	if li > 0 {
		p.ExtBase = string(p.File[li:])
	}

	// FileVer specific
	p.DelimVer = VerAnySizeRegexC.FindString(p.FileBase)
	_, p.Version, found = strings.Cut(p.DelimVer, Delim)
	if found {
		p.FileVer = p.File
	}
	p.BareFile = VerAnySizeRegexC.ReplaceAllString(p.File, "")
	p.BarePath = p.Dir + p.BareFile
	p.Bare = VerAnySizeRegexC.ReplaceAllString(p.FileBase, "")

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
	if path == "" {
		return "", ""
	}

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
	if li >= 0 {
		dir = path[:li+1]
		base = path[li+1:]
		return
	}
	return "", path
}

package filever

import (
	"strings"

	path "github.com/cyphrme/path"
)

// PathParts holds path's parts.  Values are empty when not applicable. See
// README for naming.
type PathParts struct {
	path.PathParts

	// FileVer specific.
	FileVer  string `json:"filever,omitempty"`   // E.g. `app~fv=4mIbJJPq.min.js`.
	DelimVer string `json:"delim_ver,omitempty"` // E.g. `~fv=4mIbJJPq`.
	Version  string `json:"version,omitempty"`   // E.g. `4mIbJJPq`.
	BarePath string `json:"bare_path,omitempty"` // E.g. `e/app.min.js`
	BareFile string `json:"bare_file,omitempty"` // No dir or version.  E.g. `app.min.js`.
	Bare     string `json:"bare,omitempty"`      // No dir, version, or extension.  E.g. `app`.
}

// Populate populates PathParts from FullPath.  Populates FileVer only if
// FileVer exists, but other FileVer specific fields will populate regardless
// (such as Bare).
func (p *PathParts) Populate() {
	var found bool

	p.PathParts.Populate()

	// FileVer specific
	p.DelimVer = VerAnySizeRegexC.FindString(p.FileBase)
	_, p.Version, found = strings.Cut(p.DelimVer, Delim)
	if found {
		p.FileVer = p.File
	}
	p.BareFile = VerAnySizeRegexC.ReplaceAllString(p.File, "")
	p.BarePath = p.Dir + p.BareFile
	p.Bare = VerAnySizeRegexC.ReplaceAllString(p.FileBase, "")
}

func Populated(fullPath string) *PathParts {
	p := new(PathParts)
	p.Full = fullPath
	p.Populate()
	return p
}

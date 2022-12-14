// Package filever (FileVer) is designed for automatic file versioning and
// distribution.
package filever

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/cyphrme/coze"
	"golang.org/x/exp/slices"
)

// VersionSize is the number of digest characters in the Version.
var VersionSize = 8

// Delim is the FileVer Delimiter.  FileVers are end delimited by any non-b64u
// character, such as "." or another "~".
var Delim = "~fv="

// VerRegex is the version string regex e.g. `\~fv=00000000`. The version is
// appended to the base file name. VerRegex is set here to a default value, but
// may be changed as desired.
var VerRegex = regexp.QuoteMeta(Delim) + `[0-9A-Za-z_-]{` + fmt.Sprint(VersionSize) + `}`

// Compiled VerRegex.
var VerRegexC *regexp.Regexp

// VerAnySizeRegex is the version string regex for any size version, e.g.
// (~fv=000). This is especially useful for cleaning out versions that may be of
// a different size. By default this regex will match multiple versions in a
// single file name for sanitization, e.g. this regex will match both version in
// this string:  `test.txt~fv=000~fv=JPq`
var VerAnySizeRegex = regexp.QuoteMeta(Delim) + `[0-9A-Za-z_-]*`

// Compiled VerAnySizeRegex
var VerAnySizeRegexC *regexp.Regexp

// FileVerPathReg "midVer" format and should match any path in source files that
// includes a fileVer. The path regex should be delimited by non-path
// characters, like `"`.  This regex including "startPath", e.g. `../`.
//
// Currently, this regex is very simple and doesn't support all valid paths.
var FileVerPathReg = `[0-9A-Za-z_\-\/]*~fv=[0-9A-Za-z_\-]*.[0-9A-Za-z_\-.]*`

// HashAlg is the hash alg used for versioning.
var HashAlg = coze.SHA256

// Config holds the settings for operating the main FileVer functions.
//
//	Src      - Source directory starting point.  Must be nil for SrcFiles.
//	SrcFiles - Manually provided src files.  Will be set to Src's files if nil (default behavior).
//	SrcReg   - Compiled Regex used to search source files for FileVersions for Replace().
//	             May be set by external program.
//	Dist     - destination directory.  Default: Output will be on one level.
type Config struct {
	Src      string
	SrcFiles []string
	SrcReg   *regexp.Regexp
	Dist     string
	UseSAVR  bool

	// Use Internally
	Info *Info
}

type Info struct {
	// PV "Path:Version"  which is [relative path + baseName: "New" Version]. e.g.
	//`"test_1.js":"4WYoW0MN"``
	PV map[string]string

	// SAVR "Search All Versioned, Regex". Regex that matches all filevers at once
	// in the current directory. This results in a large regex, but allows single
	// pass searching for each file.  The downside is that a match matches all
	// "versions", so the match still needs to be matched to a single version
	// after this regex returns a match.
	// Example of the resulting regex:
	// `(app\.min\.js\~fv=[0-9A-Za-z_-]*)|(coze\.min\.js\~fv=[0-9A-Za-z_-]*)`
	// SAVR needs to be regenerated for each subdirectory.
	// key: current directory with root being "" and no preceding `/`.  E.g. `subdir/subdir2`
	// Value: current directory SAVR.
	SAVR string

	// Versioned files (relatively pathed to c.Dist).  Generated when calling
	// `Version()`.  Used by `Replace()` to replace in source files the correct FileVer.
	VersionedFiles []string

	// Index is built by Index().  Key is a versioned file (e.g. subdir/test_3.js)
	// value is a list of source files that refer to that versioned file, e.g.
	//["subdir/test_4.js","test_1.min.js","test_2.js"]
	Index map[string][]string

	// Total number of references that were replaced after running `Replace()`.
	TotalSourceReplaces int

	// CheckedFilePaths are files that were checked by `Replace()`, but the
	// existing source file was current.  Paths are
	// relative to pwd.
	CheckedFilePaths []string

	// UpdatedFilePaths are files that were updated by `Replace()`.  Paths are
	// relative to pwd.
	UpdatedFilePaths []string

	// Used for processing
	// Current path should be relative to dist or src, not including dist, and not
	// from pwd.
	CurrentPath    string `json:"-"`
	CurrentMatches int    `json:"-"`
}

func init() {
	VerRegexC = regexp.MustCompile(VerRegex)
	VerAnySizeRegexC = regexp.MustCompile(VerAnySizeRegex)
}

// VersionReplace see notes on Version() and Replace()
func VersionReplace(c *Config) (err error) {
	err = Version(c)
	if err != nil {
		return err
	}
	return Replace(c)
}

// Version versions all dummied FileVer files in input directory including
// subdirectories, copies them into c.dist, and removes any existing versions.
//
// Populates c.Info.PV and c.Info.VersionedFiles.
func Version(c *Config) (err error) {
	c.Info = new(Info)
	c.Info.PV = map[string]string{}

	if c.SrcFiles == nil {
		c.SrcFiles, err = ExistingVersionedFiles(c.Src)
		if err != nil {
			return err
		}
	}

	c.Info.VersionedFiles = []string{} // Files without paths.
	for _, path := range c.SrcFiles {
		file, err := FileToFileVerOutputDelete(path, c)
		if err != nil {
			return err
		}
		// Set the Info struct
		basePath := VerAnySizeRegexC.ReplaceAllString(file, "")                            // get bare file name no version with path
		c.Info.PV[basePath] = strings.TrimPrefix(VerAnySizeRegexC.FindString(file), Delim) // Set version, e.g. "4WYoW0MN"
		c.Info.VersionedFiles = append(c.Info.VersionedFiles, file)
	}
	return nil
}

// Replace updates all source file references to versioned files with the
// current version. c.Info.PV and c.Info.VersionedFiles must be set correctly.
func Replace(c *Config) (err error) {
	//fmt.Printf("\nReplace Config  %+v Info: %+v\n", c, c.Info)
	if c.Info == nil {
		return fmt.Errorf("c.Info must be set.")
	}

	genSrcReg(c)

	// PathedVersionedReplace is called on each match.  Input is the matched string.
	// Variable "in" is file name, without the relative subdirectory path.
	var PathedVersionedReplace = func(in []byte) []byte {
		//fmt.Printf("PathedVersionedReplace - match: %s\n", in)
		c.Info.CurrentMatches++
		// Match will include version.  Get bare file name without versioning.
		bare := VerAnySizeRegexC.ReplaceAllString(string(in), "")
		// If the file has a start path, e.g. `../` in
		// `../test_1~fv=SgfqvMD3.min.js` is the startPath.
		// To find the version in PV, must remove startPath (characters [".","/"]).
		startPathReg := regexp.MustCompile(`^[\/\.]*`)
		startPath := startPathReg.FindString(bare)
		woStartPath := startPathReg.ReplaceAllString(bare, "")
		version := c.Info.PV[woStartPath]
		//fmt.Printf("bare: %s version: %s woStartPath: %s\n", bare, version, woStartPath)
		return []byte(startPath + genFileVer(woStartPath, version, c))
	}

	// Walk walks all files (recursively) in directory. Variable `path` is
	// relative to to running location of the program (program root dir).
	var walk = func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		//fmt.Printf("\nWalking path: %s\n", path)

		if d.IsDir() {
			return nil
		}
		c.Info.CurrentMatches = 0

		//fmt.Printf("walk - path: %s; d: %+v, c.Info %+v\n", path, d, c.Info)
		read, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		replaced := c.SrcReg.ReplaceAllFunc(read, PathedVersionedReplace)
		//fmt.Printf("Replaced contents: %s\n", replaced)
		if c.Info.CurrentMatches > 0 { // Only Write out on match.

			if slices.Equal(read, replaced) { // Don't write out if there are no updates.
				c.Info.CheckedFilePaths = append(c.Info.CheckedFilePaths, path)
				return nil
			}

			c.Info.TotalSourceReplaces += c.Info.CurrentMatches
			c.Info.UpdatedFilePaths = append(c.Info.UpdatedFilePaths, path)
			//fmt.Printf("info.CurrentMatches: %d.  Writing updated file: %s\n", c.Info.CurrentMatches, path)
			err = os.WriteFile(path, replaced, 0)
			if err != nil {
				return err
			}
		}
		return nil
	}

	err = filepath.WalkDir(c.Dist, walk)
	return err
}

// Index builds an index of what files The index map, has the key of the version
// that being replaced, and the value of the files that the version exists.
func Index(c *Config) {
	// TODO
}

func IndexReplace(c *Config) {
	// Index
}

// genFileVer generates the fileVer (e.g. app.min.js~fv=0000 or
// app~fv=0000.min.js) from the bare relative file name (`app.min.js` or
// `subdir/app.min.js`) and digest (0000...). Input `digest` may be the
// shortened be the version (0000) instead of the whole 64 character digest.
func genFileVer(file, digest string, c *Config) string {
	// On no digest, generate dummy, e.g. `app.min.js~fv=0000`
	if digest == "" {
		fmt.Printf("***WARNING*** Digest empty for %s\n", file)
		digest = strings.Repeat("0", VersionSize)
	}
	// strings.Cut splits on first instance of char.  Resulting `ext` will be missing first "."
	baseWithoutExt, ext, _ := strings.Cut(file, ".")
	//fmt.Printf("genFileVer:%s, %s \n", baseWithoutExt, ext)
	fv := baseWithoutExt + Delim + digest[:VersionSize] + "." + ext
	return fv

}

// genFileVerRegex Returns the regex to find the versioned file encapsulated in
// parentheses, e.g. `(subdir/test_3\~fv=[0-9A-Za-z_-]*.js)`.  Variable `file`
// should be the bare relative file path.  E.g. `subdir/test_3.js` or
// `test_1.js`.
func genFileVerRegex(file string, c *Config) (regex string) {
	dir, base := PathParts(file)
	// strings.Cut splits on first instance.  Resulting excludes match.
	baseWithoutExt, ext, _ := strings.Cut(base, ".")
	return "(" + regexp.QuoteMeta(dir+baseWithoutExt) + VerAnySizeRegex + "." + ext + ")"
}

// CleanVersionFiles removes any versioned files recursively in the given path.
func CleanVersionFiles(path string) error {
	// Walk walks all files (recursively) in directory.
	// Variable path is relative to to running location of the program (program root dir).
	var walk = func(path string, d fs.DirEntry, err error) error {
		//fmt.Printf("Clean Walk - path: %s; d: %+v\n", path, d)
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		if VerAnySizeRegexC.Match([]byte(path)) {
			//fmt.Printf("Matched: %s; removing\n", path)
			err := os.RemoveAll(path)
			if err != nil {
				panic(err)
			}
		}

		return nil
	}

	return filepath.WalkDir(path, walk)
}

// ExistingVersionedFiles returns all existing versioned files in  `directory`
// including dummies (e.g. "app.min.js~fv=00000000") and subdirectories.
// `directory` should not have trailing slash.  Returns paths relative to
// `directory`.
func ExistingVersionedFiles(directory string) (fileVers []string, err error) {
	var walk = func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			i := VerRegexC.FindStringIndex(d.Name())
			if i == nil {
				return nil
			}
			// Remove directory from path.  Returned path should be relative to `directory`
			path = strings.Replace(path, directory+string(os.PathSeparator), "", 1)

			fileVers = append(fileVers, path)
		}
		return nil
	}

	return fileVers, filepath.WalkDir(directory, walk)
}

// DirFiles returns all files in a directory including subdirectories.  Relative
// is to include relative path in file name.
func DirFiles(directory string, relative bool) (files []string, err error) {
	var walk = func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			if !relative {
				files = append(files, d.Name())
			} else {
				files = append(files, path)
			}
		}
		return nil
	}

	return files, filepath.WalkDir(directory, walk)
}

// ListFilesInPath will return a sorted of file names on the first level in a
// given directory. It is not recursive. If path is a file, it will return the
// name of that file.
func ListFilesInPath(path string) (files []string, err error) {
	// os.Stat first because os.ReadDir will error on file.
	fi, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	// Single file
	if !fi.IsDir() {
		files = append(files, fi.Name())
		return files, nil
	}

	// Is a dir
	dirEntries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	files = make([]string, 0, len(dirEntries))
	for _, d := range dirEntries {
		if !d.IsDir() {
			files = append(files, d.Name())
		}
	}

	return files, nil
}

// Returns directory and base.  e.g. for `..subdir/test_5~fv=wwjNHrIw.js`
// returns `..subdir/` and `test_5~fv=wwjNHrIw.js`
func PathParts(path string) (dir, base string) {
	base = filepath.Base(path)
	dir = path[:len(path)-len(base)]
	return
}

// FileToFileVerOutputDelete accepts a bare or dummied versioned file, copies
// the file renamed with its correct FileVer to an output directory, and deletes
// any previous version in the output directory. Returns outFilePath, relative
// to c.Dist, which itself is relative to `pwd`.  This function ignores
// subdirectories, and does not re-hash files already in output directory.
//
// `filePath` is input file path relative to `c.Src`, which itself is
// relative to pwd (unless absolute).
//
// c.Dist and c.Src must be set. If pwd == c.Src, it may be left blank.
func FileToFileVerOutputDelete(filePath string, c *Config) (outFilePath string, err error) {
	//log.Debug().Msgf("FileToFileVerOutputDelete filePath: %s src: %s dist: %s ", filePath, c.Src, c.Dist)
	dig, _, err := HashFile(c.Src+string(os.PathSeparator)+filePath, HashAlg)
	if err != nil {
		return "", err
	}

	base := VerAnySizeRegexC.ReplaceAllString(filepath.Base(filePath), "")
	fileVer := genFileVer(base, dig.String(), c)

	// If file exists, remove any existing version (including dummy version and/or
	// duplicate versions in a single file name) from file name.
	//
	// Relative path from `src` (excluding src, base name). Starts with "/".
	rDir := filepath.Dir(filePath)
	if rDir == "." { // Sanitize filepath, which adds a "." on empty. ????
		rDir = ""
	}
	rDir = strings.Replace(rDir, c.Src, "", 1)
	distRDir := c.Dist + string(os.PathSeparator) + rDir
	outFilePath = rDir + string(os.PathSeparator) + fileVer

	if rDir == "" {
		outFilePath = fileVer
	} else {
		//log.Debug().Msgf("creating relative dirs: %s", distRDir)
		os.MkdirAll(distRDir, 0755)
		if err != nil {
			return "", err
		}
		outFilePath = rDir + string(os.PathSeparator) + fileVer
	}
	//log.Debug().Msgf("rDir: %s, outFilePath: %+v", rDir, outFilePath)

	// Check if the FileVer already exists in output directory, if it does, don't
	// copy.  Regardless, also check for existing and/or previous versions in
	// output directory.
	files, err := ListFilesInPath(distRDir)
	if err != nil {
		return "", err
	}

	// Search for outFilePath FileVer, e.g. `app.min.js~fv=00000000` or
	// `app~fv=00000000.min.js.map`. Must match whole file name since if just
	// matching substring files with alternative extensions, e.g. `min.js` vs
	// `min.js.map`, will match when they shouldn't.
	anyVersionReg := regexp.MustCompile(genFileVerRegex(base, c))
	matchedExisting := false
	// Even though there should only ever be one existing versioned output, this
	// for loop continues even after match, to search for and remove other errant
	// versioned files.
	for _, f := range files {
		if !matchedExisting && f == fileVer { // File is the current version.
			matchedExisting = true
			continue // Continue in case of other errant copies.
		}

		matched := anyVersionReg.Match([]byte(f))
		if !matched {
			continue
		}

		// Existing file is a different version than new file.  Delete Existing.
		del := c.Dist + string(os.PathSeparator) + rDir + string(os.PathSeparator) + f
		//log.Debug().Msgf("Delete matched: %s file: %s", escapedAnyVersion, del)
		err := os.Remove(del)
		if err != nil {
			return "", err
		}
		// Continue in case of other errant copies.
	}

	if matchedExisting { // Match with current FileVer found.  Don't re-copy.
		return outFilePath, nil
	}

	// Copy into output directory.
	in := c.Src + string(os.PathSeparator) + filePath
	//log.Debug().Msgf("in: %s", in)
	input, err := os.ReadFile(in)
	if err != nil {
		return "", err
	}

	o := c.Dist + string(os.PathSeparator) + outFilePath
	//log.Debug().Msgf("Writing copy to: %s", o)
	err = os.WriteFile(o, input, 0644)
	if err != nil {
		return "", err
	}

	return outFilePath, nil
}

func genSrcReg(c *Config) {
	if c.SrcReg != nil { // Don't recompile if set.
		return
	}
	if c.UseSAVR {
		genSAVR(c)
		c.SrcReg = regexp.MustCompile(c.Info.SAVR)
		return
	}

	c.SrcReg = regexp.MustCompile(FileVerPathReg)
}

// Generate SAVR.  e.g.
// (subdir/test_3\\~fv=[0-9A-Za-z_-]*.js)|(test_1\\~fv=[0-9A-Za-z_-]*.js)
func genSAVR(c *Config) {
	c.Info.SAVR = ""
	bookEnd := false
	for _, k := range c.Info.VersionedFiles {
		nv := VerAnySizeRegexC.ReplaceAllString(k, "") // get bare file name without version.
		nv = strings.ReplaceAll(nv, c.Dist, "")        // Remove dist (imports are relative to dist).

		if bookEnd { // Fencepost
			c.Info.SAVR += "|"
		} else {
			bookEnd = true
		}

		c.Info.SAVR += genFileVerRegex(nv, c)
	}
}

// HashFile accepts a path, a hashing algorithm, return digest and pointer to file.
func HashFile(path string, alg coze.HshAlg) (digest coze.B64, file *[]byte, err error) {
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	d, err := coze.Hash(alg, fileBytes)
	if err != nil {
		return nil, nil, err
	}
	//log.Debug().Msgf("FileDigest: %X", d)
	return coze.B64(d), &fileBytes, nil
}

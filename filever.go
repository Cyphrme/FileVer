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
	"github.com/rs/zerolog/log"
)

// VersionSize is the number of digest characters in the Version.
var VersionSize = 8

// Delim is the FileVer Delimiter.  FileVers are end delimited by any non-b64u
// character, such as "." or another "?".
var Delim = "?fv="

// VerRegex is the version string regex e.g. `\?fv=00000000`. The version is
// appended to the base file name. VerRegex is set here to a default value, but
// may be changed as desired.
var VerRegex = regexp.QuoteMeta(Delim) + `[0-9A-Za-z_-]{` + fmt.Sprint(VersionSize) + `}`

// Compiled VerRegex.
var VerRegexC *regexp.Regexp

// VerAnySizeRegex is the version string regex for any size version, e.g.
// (?fv=000). This is especially useful for cleaning out versions that may be of
// a different size. By default this regex will match multiple versions in a
// single file name for sanitization, e.g. this regex will match both version in
// this string:  `test.txt?fv=000?fv=JPq`
var VerAnySizeRegex = regexp.QuoteMeta(Delim) + `[0-9A-Za-z_-]*`

// Compiled VerAnySizeRegex
var VerAnySizeRegexC *regexp.Regexp

// HashAlg is the hash alg used for versioning.
var HashAlg = coze.SHA256

// Config holds the settings for operating the main FileVer functions.
//
//	Src   - source directory starting point.
//	Dist  - destination directory.  Default: Output will be on one level.
//	MidVer - Mid File version format, e.g. test?fv=000.txt instead of test.txt?fv=000
type Config struct {
	Src    string
	Dist   string
	MidVer bool

	// Use Internally
	Info *Info
}

type Info struct {
	// Total number of replaces after running `Replace()`
	TotalSourceReplaces int

	// UpdatedFilePaths are files that were updated by `Replace()`.  Paths are
	// relative to pwd.
	UpdatedFilePaths []string

	// SAVR "Search All Versioned, Regex". Create the regex to match a file with
	// all versions at once. This results in a large regex, but files do not need
	// to be re-searched for each FileVer.  The downside is that a match matches
	// all "versions", so the match still needs to be matched to a single version
	// after this regex returns a match.
	// Example of the resulting regex:
	// `(app\.min\.js\?fv=[0-9A-Za-z_-]*)|(coze\.min\.js\?fv=[0-9A-Za-z_-]*)`
	SAVR string
	//PV "Path:Version"  which is [dist path baseName: "New" Version]. e.g.
	//`"test/dummy/dist/test_1.js":"4WYoW0MN"``
	PV map[string]string

	// Versioned files (relatively pathed to c.Dist).  Generated when calling
	// `Version()`.  Used by `Replace()` to know what files are versioned.
	VersionedFiles []string

	// Used for processing
	CurrentPath    string `json:"-"` // Relative to dist or src, not from program root dir.
	CurrentVersion string `json:"-"` // The Version, e.g. `4WYoW0MN`
	CurrentBare    string `json:"-"` // File name with no pathing or version.
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

	scrFilesPaths, err := ExistingVersionedFiles(c.Src)
	if err != nil {
		return err
	}

	c.Info.VersionedFiles = []string{} // Files without paths.
	for _, path := range scrFilesPaths {
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
	fmt.Printf("\nReplace Config  %+v Info: %+v", c, c.Info)
	if c.Info == nil {
		return fmt.Errorf("c.Info must be set.")
	}
	c.Info.SAVR = ""

	bookEnd := false
	for _, k := range c.Info.VersionedFiles {
		nv := VerAnySizeRegexC.ReplaceAllString(k, "") // get bare file name without version.
		nv = strings.ReplaceAll(nv, c.Dist, "")        // Remove dist (imports are relative to dist).
		escapedBase := regexp.QuoteMeta(nv)
		if bookEnd { // Fencepost
			c.Info.SAVR += "|"
		} else {
			bookEnd = true
		}
		// TODO probably needs to be modified to support mid version.
		c.Info.SAVR += "(" + escapedBase + VerAnySizeRegex + ")"
	}
	reg := regexp.MustCompile(c.Info.SAVR)

	// 20221229 OLD
	// Implements regexp.ReplaceFunc
	// var PathedVersionedReplace = func(in []byte) []byte {
	// 	c.Info.CurrentMatches++
	// 	bare := VerAnySizeRegexC.ReplaceAllString(string(in), "") // get bare file name.
	// 	path := filepath.Dir(c.Info.CurrentPath) + string(os.PathSeparator) + bare
	// 	version := c.Info.PV[path]
	// 	fmt.Printf("\n\nPathedVersionedReplace - match: %s, c.Info.CurrentPath: %s; bare: %s; path: %s; New Version: %s\n", in, c.Info.CurrentPath, bare, path, version)
	// 	return []byte(bare + version)
	// }

	// "in" will only be file name, without the relative subdirectory path.
	var PathedVersionedReplace = func(in []byte) []byte {
		fmt.Printf("PathedVersionedReplace - match: %s\n", in)
		c.Info.CurrentMatches++
		return []byte(genFileVer(c.Info.CurrentBare, c.Info.CurrentVersion, c))
	}

	// Walk walks all files (recursively) in directory.
	// Variable path is relative to to running location of the program (program root dir).
	var walk = func(path string, d fs.DirEntry, err error) error {
		fmt.Printf("Walk - path: %s; d: %+v\n", path, d)
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		// Current path should be relative to dist, not including dist.
		c.Info.CurrentPath = strings.TrimPrefix(path, c.Dist+"/")
		c.Info.CurrentBare = VerAnySizeRegexC.ReplaceAllString(c.Info.CurrentPath, "") // get bare file name, without pathing.
		c.Info.CurrentVersion = c.Info.PV[c.Info.CurrentBare]

		fmt.Printf("Walk - c.Info %+v\n", c.Info)

		read, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		// Set matches to 0 for each file
		c.Info.CurrentMatches = 0
		replaced := reg.ReplaceAllFunc(read, PathedVersionedReplace)
		fmt.Printf("Replaced contents: %s\n", replaced)
		if c.Info.CurrentMatches > 0 { // Only write out on match.
			c.Info.TotalSourceReplaces += c.Info.CurrentMatches
			c.Info.UpdatedFilePaths = append(c.Info.UpdatedFilePaths, path)
			log.Debug().Msgf("info.CurrentMatches: %d.  Writing updated file: %s", c.Info.CurrentMatches, path)
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

// genFileVer generates the fileVer (e.g.  app.min.js?fv=00000000 or
// app?fv=00000000.min.js) from the bare file name (app.min.js) and digest (00000000...).
// Input `digest` may be the shortened be the version (00000000) instead of the
// whole 64 character digest.
func genFileVer(bareFileName, digest string, c *Config) string {
	if !c.MidVer {
		// Not mid version, e.g. app.min.js?fv=00000000
		return bareFileName + Delim + digest[:VersionSize]
	} else {
		// Mid version, e.g. app?fv=00000000.min.js
		ext := filepath.Ext(bareFileName)
		baseWithoutExt := bareFileName[:len(bareFileName)-len(ext)]
		return baseWithoutExt + Delim + digest[:VersionSize] + ext
	}
}

// ExistingVersionedFiles returns all existing versioned files in  `directory`
// including dummies (e.g. "app.min.js?fv=00000000") and subdirectories.
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
	if rDir == "." { // Sanitize filepath, which adds a "." on empty. 😞
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

	// Search for outFilePath FileVer, `file name + delim + any version size`, e.g.
	// `app.min.js?fv=00000000` or `app.min.js.map?fv=00000000`. Must match whole
	// file name since if just matching substring files with alternative
	// extensions, e.g. `min.js` vs `min.js.map`, will match when they shouldn't.
	escapedAnyVersion := regexp.QuoteMeta(base) + VerAnySizeRegex
	anyVersionReg := regexp.MustCompile(escapedAnyVersion)
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

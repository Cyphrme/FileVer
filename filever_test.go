package filever

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/cyphrme/watch"
)

// Each test that writes files has its own directory, so test outputs can be
// pushed to git and used in explanations.  Init calls clean before executing
// tests.
var dummySrc = "test/dummy/src"
var dummyDist = "test/dummy/dist"
var dummyEndSrc = "test/dummy_end/src"
var dummyEndDist = "test/dummy_end/dist"
var dummyNoSrc = "test/dummy_no/src"
var dummyNoDist = "test/dummy_no/dist"
var watchSrc = "test/watch/src"
var watchDist = "test/watch/dist"
var cleanDist = "test/clean" // For ExampleCleanVersionFiles. Uses dummySrc as src.

func init() {
	clean()
}

// Example VersionReplace with mid version with "dummy" input files.
// go test -run '^ExampleVersionReplace$'
func ExampleVersionReplace() {
	c := &Config{Src: dummySrc, Dist: dummyDist}
	err := VersionReplace(c)
	if err != nil {
		panic(err)
	}
	PrintPretty(c)
	PrintFile(dummyDist + "/" + c.Info.VersionedFiles[0])

	// Output:
	// {
	// 	"Src": "test/dummy/src",
	// 	"SrcFiles": [
	// 		"subdir/test_3?fv=00000000.js",
	// 		"subdir/test_4?fv=00000000.js",
	// 		"test_1?fv=00000000.js",
	// 		"test_2?fv=00000000.js"
	// 	],
	// 	"Dist": "test/dummy/dist",
	// 	"EndVer": false,
	// 	"Info": {
	// 		"PV": {
	// 			"subdir/test_3.js": "7gnMxWXQ",
	// 			"subdir/test_4.js": "CeWrTnIH",
	// 			"test_1.js": "oG9WcWOW",
	// 			"test_2.js": "4UZ0_4xw"
	// 		},
	// 		"SAVR": "(subdir/test_3\\?fv=[0-9A-Za-z_-]*.js)|(subdir/test_4\\?fv=[0-9A-Za-z_-]*.js)|(test_1\\?fv=[0-9A-Za-z_-]*.js)|(test_2\\?fv=[0-9A-Za-z_-]*.js)",
	// 		"VersionedFiles": [
	// 			"subdir/test_3?fv=7gnMxWXQ.js",
	// 			"subdir/test_4?fv=CeWrTnIH.js",
	// 			"test_1?fv=oG9WcWOW.js",
	// 			"test_2?fv=4UZ0_4xw.js"
	// 		],
	// 		"TotalSourceReplaces": 14,
	// 		"UpdatedFilePaths": [
	// 			"test/dummy/dist/subdir/test_3?fv=7gnMxWXQ.js",
	// 			"test/dummy/dist/subdir/test_4?fv=CeWrTnIH.js",
	// 			"test/dummy/dist/test_1?fv=oG9WcWOW.js",
	// 			"test/dummy/dist/test_2?fv=4UZ0_4xw.js"
	// 		]
	// 	}
	// }
	// File test/dummy/dist/subdir/test_3?fv=7gnMxWXQ.js:
	// ////////////////
	// import * as test1 from '../test_1?fv=oG9WcWOW.js';
	// import * as test2 from '../test_2?fv=4UZ0_4xw.js';
	// import * as test4 from '../subdir/test_4?fv=CeWrTnIH.js';
	// ////////////////
}

func ExamplePathParts() {
	fmt.Println(PathParts("..subdir/test_5?fv=wwjNHrIw.js"))
	// Output:
	// ..subdir/ test_5?fv=wwjNHrIw.js
}

// Example VersionReplace using end ver format and dummy file inputs.
func ExampleVersionReplace_end() {
	c := &Config{Src: dummyEndSrc, Dist: dummyEndDist, EndVer: true}
	err := VersionReplace(c)
	if err != nil {
		panic(err)
	}
	// PrintPretty(c)
	PrintFile(dummyEndDist + "/" + c.Info.VersionedFiles[0])

	// Output:
	// File test/dummy_end/dist/subdir/test_3.js?fv=gGkfoWxG:
	// ////////////////
	// import * as test1 from '../test_1.js?fv=Jmp9dlP7';
	// import * as test2 from '../test_2.js?fv=ZbopKA8M';
	// ////////////////
}

// Example_watchVersionAndReplace demonstrates using FileVer with the external
// program "watch". Uses the "mid version" format.
func Example_watchVersionAndReplace() {
	// Set up and test with Watch. Normally (outside of testing) watch must call
	// filever.  For testing, filever will call watch so that `go test` works.
	// Also see notes in `watch_src.sh`
	watch.ParseFlags()
	watch.FC.Daemon = false
	watch.FC.ConfigPath = "test/watch.json5"
	watch.Run()

	// Normal FileVer setup.
	c := &Config{Src: watchSrc, Dist: watchDist}
	err := VersionReplace(c)
	if err != nil {
		panic(err)
	}
	// PrintPretty(c)
	PrintFile(watchDist + "/" + c.Info.VersionedFiles[0])

	// Output:
	// Flag `daemon` set to false.  Running commands in config and exiting.
	// File test/watch/dist/subdir/test_3?fv=fGF8m_Po.js:
	// ////////////////
	// import * as test1 from '../test_1?fv=pfQjFHV-.min.js';
	// import * as test2 from '../test_2?fv=7XyeFLlY.js';
	// import * as test4 from '../subdir/test_4?fv=NszgzyIB.js';
	// // Comments referring to './test_1?fv=pfQjFHV-.min.js' should be updated as well.
	// ////////////////
}

// Example_noDummy demonstrates inputting "manually" enumerated files to be
// processes by FileVer, aka it does not use dummy file inputs. Does not do
// Replace().
func Example_noDummy() {
	c := &Config{
		Src:  dummyNoSrc,
		Dist: dummyNoDist,
		SrcFiles: []string{
			"test_1.js",
			"test_2.js",
			"subdir/test_3.js",
			"subdir/test_4.js",
		},
	}
	Version(c)
	// TODO Replace does not appear to be working.
	fmt.Println(c.Info.VersionedFiles)
	// PrintPretty(c.Info)
	// PrintFile(dummyNoDist + "/" + c.Info.VersionedFiles[0])

	// Output:
	// [test_1?fv=DneDuFnM.js test_2?fv=62jVtATD.js subdir/test_3?fv=7gnMxWXQ.js subdir/test_4?fv=WMNO8lQ4.js]
}

func ExampleListFilesInPath() {
	f, err := ListFilesInPath(dummyNoSrc)
	if err != nil {
		panic(err)
	}

	fmt.Println(f)

	// Output:
	// [test_1.js test_2.js unversioned_file.md]
}

func ExampleExistingVersionedFiles() {
	// Mid Version format
	files, err := ExistingVersionedFiles(dummySrc)
	if err != nil {
		panic(err)
	}
	fmt.Println(files)

	// End version format.
	files, err = ExistingVersionedFiles(dummyEndSrc)
	if err != nil {
		panic(err)
	}
	fmt.Println(files)

	// Output:
	// [subdir/test_3?fv=00000000.js subdir/test_4?fv=00000000.js test_1?fv=00000000.js test_2?fv=00000000.js]
	// [subdir/test_3?fv=00000000.js subdir/test_4?fv=00000000.js test_1.js?fv=00000000 test_2.js?fv=00000000]
}

func ExampleCleanVersionFiles() {
	// Generate test files.
	c := &Config{Src: dummySrc, Dist: cleanDist}
	err := VersionReplace(c)
	if err != nil {
		panic(err)
	}
	// Clean out generate test files.
	CleanVersionFiles(cleanDist)
	f, err := ListFilesInPath(cleanDist)
	if err != nil {
		panic(err)
	}

	fmt.Println(c.Info.VersionedFiles)
	fmt.Println(f)
	// Output:
	// [subdir/test_3?fv=7gnMxWXQ.js subdir/test_4?fv=CeWrTnIH.js test_1?fv=oG9WcWOW.js test_2?fv=4UZ0_4xw.js]
	// [not_versioned_example.txt]
}

func Test_clean(t *testing.T) {
	clean()
}

// Clean removes versioned files from dist.
func clean() {
	c := []string{
		dummyDist,
		dummyEndDist,
		dummyNoDist,
		watchDist,
		cleanDist,
	}

	for _, v := range c {
		CleanVersionFiles(v)
	}
}

func Test__nuke(t *testing.T) {
	nukeAndRebuildTestDirs()
}

// Completely delete all test dirs and recreate.
func nukeAndRebuildTestDirs() {
	c := []string{
		dummyDist,
		dummyEndDist,
		dummyNoDist,
		watchDist,
		cleanDist,
	}

	for _, v := range c {
		// Recreate the dir via remove and create.
		err := os.RemoveAll(v)
		if err != nil {
			panic(err)
		}

		err = os.Mkdir(v, 0777)
		if err != nil {
			panic(err)
		}

		d1 := []byte("This example file exists in `dist` directory and is not versioned.")
		err = os.WriteFile(v+"/not_versioned_example.txt", d1, 0644)
		if err != nil {
			panic(err)
		}
	}
}

// PrintFile is a helper function.
func PrintFile(filePath string) {

	// Print out the first file and verify contents
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = file.Close(); err != nil {
			panic(err)
		}
	}()

	b, err := ioutil.ReadAll(file)

	fmt.Printf("File %s:\n////////////////\n%s\n////////////////\n", filePath, b)

}

func PrintPretty(s any) {
	json, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(json))
}

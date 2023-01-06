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
var dummyNoSrc = "test/dummy_no/src"
var dummyNoDist = "test/dummy_no/dist"
var watchSrc = "test/watch/src"
var watchDist = "test/watch/dist"
var cleanDist = "test/clean" // For ExampleCleanVersionFiles. Uses dummySrc as src.

func init() {
	clean()
}

func ExampleFileVerPathReg() {
	ts :=
		`// test_1?fv=00000000.js
		import * as test2 from './test_2?fv=00000000.js';
		import * as test3 from './subdir/test_3?fv=00000000.js';
		import * as test4 from './subdir/test_4?fv=00000000.js';
`

	c := &Config{Src: dummySrc, Dist: dummyDist}
	genSrcReg(c)

	matches := c.SrcReg.FindAllString(ts, -1)
	fmt.Println(matches)
	// Output:
	//[test_1?fv=00000000.js /test_2?fv=00000000.js /subdir/test_3?fv=00000000.js /subdir/test_4?fv=00000000.js]
}

func ExamplePathParts() {
	fmt.Println(PathParts("..subdir/test_5?fv=wwjNHrIw.js"))
	// Output:
	// ..subdir/ test_5?fv=wwjNHrIw.js
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
	// ***WARNING*** Digest empty for test_3.js
	// {
	// 	"Src": "test/dummy/src",
	// 	"SrcFiles": [
	// 		"subdir/test_3?fv=00000000.js",
	// 		"subdir/test_4?fv=00000000.js",
	// 		"test_1?fv=00000000.js",
	// 		"test_2?fv=00000000.js"
	// 	],
	// 	"SrcReg": {},
	// 	"Dist": "test/dummy/dist",
	// 	"UseSAVR": false,
	// 	"Info": {
	// 		"PV": {
	// 			"subdir/test_3.js": "7gnMxWXQ",
	// 			"subdir/test_4.js": "CeWrTnIH",
	// 			"test_1.js": "oG9WcWOW",
	// 			"test_2.js": "4UZ0_4xw"
	// 		},
	// 		"SAVR": "",
	// 		"VersionedFiles": [
	// 			"subdir/test_3?fv=7gnMxWXQ.js",
	// 			"subdir/test_4?fv=CeWrTnIH.js",
	// 			"test_1?fv=oG9WcWOW.js",
	// 			"test_2?fv=4UZ0_4xw.js"
	// 		],
	// 		"TotalSourceReplaces": 15,
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
	PrintPretty(c)
	PrintFile(watchDist + "/" + c.Info.VersionedFiles[0])

	// Output:
	// Flag `daemon` set to false.  Running commands in config and exiting.
	// ***WARNING*** Digest empty for test_1.js
	// {
	// 	"Src": "test/watch/src",
	// 	"SrcFiles": [
	// 		"subdir/test_3?fv=00000000.js",
	// 		"subdir/test_4?fv=00000000.js",
	// 		"test_1?fv=00000000.min.js",
	// 		"test_2?fv=00000000.js"
	// 	],
	// 	"SrcReg": {},
	// 	"Dist": "test/watch/dist",
	// 	"UseSAVR": false,
	// 	"Info": {
	// 		"PV": {
	// 			"subdir/test_3.js": "fGF8m_Po",
	// 			"subdir/test_4.js": "NszgzyIB",
	// 			"test_1.min.js": "SgfqvMD3",
	// 			"test_2.js": "7XyeFLlY"
	// 		},
	// 		"SAVR": "",
	// 		"VersionedFiles": [
	// 			"subdir/test_3?fv=fGF8m_Po.js",
	// 			"subdir/test_4?fv=NszgzyIB.js",
	// 			"test_1?fv=SgfqvMD3.min.js",
	// 			"test_2?fv=7XyeFLlY.js"
	// 		],
	// 		"TotalSourceReplaces": 20,
	// 		"UpdatedFilePaths": [
	// 			"test/watch/dist/subdir/test_3?fv=fGF8m_Po.js",
	// 			"test/watch/dist/subdir/test_4?fv=NszgzyIB.js",
	// 			"test/watch/dist/test_1.min.js.map",
	// 			"test/watch/dist/test_1?fv=SgfqvMD3.min.js",
	// 			"test/watch/dist/test_2?fv=7XyeFLlY.js"
	// 		]
	// 	}
	// }
	// File test/watch/dist/subdir/test_3?fv=fGF8m_Po.js:
	// ////////////////
	// import * as test1 from '../test_1?fv=SgfqvMD3.min.js';
	// import * as test2 from '../test_2?fv=7XyeFLlY.js';
	// import * as test4 from '../subdir/test_4?fv=NszgzyIB.js';
	// // Comments referring to './test_1?fv=SgfqvMD3.min.js' should be updated as well.
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

	// Output:
	// [subdir/test_3?fv=00000000.js subdir/test_4?fv=00000000.js test_1?fv=00000000.js test_2?fv=00000000.js]
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
	// ***WARNING*** Digest empty for test_3.js
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

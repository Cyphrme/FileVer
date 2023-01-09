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
		`// test_1~fv=00000000.js
		import * as test2 from './test_2~fv=00000000.js';
		import * as test3 from './subdir/test_3~fv=00000000.js';
		import * as test4 from './subdir/test_4~fv=00000000.js';
`

	c := &Config{Src: dummySrc, Dist: dummyDist}
	genSrcReg(c)

	matches := c.SrcReg.FindAllString(ts, -1)
	fmt.Println(matches)
	// Output:
	//[test_1~fv=00000000.js /test_2~fv=00000000.js /subdir/test_3~fv=00000000.js /subdir/test_4~fv=00000000.js]
}

func ExamplePathParts() {
	fmt.Println(PathParts("..subdir/test_5~fv=wwjNHrIw.js"))
	// Output:
	// ..subdir/ test_5~fv=wwjNHrIw.js
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
	// 		"subdir/test_3~fv=00000000.js",
	// 		"subdir/test_4~fv=00000000.js",
	// 		"test_1~fv=00000000.js",
	// 		"test_2~fv=00000000.js"
	// 	],
	// 	"SrcReg": {},
	// 	"Dist": "test/dummy/dist",
	// 	"UseSAVR": false,
	// 	"Info": {
	// 		"PV": {
	// 			"subdir/test_3.js": "_X83uO__",
	// 			"subdir/test_4.js": "GJIrg6k1",
	// 			"test_1.js": "vPCb4GVO",
	// 			"test_2.js": "BOl7h9TM"
	// 		},
	// 		"SAVR": "",
	// 		"VersionedFiles": [
	// 			"subdir/test_3~fv=_X83uO__.js",
	// 			"subdir/test_4~fv=GJIrg6k1.js",
	// 			"test_1~fv=vPCb4GVO.js",
	// 			"test_2~fv=BOl7h9TM.js"
	// 		],
	// 		"Index": null,
	// 		"TotalSourceReplaces": 15,
	// 		"UpdatedFilePaths": [
	// 			"test/dummy/dist/subdir/test_3~fv=_X83uO__.js",
	// 			"test/dummy/dist/subdir/test_4~fv=GJIrg6k1.js",
	// 			"test/dummy/dist/test_1~fv=vPCb4GVO.js",
	// 			"test/dummy/dist/test_2~fv=BOl7h9TM.js"
	// 		]
	// 	}
	// }
	// File test/dummy/dist/subdir/test_3~fv=_X83uO__.js:
	// ////////////////
	// import * as test1 from '../test_1~fv=vPCb4GVO.js';
	// import * as test2 from '../test_2~fv=BOl7h9TM.js';
	// import * as test4 from '../subdir/test_4~fv=GJIrg6k1.js';
	// ////////////////

}

// // TestVersionReplace VersionReplace with mid version with "dummy" input files.
// func TestVersionReplace(t *testing.T) {
// 	c := &Config{Src: dummySrc, Dist: dummyDist}
// 	err := VersionReplace(c)
// 	if err != nil {
// 		panic(err)
// 	}
// 	PrintPretty(c)
// 	PrintFile(dummyDist + "/" + c.Info.VersionedFiles[0])
// }

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

	// Replace Config  &{Src:test/watch/src SrcFiles:[subdir/test_3~fv=00000000.js subdir/test_4~fv=00000000.js test_1~fv=00000000.min.js test_2~fv=00000000.js] SrcReg:<nil> Dist:test/watch/dist UseSAVR:false Info:0xc0001dc5b0} Info: &{PV:map[subdir/test_3.js:gia0-_Z_ subdir/test_4.js:da1EKBXZ test_1.min.js:1sTEzePc test_2.js:qBbNrrTr] SAVR: VersionedFiles:[subdir/test_3~fv=gia0-_Z_.js subdir/test_4~fv=da1EKBXZ.js test_1~fv=1sTEzePc.min.js test_2~fv=qBbNrrTr.js] Index:map[] TotalSourceReplaces:0 UpdatedFilePaths:[] CurrentPath: CurrentMatches:0}
	// ***WARNING*** Digest empty for test_3.js
	// ***WARNING*** Digest empty for test_1.js
	// {
	// 	"Src": "test/watch/src",
	// 	"SrcFiles": [
	// 		"subdir/test_3~fv=00000000.js",
	// 		"subdir/test_4~fv=00000000.js",
	// 		"test_1~fv=00000000.min.js",
	// 		"test_2~fv=00000000.js"
	// 	],
	// 	"SrcReg": {},
	// 	"Dist": "test/watch/dist",
	// 	"UseSAVR": false,
	// 	"Info": {
	// 		"PV": {
	// 			"subdir/test_3.js": "gia0-_Z_",
	// 			"subdir/test_4.js": "da1EKBXZ",
	// 			"test_1.min.js": "1sTEzePc",
	// 			"test_2.js": "qBbNrrTr"
	// 		},
	// 		"SAVR": "",
	// 		"VersionedFiles": [
	// 			"subdir/test_3~fv=gia0-_Z_.js",
	// 			"subdir/test_4~fv=da1EKBXZ.js",
	// 			"test_1~fv=1sTEzePc.min.js",
	// 			"test_2~fv=qBbNrrTr.js"
	// 		],
	// 		"Index": null,
	// 		"TotalSourceReplaces": 20,
	// 		"UpdatedFilePaths": [
	// 			"test/watch/dist/subdir/test_3~fv=gia0-_Z_.js",
	// 			"test/watch/dist/subdir/test_4~fv=da1EKBXZ.js",
	// 			"test/watch/dist/test_1.min.js.map",
	// 			"test/watch/dist/test_1~fv=1sTEzePc.min.js",
	// 			"test/watch/dist/test_2~fv=qBbNrrTr.js"
	// 		]
	// 	}
	// }
	// File test/watch/dist/subdir/test_3~fv=gia0-_Z_.js:
	// ////////////////
	// import * as test1 from '../test_1~fv=1sTEzePc.min.js';
	// import * as test2 from '../test_2~fv=qBbNrrTr.js';
	// import * as test4 from '../subdir/test_4~fv=da1EKBXZ.js';
	// // Comments referring to './test_1~fv=1sTEzePc.min.js' should be updated as well.
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
	// [test_1~fv=qlJgGoFM.js test_2~fv=8RBMUSqr.js subdir/test_3~fv=_X83uO__.js subdir/test_4~fv=2lfgwhXI.js]
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
	// [subdir/test_3~fv=00000000.js subdir/test_4~fv=00000000.js test_1~fv=00000000.js test_2~fv=00000000.js]
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
	// [subdir/test_3~fv=_X83uO__.js subdir/test_4~fv=GJIrg6k1.js test_1~fv=vPCb4GVO.js test_2~fv=BOl7h9TM.js]
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

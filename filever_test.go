package filever

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/cyphrme/watch"
)

var dummySrc = "test/dummy/src"
var dummyDist = "test/dummy/dist"
var dummyEndSrc = "test/dummy_end/src"
var dummyEndDist = "test/dummy_end/dist"
var dummyNoSrc = "test/dummy_no/src"
var dummyNoDist = "test/dummy_no/dist"
var watchSrc = "test/watch/src"
var watchDist = "test/watch/dist"

// Example VersionReplace with mid version with "dummy" input files.
// go test -run '^ExampleVersionReplace$'
func ExampleVersionReplace() {
	clean()
	c := &Config{Src: dummySrc, Dist: dummyDist}
	err := VersionReplace(c)
	if err != nil {
		panic(err)
	}

	PrintPretty(c.Info)
	PrintFile(dummyDist + "/" + c.Info.VersionedFiles[0])

	// Output:
	// {
	// 	"PV": {
	// 		"subdir/test_3.js": "JPWlGq2g",
	// 		"test_1.js": "vk-N-Yhv",
	// 		"test_2.js": "HNWQ4g9G",
	// 		"test_4.js": "Dgc1dsET"
	// 	},
	// 	"SAVR": "(subdir/test_3\\?fv=[0-9A-Za-z_-]*.js)|(test_1\\?fv=[0-9A-Za-z_-]*.js)|(test_2\\?fv=[0-9A-Za-z_-]*.js)|(test_4\\?fv=[0-9A-Za-z_-]*.js)",
	// 	"VersionedFiles": [
	// 		"subdir/test_3?fv=JPWlGq2g.js",
	// 		"test_1?fv=vk-N-Yhv.js",
	// 		"test_2?fv=HNWQ4g9G.js",
	// 		"test_4?fv=Dgc1dsET.js"
	// 	],
	// 	"TotalSourceReplaces": 8,
	// 	"UpdatedFilePaths": [
	// 		"test/dummy/dist/subdir/test_3?fv=JPWlGq2g.js",
	// 		"test/dummy/dist/test_1?fv=vk-N-Yhv.js",
	// 		"test/dummy/dist/test_2?fv=HNWQ4g9G.js",
	// 		"test/dummy/dist/test_4?fv=Dgc1dsET.js"
	// 	]
	// }
	// File test/dummy/dist/subdir/test_3?fv=JPWlGq2g.js:
	// ////////////////
	// import * as test1 from './test_1?fv=vk-N-Yhv.js';
	// import * as test2 from './test_2?fv=HNWQ4g9G.js';
	// ////////////////
}

// Example VersionReplace using end ver format and dummy file inputs.
func ExampleVersionReplace_end() {
	clean()

	c := &Config{Src: dummyEndSrc, Dist: dummyEndDist, EndVer: true}
	err := VersionReplace(c)
	if err != nil {
		panic(err)
	}

	PrintPretty(c.Info)
	// Print the first file to ensure contents updated with appropriate version.
	PrintFile(dummyEndDist + "/" + c.Info.VersionedFiles[0])

	// Output:
	// {
	// 	"PV": {
	// 		"subdir/test_3.js": "gGkfoWxG",
	// 		"test_1.js": "iwjsfpt6",
	// 		"test_2.js": "k4Ti_VfV",
	// 		"test_4.js": "PRrfCtOC"
	// 	},
	// 	"SAVR": "(subdir/test_3\\.js\\?fv=[0-9A-Za-z_-]*)|(test_1\\.js\\?fv=[0-9A-Za-z_-]*)|(test_2\\.js\\?fv=[0-9A-Za-z_-]*)|(test_4\\.js\\?fv=[0-9A-Za-z_-]*)",
	// 	"VersionedFiles": [
	// 		"subdir/test_3.js?fv=gGkfoWxG",
	// 		"test_1.js?fv=iwjsfpt6",
	// 		"test_2.js?fv=k4Ti_VfV",
	// 		"test_4.js?fv=PRrfCtOC"
	// 	],
	// 	"TotalSourceReplaces": 8,
	// 	"UpdatedFilePaths": [
	// 		"test/dummy_end/dist/subdir/test_3.js?fv=gGkfoWxG",
	// 		"test/dummy_end/dist/test_1.js?fv=iwjsfpt6",
	// 		"test/dummy_end/dist/test_2.js?fv=k4Ti_VfV",
	// 		"test/dummy_end/dist/test_4.js?fv=PRrfCtOC"
	// 	]
	// }
	// File test/dummy_end/dist/subdir/test_3.js?fv=gGkfoWxG:
	// ////////////////
	// import * as test1 from '../test_1.js?fv=iwjsfpt6';
	// import * as test2 from '../test_2.js?fv=k4Ti_VfV';
	// ////////////////
}

// Example_watchVersionAndReplace demonstrates using FileVer with the external
// program "watch". Uses the "mid version" format.
func Example_watchVersionAndReplace() {
	clean()

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

	PrintPretty(c.Info)
	PrintFile(watchDist + "/" + c.Info.VersionedFiles[0])

	// Output:
	// Flag `daemon` set to false.  Running commands in config and exiting.
	// {
	// 	"PV": {
	// 		"subdir/test_3.js": "JPWlGq2g",
	// 		"test_1.min.js": "QCdhFObN",
	// 		"test_1.min.js.map": "iBbcDWoU",
	// 		"test_2.js": "mighoZca",
	// 		"test_4.js": "BRW8FEPq"
	// 	},
	// 	"SAVR": "(subdir/test_3\\?fv=[0-9A-Za-z_-]*.js)|(test_1\\?fv=[0-9A-Za-z_-]*.min.js)|(test_1\\?fv=[0-9A-Za-z_-]*.min.js.map)|(test_2\\?fv=[0-9A-Za-z_-]*.js)|(test_4\\?fv=[0-9A-Za-z_-]*.js)",
	// 	"VersionedFiles": [
	// 		"subdir/test_3?fv=JPWlGq2g.js",
	// 		"test_1?fv=QCdhFObN.min.js",
	// 		"test_1?fv=iBbcDWoU.min.js.map",
	// 		"test_2?fv=mighoZca.js",
	// 		"test_4?fv=BRW8FEPq.js"
	// 	],
	// 	"TotalSourceReplaces": 8,
	// 	"UpdatedFilePaths": [
	// 		"test/watch/dist/subdir/test_3?fv=JPWlGq2g.js",
	// 		"test/watch/dist/test_1?fv=QCdhFObN.min.js",
	// 		"test/watch/dist/test_1?fv=iBbcDWoU.min.js.map",
	// 		"test/watch/dist/test_2?fv=mighoZca.js",
	// 		"test/watch/dist/test_4?fv=BRW8FEPq.js"
	// 	]
	// }
	// File test/watch/dist/subdir/test_3?fv=JPWlGq2g.js:
	// ////////////////
	// import * as test1 from './test_1?fv=00000000.js';
	// import * as test2 from './test_2?fv=mighoZca.js';
	// ////////////////
}

// Example_noDummy demonstrates inputting "manually" enumerated files to be
// processes by FileVer, aka it does not use dummy file inputs. Does not do
// Replace().
func Example_noDummy() {
	clean()
	c := &Config{
		Src:  dummyNoSrc,
		Dist: dummyNoDist,
		SrcFiles: []string{
			"test_1.js",
			"test_2.js",
			"test_4.js",
			"subdir/test_3.js",
		},
	}

	Version(c)
	// TODO Replace does not appear to be working.
	fmt.Println(c.Info.VersionedFiles)
	// PrintPretty(c.Info)
	// PrintFile(dummyNoDist + "/" + c.Info.VersionedFiles[0])
	// Output:
	// [test_1?fv=iwjsfpt6.js test_2?fv=k4Ti_VfV.js test_4?fv=1M7qm-Pd.js subdir/test_3?fv=HlNaJEAj.js]
}

func ExampleListFilesInPath() {
	f, err := ListFilesInPath(dummyNoSrc)
	if err != nil {
		panic(err)
	}

	fmt.Println(f)

	// Output:
	// [test_1.js test_2.js test_4.js unversioned_file.md]
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
	//[subdir/test_3?fv=00000000.js test_1?fv=00000000.js test_2?fv=00000000.js test_4?fv=00000000.js]
	//[subdir/test_3.js?fv=00000000 test_1.js?fv=00000000 test_2.js?fv=00000000 test_4.js?fv=00000000]
}

func ExampleCleanVersionFiles() {
	// Generate test files.
	clean()
	c := &Config{Src: dummySrc, Dist: dummyDist}
	err := VersionReplace(c)
	if err != nil {
		panic(err)
	}
	// Clean out generate test files.
	CleanVersionFiles(dummyEndDist)
	f, err := ListFilesInPath(dummyEndDist)
	if err != nil {
		panic(err)
	}

	fmt.Println(c.Info.VersionedFiles)
	fmt.Println(f)
	// Output:
	// [subdir/test_3?fv=JPWlGq2g.js test_1?fv=vk-N-Yhv.js test_2?fv=HNWQ4g9G.js test_4?fv=Dgc1dsET.js]
	// [not_versioned_example.txt]
}

func Test_clean(t *testing.T) {
	clean()
}

func Test__nuke(t *testing.T) {
	nukeAndRebuildTestDirs()
}

// Helper Functions

// Clean removes versioned files from dist.
func clean() {
	c := []string{
		dummyDist,
		dummyEndDist,
		dummyNoDist,
		watchDist,
	}

	for _, v := range c {
		CleanVersionFiles(v)
	}
}

// Completely delete all test dirs and recreate.
func nukeAndRebuildTestDirs() {
	c := []string{
		dummyDist,
		dummyEndDist,
		dummyNoDist,
		watchDist,
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

		d1 := []byte("This file lives in dist and is not versioned.")
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

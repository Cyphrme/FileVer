package filever

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cyphrme/watch"
)

var dummyNoSrc = "test/dummy_no/src"
var dummyNoDist = "test/dummy_no/dist"
var dummySrc = "test/dummy/src"
var dummyDist = "test/dummy/dist"
var dummyMidSrc = "test/dummy_mid/src"
var dummyMidDist = "test/dummy_mid/dist"
var watchSrc = "test/watch/src"
var watchDist = "test/watch/dist"

func Example_clean() {
	clean()
	fmt.Println("Clean done. ")
	// Output:
	// Clean done.
}

// Clean removes versioned files from dist.
func clean() {
	c := []string{
		dummyNoDist,
		dummyDist,
		dummyMidDist,
		watchDist,
	}

	for _, v := range c {
		CleanVersionFiles(v)
	}
}

// Completely delete all test dirs and recreate.
func nukeAndRebuildTestDirs() {
	c := []string{
		dummyNoDist,
		dummyDist,
		dummyMidDist,
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

func Example_nuke() {
	nukeAndRebuildTestDirs()
	// Output:
}

func ExampleCleanVersionFiles() {
	Example_dummy() // Generates files
	CleanVersionFiles(dummyDist)
	f, err := ListFilesInPath(dummyDist)
	if err != nil {
		panic(err)
	}

	fmt.Println(f)
	// Output:
	// [test_1.js?fv=iwjsfpt6 test_2.js?fv=k4Ti_VfV test_4.js?fv=FS0yyWFi subdir/test_3.js?fv=HlNaJEAj]
	// [not_versioned_example.txt]
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
	files, err := ExistingVersionedFiles(dummySrc)
	if err != nil {
		panic(err)
	}
	fmt.Println(files)

	// Mid Version
	files, err = ExistingVersionedFiles(dummyMidSrc)
	if err != nil {
		panic(err)
	}
	fmt.Println(files)

	// Output:
	//[subdir/test_3.js?fv=00000000 test_1.js?fv=00000000 test_2.js?fv=00000000 test_4.js?fv=00000000]
	//[subdir/test_3?fv=00000000.js test_1?fv=00000000.js test_2?fv=00000000.js test_4?fv=00000000.js]
}

// Example_noDummy demonstrates "manually" enumerating files to be processes by
// FileVer. // TODO replace appears to not be working.
func Example_noDummy() {
	clean()
	c := &Config{Src: dummyNoSrc, Dist: dummyNoDist}
	scrFiles := []string{
		"test_1.js",
		"test_2.js",
		"test_4.js",
		"subdir/test_3.js",
	}

	files := []string{}
	for _, f := range scrFiles {
		file, err := FileToFileVerOutputDelete(f, c)
		if err != nil {
			panic(err)
		}

		files = append(files, file)
	}

	fmt.Printf("%s", files)
	// Output:
	// [test_1.js?fv=iwjsfpt6 test_2.js?fv=k4Ti_VfV test_4.js?fv=1M7qm-Pd subdir/test_3.js?fv=HlNaJEAj]
}

// Example_dummy takes input files from many directories (starting at
// `/test/dummy/src`), outputs to a single dist directory  (`/test/dummy/dist`).
// This is an example of a "global", unstructured dist.
func Example_dummy_end() {
	clean()
	c := &Config{Src: dummySrc, Dist: dummyDist, EndVer: true}
	scrFiles, err := ExistingVersionedFiles(c.Src)
	if err != nil {
		panic(err)
	}
	files := []string{}
	for _, f := range scrFiles {
		file, err := FileToFileVerOutputDelete(f, c)
		if err != nil {
			panic(err)
		}

		files = append(files, file)
	}

	fmt.Println(files)

	//Replace(c)

	// Output:
	//[subdir/test_3.js?fv=HlNaJEAj test_1.js?fv=iwjsfpt6 test_2.js?fv=k4Ti_VfV test_4.js?fv=liDycGy1]
}

// Example_mid_dummy takes input files from many directories (starting at
// `/test/dummy_mid/src`), outputs to a single dist directory  (`/test/dummy_mid/dist`).
// This is an example of a "global", unstructured dist.
func Example_dummy() {
	clean()
	c := &Config{Src: dummyMidSrc, Dist: dummyMidDist}
	scrFiles, err := ExistingVersionedFiles(c.Src)
	if err != nil {
		panic(err)
	}
	files := []string{}
	for _, f := range scrFiles {
		file, err := FileToFileVerOutputDelete(f, c)
		if err != nil {
			panic(err)
		}

		files = append(files, file)
	}

	fmt.Println(files)

	// Output:
	//[test_1?fv=jd2_eypN.js test_2?fv=1XdEl8NR.js subdir/test_3?fv=1XdEl8NR.js]
}

// Example VersionReplace.  No mid ver.
func ExampleVersionReplace_end() {
	clean()

	c := &Config{Src: dummySrc, Dist: dummyDist, EndVer: true}
	err := VersionReplace(c)
	if err != nil {
		panic(err)
	}

	PrintPretty(c.Info)
	// Print the first file to ensure contents updated with appropriate version.
	PrintFile(dummyDist + "/" + c.Info.VersionedFiles[0])

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
	// 		"test/dummy/dist/subdir/test_3.js?fv=gGkfoWxG",
	// 		"test/dummy/dist/test_1.js?fv=iwjsfpt6",
	// 		"test/dummy/dist/test_2.js?fv=k4Ti_VfV",
	// 		"test/dummy/dist/test_4.js?fv=PRrfCtOC"
	// 	]
	// }
	// File test/dummy/dist/subdir/test_3.js?fv=gGkfoWxG:
	// ////////////////
	// import * as test1 from '../test_1.js?fv=iwjsfpt6';
	// import * as test2 from '../test_2.js?fv=k4Ti_VfV';
	// ////////////////
}

// Example VersionReplace with mid version.
func ExampleVersionReplace() {
	clean()

	c := &Config{Src: dummyMidSrc, Dist: dummyMidDist}
	err := VersionReplace(c)
	if err != nil {
		panic(err)
	}

	PrintPretty(c.Info)
	PrintFile(dummyMidDist + "/" + c.Info.VersionedFiles[0])

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
	// 		"test/dummy_mid/dist/subdir/test_3?fv=JPWlGq2g.js",
	// 		"test/dummy_mid/dist/test_1?fv=vk-N-Yhv.js",
	// 		"test/dummy_mid/dist/test_2?fv=HNWQ4g9G.js",
	// 		"test/dummy_mid/dist/test_4?fv=Dgc1dsET.js"
	// 	]
	// }
	// File test/dummy_mid/dist/subdir/test_3?fv=JPWlGq2g.js:
	// ////////////////
	// import * as test1 from './test_1?fv=vk-N-Yhv.js';
	// import * as test2 from './test_2?fv=HNWQ4g9G.js';
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
	// 		"subdir/test_3.js": "w34ZX0Lq",
	// 		"test_1.min.js": "nVqEy2SO",
	// 		"test_1.min.js.map": "Qk-Y21_2",
	// 		"test_2.js": "w34ZX0Lq"
	// 	},
	// 	"SAVR": "(subdir/test_3\\.js\\?fv=[0-9A-Za-z_-]*)|(test_1\\.min\\.js\\.map\\?fv=[0-9A-Za-z_-]*)|(test_1\\.min\\.js\\?fv=[0-9A-Za-z_-]*)|(test_2\\.js\\?fv=[0-9A-Za-z_-]*)",
	// 	"VersionedFiles": [
	// 		"subdir/test_3.js?fv=w34ZX0Lq",
	// 		"test_1.min.js.map?fv=Qk-Y21_2",
	// 		"test_1.min.js?fv=nVqEy2SO",
	// 		"test_2.js?fv=w34ZX0Lq"
	// 	],
	// 	"TotalSourceReplaces": 5,
	// 	"UpdatedFilePaths": [
	// 		"test/watch/dist/subdir/test_3.js?fv=w34ZX0Lq",
	// 		"test/watch/dist/test_1.min.js.map?fv=Qk-Y21_2",
	// 		"test/watch/dist/test_1.min.js?fv=nVqEy2SO",
	// 		"test/watch/dist/test_2.js?fv=w34ZX0Lq"
	// 	]
	// }

}

// PrintFile is a helper function
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

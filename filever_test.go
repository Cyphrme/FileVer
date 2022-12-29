package filever

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

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

func clean() {

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
	}
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
	//[test_1.js?fv=00000000 test_2.js?fv=00000000 test_folder/test_3.js?fv=00000000]
	//[test_1?fv=00000000.js test_2?fv=00000000.js test_folder/test_3?fv=00000000.js]
}

// Example_noDummy demonstrates "manually" enumerating files to be processes by
// FileVer.
func Example_noDummy() {
	clean()
	c := &Config{Src: dummyNoSrc, Dist: dummyNoDist}
	scrFiles := []string{
		"test_1.js",
		"test_2.js",
		"test_folder/test_3.js",
	}

	files := []string{}
	for _, f := range scrFiles {
		file, err := FileToFileVerOutputDelete(f, c)
		if err != nil {
			panic(err)
		}

		files = append(files, file)
	}

	fmt.Println("Processed: " + strings.Join(files, ", "))
	// Output:
	// Processed: test_1.js?fv=4WYoW0MN, test_2.js?fv=zmoLIyPU, test_folder/test_3.js?fv=zmoLIyPU
}

// Example_dummy takes input files from many directories (starting at
// `/test/dummy/src`), outputs to a single dist directory  (`/test/dummy/dist`).
// This is an example of a "global", unstructured dist.
func Example_dummy() {
	clean()
	c := &Config{Src: dummySrc, Dist: dummyDist}
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
	//[test_1.js?fv=4WYoW0MN test_2.js?fv=zmoLIyPU test_folder/test_3.js?fv=zmoLIyPU]
}

// Example_dummy takes input files from many directories (starting at
// `/test/dummy/src`), outputs to a single dist directory  (`/test/dummy/dist`).
// This is an example of a "global", unstructured dist.
func Test_dummy(t *testing.T) {
	clean()
	c := &Config{Src: dummySrc, Dist: dummyDist}
	Version(c)
	fmt.Printf("Post Version config: %+v, \ninfo: %+v", c, c.Info)

	Replace(c)

	// Output:
	//[test_1.js?fv=4WYoW0MN test_2.js?fv=zmoLIyPU test_folder/test_3.js?fv=zmoLIyPU]
}

// Example_mid_dummy takes input files from many directories (starting at
// `/test/dummy_mid/src`), outputs to a single dist directory  (`/test/dummy_mid/dist`).
// This is an example of a "global", unstructured dist.
func Example_mid_dummy() {
	clean()
	c := &Config{Src: dummyMidSrc, Dist: dummyMidDist, MidVer: true}
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
	//[test_1?fv=jd2_eypN.js test_2?fv=1XdEl8NR.js test_folder/test_3?fv=1XdEl8NR.js]
}

func ExampleVersionReplace() {
	clean()

	c := &Config{Src: dummySrc, Dist: dummyDist}
	err := VersionReplace(c)
	if err != nil {
		panic(err)
	}

	PrintPretty(c.Info)
	PrintFile("test/dummy/dist/test_1.js?fv=YNCUu8I7")

	// Output:
	// Replace Config  &{Src:test/dummy/src Dist:test/dummy/dist MidVer:false Info:0xc000132700} Info: &{TotalSourceReplaces:0 UpdatedFilePaths:[] SAVR: PV:map[test_1.js:?fv=4WYoW0MN test_2.js:?fv=zmoLIyPU test_folder/test_3.js:?fv=zmoLIyPU] VersionedFiles:[test_1.js?fv=4WYoW0MN test_2.js?fv=zmoLIyPU test_folder/test_3.js?fv=zmoLIyPU] CurrentPath: CurrentMatches:0}{
	// 	"TotalSourceReplaces": 3,
	// 	"UpdatedFilePaths": [
	// 		"test/dummy/dist/test_1.js?fv=4WYoW0MN",
	// 		"test/dummy/dist/test_2.js?fv=zmoLIyPU",
	// 		"test/dummy/dist/test_folder/test_3.js?fv=zmoLIyPU"
	// 	],
	// 	"SAVR": "(test_1\\.js\\?fv=[0-9A-Za-z_-]*)|(test_2\\.js\\?fv=[0-9A-Za-z_-]*)|(test_folder/test_3\\.js\\?fv=[0-9A-Za-z_-]*)",
	// 	"PV": {
	// 		"test_1.js": "?fv=4WYoW0MN",
	// 		"test_2.js": "?fv=zmoLIyPU",
	// 		"test_folder/test_3.js": "?fv=zmoLIyPU"
	// 	},
	// 	"VersionedFiles": [
	// 		"test_1.js?fv=4WYoW0MN",
	// 		"test_2.js?fv=zmoLIyPU",
	// 		"test_folder/test_3.js?fv=zmoLIyPU"
	// 	]
	// }
	//
	// File:
	// ////////////////
	// "use strict";
	//
	// import * as test2 from './test_2.js?fv=zmoLIyPU';
	// ////////////////
}

func TestVersionReplace(t *testing.T) {
	clean()

	c := &Config{Src: dummySrc, Dist: dummyDist}
	err := VersionReplace(c)
	if err != nil {
		panic(err)
	}

	PrintPretty(c.Info)
	PrintFile("test/dummy/dist/test_1.js?fv=YNCUu8I7")
}

// This example demonstrates using FileVer with the program "watch".
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

	// Pretty Print
	json, err := json.MarshalIndent(c.Info, "", "\t")
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Println(string(json))

	// Output:
	//Flag `daemon` set to false.  Running commands in config and exiting.
	// {
	// 	"TotalSourceReplaces": 5,
	// 	"UpdatedFilePaths": [
	// 		"test/watch/dist/test_1.min.js.map?fv=820OsC4y",
	// 		"test/watch/dist/test_1.min.js?fv=nVqEy2SO",
	// 		"test/watch/dist/test_2.js?fv=Gv-m9hRr",
	// 		"test/watch/dist/test_folder/test_3.js?fv=Gv-m9hRr"
	// 	],
	// 	"SAVR": "(test_1\\.min\\.js\\.map\\?fv=[0-9A-Za-z_-]*)|(test_1\\.min\\.js\\?fv=[0-9A-Za-z_-]*)|(test_2\\.js\\?fv=[0-9A-Za-z_-]*)|(test_folder/test_3\\.js\\?fv=[0-9A-Za-z_-]*)",
	// 	"PV": {
	// 		"test_1.min.js": "?fv=nVqEy2SO",
	// 		"test_1.min.js.map": "?fv=820OsC4y",
	// 		"test_2.js": "?fv=Gv-m9hRr",
	// 		"test_folder/test_3.js": "?fv=Gv-m9hRr"
	// 	},
	// 	"VersionedFiles": [
	// 		"test_1.min.js.map?fv=820OsC4y",
	// 		"test_1.min.js?fv=nVqEy2SO",
	// 		"test_2.js?fv=Gv-m9hRr",
	// 		"test_folder/test_3.js?fv=Gv-m9hRr"
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

	fmt.Printf("\nFile:\n////////////////\n%s\n////////////////", b)

}

func PrintPretty(s any) {
	json, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(json))
}

package filever

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/cyphrme/watch"
)

var dummyNoSrc = "test/dummy_no/src"
var dummyNoDist = "test/dummy_no/dist"

var dummySrc = "test/dummy/src"
var dummyDist = "test/dummy/dist"

var watchSrc = "test/watch/src"
var watchDist = "test/watch/dist"

func Example_clean() {
	clean()
	fmt.Println("Clean done. ")
	// Output:
	// Clean done.
}

func clean() {
	// Dummy No
	err := os.RemoveAll(dummyNoDist)
	if err != nil {
		panic(err)
	}

	err = os.Mkdir(dummyNoDist, 0777)
	if err != nil {
		panic(err)
	}

	// Dummy
	err = os.RemoveAll(dummyDist)
	if err != nil {
		panic(err)
	}

	err = os.Mkdir(dummyDist, 0777)
	if err != nil {
		panic(err)
	}

	// Watch
	err = os.RemoveAll(watchDist)
	if err != nil {
		panic(err)
	}

	err = os.Mkdir(watchDist, 0777)
	if err != nil {
		panic(err)
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

func ExampleDirVersionedFiles() {
	files, err := DirVersionedFiles(dummySrc)
	if err != nil {
		panic(err)
	}
	fmt.Println(files)

	// Output:
	//[test_1.js?fv=00000000 test_2.js?fv=00000000 test_folder/test_3.js?fv=00000000]
}

// Example_noDummy demonstrates "manually" enumerating files to be processes by
// FileVer.
func Example_noDummy() {
	clean()
	scrFiles := []string{
		"test_1.js",
		"test_2.js",
		"test_folder/test_3.js",
	}

	c := &Config{Src: dummyNoSrc, Dist: dummyNoDist}

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

// Input files from many directories, output dist to a single directory.  This
// is an example of a "global", unstructured dist.
func Example_dummy_from_many_dist_to_one() {
	clean()

	c := &Config{Src: dummySrc, Dist: dummyDist}

	scrFiles, err := DirVersionedFiles(dummySrc)
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
	//[test_1.js?fv=4WYoW0MN test_2.js?fv=zmoLIyPU test_folder/test_3.js?fv=zmoLIyPU]
}

func ExampleVersionReplace() {
	clean()

	c := &Config{Src: dummySrc, Dist: dummyDist}
	err := VersionReplace(c)
	if err != nil {
		panic(err)
	}

	// Pretty Print
	json, err := json.MarshalIndent(c.Info, "", "\t")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(json))

	// Output:
	// {
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

package filever

import "fmt"

func ExamplePopulated() {

	paths := []string{
		"app",
		"app.js",
		"app.min.js",
		"e/app.min.js",
		"/a/e/app.min.js",

		// FileVer'd paths
		"app~fv=4mIbJJPq",
		"e/app~fv=4mIbJJPq",
		"e/app~fv=4mIbJJPq.js",
		"e/app~fv=4mIbJJPq.min.js",
		"/a/e/app~fv=4mIbJJPq.min.js",
	}

	for _, v := range paths {
		p := Populated(v)

		//fmt.Printf("%+v\n", p)
		PrintPretty(p)

	}

	// Output:
	// {
	// 	"Full": "app",
	// 	"File": "app",
	// 	"Base": "app",
	// 	"BarePath": "app",
	// 	"BareFile": "app",
	// 	"Bare": "app"
	// }
	// {
	// 	"Full": "app.js",
	// 	"File": "app.js",
	// 	"Base": "app",
	// 	"Ext": ".js",
	// 	"BaseExt": ".js",
	// 	"BarePath": "app.js",
	// 	"BareFile": "app.js",
	// 	"Bare": "app"
	// }
	// {
	// 	"Full": "app.min.js",
	// 	"File": "app.min.js",
	// 	"Base": "app",
	// 	"Ext": ".min.js",
	// 	"BaseExt": ".js",
	// 	"BarePath": "app.min.js",
	// 	"BareFile": "app.min.js",
	// 	"Bare": "app"
	// }
	// {
	// 	"Full": "e/app.min.js",
	// 	"Dir": "e/",
	// 	"File": "app.min.js",
	// 	"Base": "app",
	// 	"Ext": ".min.js",
	// 	"BaseExt": ".js",
	// 	"BarePath": "e/app.min.js",
	// 	"BareFile": "app.min.js",
	// 	"Bare": "app"
	// }
	// {
	// 	"Full": "/a/e/app.min.js",
	// 	"Dir": "/a/e/",
	// 	"File": "app.min.js",
	// 	"Base": "app",
	// 	"Ext": ".min.js",
	// 	"BaseExt": ".js",
	// 	"BarePath": "/a/e/app.min.js",
	// 	"BareFile": "app.min.js",
	// 	"Bare": "app"
	// }
	// {
	// 	"Full": "app~fv=4mIbJJPq",
	// 	"File": "app~fv=4mIbJJPq",
	// 	"Base": "app~fv=4mIbJJPq",
	// 	"FileVer": "app~fv=4mIbJJPq",
	// 	"DelimVer": "~fv=4mIbJJPq",
	// 	"Version": "4mIbJJPq",
	// 	"BarePath": "app",
	// 	"BareFile": "app",
	// 	"Bare": "app"
	// }
	// {
	// 	"Full": "e/app~fv=4mIbJJPq",
	// 	"Dir": "e/",
	// 	"File": "app~fv=4mIbJJPq",
	// 	"Base": "app~fv=4mIbJJPq",
	// 	"FileVer": "app~fv=4mIbJJPq",
	// 	"DelimVer": "~fv=4mIbJJPq",
	// 	"Version": "4mIbJJPq",
	// 	"BarePath": "e/app",
	// 	"BareFile": "app",
	// 	"Bare": "app"
	// }
	// {
	// 	"Full": "e/app~fv=4mIbJJPq.js",
	// 	"Dir": "e/",
	// 	"File": "app~fv=4mIbJJPq.js",
	// 	"Base": "app~fv=4mIbJJPq",
	// 	"Ext": ".js",
	// 	"BaseExt": ".js",
	// 	"FileVer": "app~fv=4mIbJJPq.js",
	// 	"DelimVer": "~fv=4mIbJJPq",
	// 	"Version": "4mIbJJPq",
	// 	"BarePath": "e/app.js",
	// 	"BareFile": "app.js",
	// 	"Bare": "app"
	// }
	// {
	// 	"Full": "e/app~fv=4mIbJJPq.min.js",
	// 	"Dir": "e/",
	// 	"File": "app~fv=4mIbJJPq.min.js",
	// 	"Base": "app~fv=4mIbJJPq",
	// 	"Ext": ".min.js",
	// 	"BaseExt": ".js",
	// 	"FileVer": "app~fv=4mIbJJPq.min.js",
	// 	"DelimVer": "~fv=4mIbJJPq",
	// 	"Version": "4mIbJJPq",
	// 	"BarePath": "e/app.min.js",
	// 	"BareFile": "app.min.js",
	// 	"Bare": "app"
	// }
	// {
	// 	"Full": "/a/e/app~fv=4mIbJJPq.min.js",
	// 	"Dir": "/a/e/",
	// 	"File": "app~fv=4mIbJJPq.min.js",
	// 	"Base": "app~fv=4mIbJJPq",
	// 	"Ext": ".min.js",
	// 	"BaseExt": ".js",
	// 	"FileVer": "app~fv=4mIbJJPq.min.js",
	// 	"DelimVer": "~fv=4mIbJJPq",
	// 	"Version": "4mIbJJPq",
	// 	"BarePath": "/a/e/app.min.js",
	// 	"BareFile": "app.min.js",
	// 	"Bare": "app"
	// }

}

func ExamplePathCut() {
	fmt.Println(PathCut("..subdir/test_5~fv=wwjNHrIw.js"))
	// Output:
	// ..subdir/ test_5~fv=wwjNHrIw.js
}

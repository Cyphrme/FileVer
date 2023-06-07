package filever

import "testing"

func TestPop(t *testing.T) {
	Populated("/a/e/app.min.js")
}

// go test -run ExamplePopulated$
func ExamplePopulated() {
	paths := []string{
		"",
		"/",
		"app",
		"/app",
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
		PrintPretty(p)
	}

	// Output:
	// {
	// 	"full": ""
	// }
	// {
	// 	"full": "/"
	// }
	// {
	// 	"full": "app",
	// 	"file": "app",
	// 	"file_base": "app",
	// 	"bare_path": "app",
	// 	"bare_file": "app",
	// 	"bare": "app"
	// }
	// {
	// 	"full": "/app",
	// 	"dir": "/",
	// 	"file": "app",
	// 	"file_base": "app",
	// 	"bare_path": "/app",
	// 	"bare_file": "app",
	// 	"bare": "app"
	// }
	// {
	// 	"full": "app.js",
	// 	"file": "app.js",
	// 	"file_base": "app",
	// 	"ext": ".js",
	// 	"ext_base": ".js",
	// 	"bare_path": "app.js",
	// 	"bare_file": "app.js",
	// 	"bare": "app"
	// }
	// {
	// 	"full": "app.min.js",
	// 	"file": "app.min.js",
	// 	"file_base": "app",
	// 	"ext": ".min.js",
	// 	"ext_base": ".js",
	// 	"bare_path": "app.min.js",
	// 	"bare_file": "app.min.js",
	// 	"bare": "app"
	// }
	// {
	// 	"full": "e/app.min.js",
	// 	"dir": "e/",
	// 	"file": "app.min.js",
	// 	"file_base": "app",
	// 	"ext": ".min.js",
	// 	"ext_base": ".js",
	// 	"bare_path": "e/app.min.js",
	// 	"bare_file": "app.min.js",
	// 	"bare": "app"
	// }
	// {
	// 	"full": "/a/e/app.min.js",
	// 	"dir": "/a/e/",
	// 	"file": "app.min.js",
	// 	"file_base": "app",
	// 	"ext": ".min.js",
	// 	"ext_base": ".js",
	// 	"bare_path": "/a/e/app.min.js",
	// 	"bare_file": "app.min.js",
	// 	"bare": "app"
	// }
	// {
	// 	"full": "app~fv=4mIbJJPq",
	// 	"file": "app~fv=4mIbJJPq",
	// 	"file_base": "app~fv=4mIbJJPq",
	// 	"filever": "app~fv=4mIbJJPq",
	// 	"delim_ver": "~fv=4mIbJJPq",
	// 	"version": "4mIbJJPq",
	// 	"bare_path": "app",
	// 	"bare_file": "app",
	// 	"bare": "app"
	// }
	// {
	// 	"full": "e/app~fv=4mIbJJPq",
	// 	"dir": "e/",
	// 	"file": "app~fv=4mIbJJPq",
	// 	"file_base": "app~fv=4mIbJJPq",
	// 	"filever": "app~fv=4mIbJJPq",
	// 	"delim_ver": "~fv=4mIbJJPq",
	// 	"version": "4mIbJJPq",
	// 	"bare_path": "e/app",
	// 	"bare_file": "app",
	// 	"bare": "app"
	// }
	// {
	// 	"full": "e/app~fv=4mIbJJPq.js",
	// 	"dir": "e/",
	// 	"file": "app~fv=4mIbJJPq.js",
	// 	"file_base": "app~fv=4mIbJJPq",
	// 	"ext": ".js",
	// 	"ext_base": ".js",
	// 	"filever": "app~fv=4mIbJJPq.js",
	// 	"delim_ver": "~fv=4mIbJJPq",
	// 	"version": "4mIbJJPq",
	// 	"bare_path": "e/app.js",
	// 	"bare_file": "app.js",
	// 	"bare": "app"
	// }
	// {
	// 	"full": "e/app~fv=4mIbJJPq.min.js",
	// 	"dir": "e/",
	// 	"file": "app~fv=4mIbJJPq.min.js",
	// 	"file_base": "app~fv=4mIbJJPq",
	// 	"ext": ".min.js",
	// 	"ext_base": ".js",
	// 	"filever": "app~fv=4mIbJJPq.min.js",
	// 	"delim_ver": "~fv=4mIbJJPq",
	// 	"version": "4mIbJJPq",
	// 	"bare_path": "e/app.min.js",
	// 	"bare_file": "app.min.js",
	// 	"bare": "app"
	// }
	// {
	// 	"full": "/a/e/app~fv=4mIbJJPq.min.js",
	// 	"dir": "/a/e/",
	// 	"file": "app~fv=4mIbJJPq.min.js",
	// 	"file_base": "app~fv=4mIbJJPq",
	// 	"ext": ".min.js",
	// 	"ext_base": ".js",
	// 	"filever": "app~fv=4mIbJJPq.min.js",
	// 	"delim_ver": "~fv=4mIbJJPq",
	// 	"version": "4mIbJJPq",
	// 	"bare_path": "/a/e/app.min.js",
	// 	"bare_file": "app.min.js",
	// 	"bare": "app"
	// }
}

package filever

import (
	"fmt"
)

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

func ExamplePopulated_uri() {
	paths := []string{
		"",
		"https://cyphr.me/coze",
		"https://cyphr.me/assets/img/cyphrme_long.png",
		"https://localhost:8081/",
		"https://localhost:8081",
		"sftp://example.com/joe/bob/file.txt",
	}

	for _, v := range paths {
		p := Populated(v)
		//fmt.Println(p)
		PrintPretty(p)
	}

	// Output:
	// {
	// 	"full": ""
	// }
	// {
	// 	"full": "https://cyphr.me/coze",
	// 	"dir": "https://cyphr.me/",
	// 	"file": "coze",
	// 	"file_base": "coze",
	// 	"bare_path": "https://cyphr.me/coze",
	// 	"bare_file": "coze",
	// 	"bare": "coze",
	// 	"scheme": "https",
	// 	"authority": "cyphr.me",
	// 	"host": "cyphr.me",
	// 	"uri_path": "/coze"
	// }
	// {
	// 	"full": "https://cyphr.me/assets/img/cyphrme_long.png",
	// 	"dir": "https://cyphr.me/assets/img/",
	// 	"file": "cyphrme_long.png",
	// 	"file_base": "cyphrme_long",
	// 	"ext": ".png",
	// 	"ext_base": ".png",
	// 	"bare_path": "https://cyphr.me/assets/img/cyphrme_long.png",
	// 	"bare_file": "cyphrme_long.png",
	// 	"bare": "cyphrme_long",
	// 	"scheme": "https",
	// 	"authority": "cyphr.me",
	// 	"host": "cyphr.me",
	// 	"uri_path": "/assets/img/cyphrme_long.png"
	// }
	// {
	// 	"full": "https://localhost:8081/",
	// 	"dir": "https://",
	// 	"file": "localhost:8081",
	// 	"file_base": "localhost:8081",
	// 	"bare_path": "https://localhost:8081",
	// 	"bare_file": "localhost:8081",
	// 	"bare": "localhost:8081",
	// 	"scheme": "https",
	// 	"authority": "localhost:8081",
	// 	"host": "localhost",
	// 	"port": ":8081"
	// }
	// {
	// 	"full": "https://localhost:8081",
	// 	"dir": "https://",
	// 	"file": "localhost:8081",
	// 	"file_base": "localhost:8081",
	// 	"bare_path": "https://localhost:8081",
	// 	"bare_file": "localhost:8081",
	// 	"bare": "localhost:8081",
	// 	"scheme": "https",
	// 	"authority": "localhost:8081",
	// 	"host": "localhost",
	// 	"port": ":8081"
	// }
	// {
	// 	"full": "sftp://example.com/joe/bob/file.txt",
	// 	"dir": "sftp://example.com/joe/bob/",
	// 	"file": "file.txt",
	// 	"file_base": "file",
	// 	"ext": ".txt",
	// 	"ext_base": ".txt",
	// 	"bare_path": "sftp://example.com/joe/bob/file.txt",
	// 	"bare_file": "file.txt",
	// 	"bare": "file",
	// 	"scheme": "sftp",
	// 	"authority": "example.com",
	// 	"host": "example.com",
	// 	"uri_path": "/joe/bob/file.txt"
	// }
}

func ExamplePathCut() {
	paths := []string{
		"app",
		"..subdir/test_5~fv=wwjNHrIw.js",
		"sftp://example.com/joe/bob/file.txt",
		"https://cyphr.me/coze",
		"https://cyphr.me/assets/img/cyphrme_long.png",
		"https://localhost:8081/",
		"https://localhost:8081",
	}

	for _, v := range paths {
		d, f := PathCut(v)
		fmt.Printf("dir: %s file: %s\n", d, f)
	}

	// Output:
	// dir:  file: app
	// dir: ..subdir/ file: test_5~fv=wwjNHrIw.js
	// dir: sftp://example.com/joe/bob/ file: file.txt
	// dir: https://cyphr.me/ file: coze
	// dir: https://cyphr.me/assets/img/ file: cyphrme_long.png
	// dir: https:// file: localhost:8081
	// dir: https:// file: localhost:8081
}

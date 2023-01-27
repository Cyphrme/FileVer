package filever

import "fmt"

// go test -run ExamplePopulated$
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

func ExamplePopulated_uri() {
	paths := []string{
		"https://cyphr.me/coze",
		"https://cyphr.me/assets/img/cyphrme_long.png",
		"https://localhost:8081/",
		"https://localhost:8081",
		"sftp://example.com/joe/bob/file.txt",
	}

	for _, v := range paths {
		p := Populated(v)
		PrintPretty(p)
	}

	// Output:
	// {
	// 	"Full": "https://cyphr.me/coze",
	// 	"Dir": "https://cyphr.me/",
	// 	"File": "coze",
	// 	"Base": "coze",
	// 	"BarePath": "https://cyphr.me/coze",
	// 	"BareFile": "coze",
	// 	"Bare": "coze",
	// 	"Scheme": "https",
	// 	"Authority": "cyphr.me",
	// 	"Host": "cyphr.me",
	// 	"URIPath": "/coze"
	// }
	// {
	// 	"Full": "https://cyphr.me/assets/img/cyphrme_long.png",
	// 	"Dir": "https://cyphr.me/assets/img/",
	// 	"File": "cyphrme_long.png",
	// 	"Base": "cyphrme_long",
	// 	"Ext": ".png",
	// 	"BaseExt": ".png",
	// 	"BarePath": "https://cyphr.me/assets/img/cyphrme_long.png",
	// 	"BareFile": "cyphrme_long.png",
	// 	"Bare": "cyphrme_long",
	// 	"Scheme": "https",
	// 	"Authority": "cyphr.me",
	// 	"Host": "cyphr.me",
	// 	"URIPath": "/assets/img/cyphrme_long.png"
	// }
	// {
	// 	"Full": "https://localhost:8081/",
	// 	"Dir": "https://",
	// 	"File": "localhost:8081",
	// 	"Base": "localhost:8081",
	// 	"BarePath": "https://localhost:8081",
	// 	"BareFile": "localhost:8081",
	// 	"Bare": "localhost:8081",
	// 	"Scheme": "https",
	// 	"Authority": "localhost:8081",
	// 	"Host": "localhost",
	// 	"Port": ":8081"
	// }
	// {
	// 	"Full": "https://localhost:8081",
	// 	"Dir": "https://",
	// 	"File": "localhost:8081",
	// 	"Base": "localhost:8081",
	// 	"BarePath": "https://localhost:8081",
	// 	"BareFile": "localhost:8081",
	// 	"Bare": "localhost:8081",
	// 	"Scheme": "https",
	// 	"Authority": "localhost:8081",
	// 	"Host": "localhost",
	// 	"Port": ":8081"
	// }
	// {
	// 	"Full": "sftp://example.com/joe/bob/file.txt",
	// 	"Dir": "sftp://example.com/joe/bob/",
	// 	"File": "file.txt",
	// 	"Base": "file",
	// 	"Ext": ".txt",
	// 	"BaseExt": ".txt",
	// 	"BarePath": "sftp://example.com/joe/bob/file.txt",
	// 	"BareFile": "file.txt",
	// 	"Bare": "file",
	// 	"Scheme": "sftp",
	// 	"Authority": "example.com",
	// 	"Host": "example.com",
	// 	"URIPath": "/joe/bob/file.txt"
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

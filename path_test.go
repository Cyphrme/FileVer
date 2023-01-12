package filever

import "fmt"

func ExamplePathParts_Populate() {

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

		p := PathParts{FullPath: v}
		p.Populate()
		fmt.Printf("%+v\n", p)

	}

	// Output:
	// {FullPath:app Path: File:app Base:app Ext: BaseExt: FileVer: BarePath:app BareFile:app Bare:app DelimVer: Version:}
	// {FullPath:app.js Path: File:app.js Base:app Ext:.js BaseExt:.js FileVer: BarePath:app.js BareFile:app.js Bare:app DelimVer: Version:}
	// {FullPath:app.min.js Path: File:app.min.js Base:app Ext:.min.js BaseExt:.js FileVer: BarePath:app.min.js BareFile:app.min.js Bare:app DelimVer: Version:}
	// {FullPath:e/app.min.js Path:e/ File:app.min.js Base:app Ext:.min.js BaseExt:.js FileVer: BarePath:e/app.min.js BareFile:app.min.js Bare:app DelimVer: Version:}
	// {FullPath:/a/e/app.min.js Path:/a/e/ File:app.min.js Base:app Ext:.min.js BaseExt:.js FileVer: BarePath:/a/e/app.min.js BareFile:app.min.js Bare:app DelimVer: Version:}
	// {FullPath:app~fv=4mIbJJPq Path: File:app~fv=4mIbJJPq Base:app~fv=4mIbJJPq Ext: BaseExt: FileVer:app~fv=4mIbJJPq BarePath:app BareFile:app Bare:app DelimVer:~fv=4mIbJJPq Version:4mIbJJPq}
	// {FullPath:e/app~fv=4mIbJJPq Path:e/ File:app~fv=4mIbJJPq Base:app~fv=4mIbJJPq Ext: BaseExt: FileVer:app~fv=4mIbJJPq BarePath:e/app BareFile:app Bare:app DelimVer:~fv=4mIbJJPq Version:4mIbJJPq}
	// {FullPath:e/app~fv=4mIbJJPq.js Path:e/ File:app~fv=4mIbJJPq.js Base:app~fv=4mIbJJPq Ext:.js BaseExt:.js FileVer:app~fv=4mIbJJPq.js BarePath:e/app.js BareFile:app.js Bare:app DelimVer:~fv=4mIbJJPq Version:4mIbJJPq}
	// {FullPath:e/app~fv=4mIbJJPq.min.js Path:e/ File:app~fv=4mIbJJPq.min.js Base:app~fv=4mIbJJPq Ext:.min.js BaseExt:.js FileVer:app~fv=4mIbJJPq.min.js BarePath:e/app.min.js BareFile:app.min.js Bare:app DelimVer:~fv=4mIbJJPq Version:4mIbJJPq}
	// {FullPath:/a/e/app~fv=4mIbJJPq.min.js Path:/a/e/ File:app~fv=4mIbJJPq.min.js Base:app~fv=4mIbJJPq Ext:.min.js BaseExt:.js FileVer:app~fv=4mIbJJPq.min.js BarePath:/a/e/app.min.js BareFile:app.min.js Bare:app DelimVer:~fv=4mIbJJPq Version:4mIbJJPq}

}

func ExamplePathCut() {
	fmt.Println(PathCut("..subdir/test_5~fv=wwjNHrIw.js"))
	// Output:
	// ..subdir/ test_5~fv=wwjNHrIw.js
}

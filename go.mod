module github.com/cyphrme/filever

go 1.19

// Go Mod and Go Work Zami Tutorial
// Go work/go mod/go has a bug: https://github.com/golang/go/issues/54264
// ** Cannot use `go work` as documented.  Must use `go mod` **
// Cannot use `go work` as documented with replace.  Until a bug fix, use
// replace in `go mod`.
//
// Here are the five (!) documents that need to be
// throughly understood before using `go mod` or `go work`.  Remember to use
// `replace`, the most essential part of either.
//
// Go Mod
// - [Go Modules Reference](https://go.dev/ref/mod#go-mod-file-replace)
// - [Managing dependencies](https://go.dev/doc/modules/managing-dependencies#local_directory)
// - [Module release and versioning workflow](https://go.dev/doc/modules/release-workflow#unpublished)
// Go Work
// - [Get familiar with workspaces](https://go.dev/blog/get-familiar-with-workspaces)
// - [Tutorial: Getting started with multi-module workspaces](https://go.dev/doc/tutorial/workspaces)
//
// # Go Workspace
// To set up workspaces, do the following in $GOPATH, which adds everything:
// `go.work` belongs here: `$GOPATH/go.work` `go work init` makes the file
// `go.work` with the single line, `go 1.19`, so don't do this.  Just make it
// manually.
// `go work use -r src` will add all modules currently in `src`
// `go env GOWORK` will be populated if go is in workspace mode.
// ```
// cd $GOPATH
// go work init
// go work use -r src
// go env GOWORK
// ```
//
// (The `!` character is the escape character for upper case directories)

replace (
	github.com/cyphrme/coze => ../coze
	github.com/cyphrme/watch => ../watch
)

require (
	github.com/cyphrme/coze v0.0.0-00010101000000-000000000000
	github.com/cyphrme/watch v0.0.0-00010101000000-000000000000
)

require (
	github.com/DisposaBoy/JsonConfigReader v0.0.0-20201129172854-99cf318d67e7 // indirect
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	golang.org/x/crypto v0.0.0-20220829220503-c86fa9a7ed90 // indirect
	golang.org/x/sys v0.0.0-20220829200755-d48e67d00261 // indirect
)

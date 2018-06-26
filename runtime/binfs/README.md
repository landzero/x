# BinFS

BinFS is a package for embedding one or more directories into Go binary.

# Usage

## Get BinFS

```
go get landzero.net/x/cmd/binfs # the cli
go get landzero.net/x/runtime/binfs # the runtime package
```

## Generate File

`PKG=pkgname binfs public view > binfs.gen.go`

This command read the content of directory `public` and `view`, output a `binfs.gen.go` file

The environment variable `PKG` is used for package name in `binfs.gen.go` file

## Use File

As long as `binfs.gen.go` is compiled with your source code, you can extract file with

```go
binfs.Open("/public/robots.txt")
```

You can also use `binfs.FileSystem()` to get a implementation of `http.FileSystem`

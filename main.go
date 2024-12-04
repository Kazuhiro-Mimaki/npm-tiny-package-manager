package main

import (
	"npm-tiny-package-manager/file"
)

func main() {
	root, err := file.ParsePackageJson()
	if err != nil {
		panic(err)
	}
	for k, v := range root.Dependencies {
		println(k, v)
	}
}

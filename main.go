package main

import (
	"npm-tiny-package-manager/file"
	"npm-tiny-package-manager/npm"
)

func main() {
	root, err := file.ParsePackageJson()
	if err != nil {
		panic(err)
	}
	for pkgName, ver := range root.Dependencies {
		nm, err := npm.FetchPackageManifest(pkgName)
		if err != nil {
			panic(err)
		}
		err = npm.InstallTarball(nm, ver)
		if err != nil {
			panic(err)
		}
	}
}

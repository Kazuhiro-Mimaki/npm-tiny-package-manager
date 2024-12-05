package main

import (
	"npm-tiny-package-manager/file"
	"npm-tiny-package-manager/npm"
	"npm-tiny-package-manager/utils"
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

		msv, err := npm.MaxSatisfyingVer(utils.MapKeysToSlice(nm.Versions), string(ver))
		if err != nil {
			panic(err)
		}

		err = npm.InstallTarball(nm, npm.Version(msv))
		if err != nil {
			panic(err)
		}
	}
}

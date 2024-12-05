package main

import (
	"npm-tiny-package-manager/file"
	"npm-tiny-package-manager/npm"
	"npm-tiny-package-manager/resolver"
)

func main() {
	root, err := file.ParsePackageJson()
	if err != nil {
		panic(err)
	}

	info := resolver.Info{
		TopLevel: make(map[npm.PackageName]resolver.TopLevel),
	}
	npmManifestCache := make(map[npm.PackageName]npm.NpmManifest)

	for pkgName, ver := range root.Dependencies {
		err = resolver.ResolveRecursively(pkgName, npm.Version(ver), info, npmManifestCache)
		if err != nil {
			panic(err)
		}

	}

	for pkgName, topLevel := range info.TopLevel {
		err := npm.InstallTarball(pkgName, topLevel.TarballUrl)
		if err != nil {
			panic(err)
		}
	}
}

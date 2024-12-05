package main

import (
	"npm-tiny-package-manager/file"
	"npm-tiny-package-manager/npm"
	"npm-tiny-package-manager/resolver"

	"golang.org/x/sync/errgroup"
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

	var eg errgroup.Group

	for pkgName, ver := range root.Dependencies {
		err = resolver.ResolveRecursively(pkgName, npm.Version(ver), info, npmManifestCache)
		if err != nil {
			panic(err)
		}
	}

	for pkgName, topLevel := range info.TopLevel {
		eg.Go(func() error {
			err := npm.InstallTarball(pkgName, topLevel.TarballUrl)
			if err != nil {
				return err
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		panic(err)
	}
}

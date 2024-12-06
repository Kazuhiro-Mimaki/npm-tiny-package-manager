package main

import (
	"npm-tiny-package-manager/file"
	"npm-tiny-package-manager/logger"
	"npm-tiny-package-manager/npm"
	"npm-tiny-package-manager/resolver"
	"npm-tiny-package-manager/types"

	"golang.org/x/sync/errgroup"
)

func main() {
	root, err := file.ParsePackageJson()
	if err != nil {
		panic(err)
	}

	info := resolver.Info{
		TopLevel: make(map[types.PackageName]resolver.TopLevel),
	}

	var eg errgroup.Group

	for pkgName, ver := range root.Dependencies {
		err = resolver.ResolveRecursively(pkgName, types.Version(ver), info)
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
		logger.InstalledLog(pkgName, topLevel.Version)
	}

	if err := eg.Wait(); err != nil {
		panic(err)
	}
}

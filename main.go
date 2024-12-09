package main

import (
	"fmt"

	"npm-tiny-package-manager/file"
	"npm-tiny-package-manager/lock"
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

	info := &resolver.Info{
		TopLevel:   make(map[types.PackageName]resolver.TopLevel),
		Conflicted: []resolver.Conflicted{},
	}

	err = lock.ReadLock()
	if err != nil {
		panic(err)
	}

	var eg errgroup.Group

	for pkgName, constraint := range root.Dependencies {
		dependencyStack := resolver.DependencyStack{Items: []resolver.DependencyStackItem{}}
		err = resolver.ResolveRecursively(pkgName, constraint, root.Dependencies, info, dependencyStack)
		if err != nil {
			panic(err)
		}
	}

	lock.SaveLock()

	for pkgName, item := range info.TopLevel {
		eg.Go(func() error {
			err := npm.InstallTarball(pkgName, item.Version, item.TarballUrl, ".")
			if err != nil {
				return err
			}
			return nil
		})
	}

	for _, item := range info.Conflicted {
		eg.Go(func() error {
			err := npm.InstallTarball(item.Name, item.Version, item.TarballUrl, fmt.Sprintf("./node_modules/%s", item.Parent))
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

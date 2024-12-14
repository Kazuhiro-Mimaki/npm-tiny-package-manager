package main

import (
	"fmt"
	"strings"

	"npm-tiny-package-manager/cli"
	"npm-tiny-package-manager/file"
	"npm-tiny-package-manager/lock"
	"npm-tiny-package-manager/npm"
	"npm-tiny-package-manager/resolver"
	"npm-tiny-package-manager/types"

	"golang.org/x/sync/errgroup"
)

func main() {
	installOptions, err := cli.Parse()
	if err != nil {
		panic(err)
	}

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

	for _, pkgName := range installOptions.Packages {
		constraint := "*"
		if strings.Contains(pkgName, "@") {
			v := pkgName[strings.LastIndex(pkgName, "@")+1:]
			if npm.IsValid(v) {
				constraint = v
				pkgName = pkgName[:strings.LastIndex(pkgName, "@")]
			}
		}
		if constraint == "*" {
			_, resolvedVer, err := resolver.ResolvePackage(types.PackageName(pkgName), types.Constraint(constraint))
			if err != nil {
				panic(err)
			}
			constraint = string(fmt.Sprintf("^%s", resolvedVer))
		}
		if installOptions.SaveDev {
			root.DevDependencies[types.PackageName(pkgName)] = types.Constraint(constraint)
		} else {
			root.Dependencies[types.PackageName(pkgName)] = types.Constraint(constraint)
		}
	}

	rootDependencies := collectDependencies(root, installOptions.Production)

	for pkgName, constraint := range rootDependencies {
		dependencyStack := resolver.DependencyStack{Items: []resolver.DependencyStackItem{}}
		err = resolver.ResolveRecursively(pkgName, constraint, rootDependencies, info, dependencyStack)
		if err != nil {
			panic(err)
		}
	}

	for pkgName, item := range info.TopLevel {
		eg.Go(func() error {
			err := npm.InstallTarball(pkgName, item.Version, item.TarballUrl, fmt.Sprintf("./node_modules/%s", pkgName))
			if err != nil {
				return err
			}
			return nil
		})
	}

	for _, item := range info.Conflicted {
		eg.Go(func() error {
			err := npm.InstallTarball(item.Name, item.Version, item.TarballUrl, fmt.Sprintf("./node_modules/%s/node_modules/%s", item.Parent, item.Name))
			if err != nil {
				return err
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		panic(err)
	}

	err = lock.SaveLock()
	if err != nil {
		panic(err)
	}
	err = file.WritePackageJson(root)
	if err != nil {
		panic(err)
	}
}

func collectDependencies(rootDependencies file.PackageJson, isPrd bool) types.Dependencies {
	allRootDependencies := make(map[types.PackageName]types.Constraint)

	if isPrd {
		return rootDependencies.Dependencies
	}

	for pkgName, constraint := range rootDependencies.Dependencies {
		allRootDependencies[pkgName] = constraint
	}
	for pkgName, constraint := range rootDependencies.DevDependencies {
		allRootDependencies[pkgName] = constraint
	}

	return allRootDependencies
}

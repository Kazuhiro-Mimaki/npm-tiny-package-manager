package resolver

import (
	"fmt"
	"strings"

	"npm-tiny-package-manager/lock"
	"npm-tiny-package-manager/logger"
	"npm-tiny-package-manager/npm"
	"npm-tiny-package-manager/types"
	"npm-tiny-package-manager/utils"
)

type Info struct {
	TopLevel   map[types.PackageName]TopLevel
	Conflicted []Conflicted
}

type TopLevel struct {
	TarballUrl string
	Version    types.Version
}

type Conflicted struct {
	Name       types.PackageName
	Version    types.Version
	TarballUrl string
	Parent     types.PackageName
}

/*
 * Resolve the dependencies recursively.
 */
func ResolveRecursively(
	pkgName types.PackageName,
	constraint types.Constraint,
	rootDependencies types.Dependencies,
	installList *Info,
	dependencyStack DependencyStack,
) error {
	matchedManifest, resolvedVer, err := ResolvePackage(pkgName, constraint)
	if err != nil {
		return err
	}
	topLevel, existsInTopLevel := installList.TopLevel[pkgName]
	rootDependencyConstraint, existsInRootDependency := rootDependencies[pkgName]

	if !existsInTopLevel && len(dependencyStack.Items) == 0 {
		installList.TopLevel[pkgName] = TopLevel{
			TarballUrl: matchedManifest.Dist.Tarball,
			Version:    types.Version(resolvedVer),
		}
	} else if !existsInTopLevel && !existsInRootDependency {
		installList.TopLevel[pkgName] = TopLevel{
			TarballUrl: matchedManifest.Dist.Tarball,
			Version:    types.Version(resolvedVer),
		}
	} else if !existsInTopLevel && existsInRootDependency {
		_, resolvedVer, err := ResolvePackage(pkgName, rootDependencyConstraint)
		if err != nil {
			return err
		}
		if !utils.Satisfies(string(resolvedVer), string(constraint)) {
			logger.ConflictLog(pkgName, constraint, types.Version(rootDependencyConstraint))
			installList.Conflicted = append(installList.Conflicted, Conflicted{
				Name:       pkgName,
				TarballUrl: matchedManifest.Dist.Tarball,
				Parent:     dependencyStack.Items[len(dependencyStack.Items)-1].Name,
				Version:    types.Version(constraint),
			})
		}
	} else {
		if !utils.Satisfies(string(topLevel.Version), string(constraint)) {
			logger.ConflictLog(pkgName, constraint, types.Version(rootDependencyConstraint))
			installList.Conflicted = append(installList.Conflicted, Conflicted{
				Name:       pkgName,
				TarballUrl: matchedManifest.Dist.Tarball,
				Parent:     dependencyStack.Items[len(dependencyStack.Items)-1].Name,
				Version:    types.Version(constraint),
			})
		}
	}

	lock.UpsertLock(lock.LockKey(fmt.Sprintf("%s@%s", pkgName, constraint)), lock.Lock{
		Version:      resolvedVer,
		Shasum:       matchedManifest.Dist.Shasum,
		Url:          matchedManifest.Dist.Tarball,
		Dependencies: matchedManifest.Dependencies,
	})

	filteredDependencies := make(types.Dependencies)
	for depName, depConstraint := range matchedManifest.Dependencies {
		if !hasCirculation(depName, depConstraint, dependencyStack) {
			filteredDependencies[depName] = depConstraint
		}
	}

	if len(filteredDependencies) > 0 {
		dependencyStack.append(
			DependencyStackItem{
				Name:         pkgName,
				Version:      resolvedVer,
				Dependencies: filteredDependencies,
			},
		)
		for depName, depConstraint := range filteredDependencies {
			err = ResolveRecursively(depName, depConstraint, rootDependencies, installList, dependencyStack)
			if err != nil {
				return err
			}
		}
		dependencyStack.pop()
	}

	return nil
}

/*
* This function is to resolve the package.
* If the package is not in the lock file, fetch the manifest from utils.
* Resolve the semantic version.
 */
func ResolvePackage(pkgName types.PackageName, constraint types.Constraint) (types.Manifest, types.Version, error) {
	//  Get package manifest from lock
	manifestVersions, existsInLock := lock.GetItem(pkgName, constraint)

	// If the package is not in the lock file, fetch the manifest from npm
	if !existsInLock {
		manifest, err := npm.FetchPackageManifest(pkgName)
		if err != nil {
			return types.Manifest{}, "", err
		}
		manifestVersions = manifest.Versions
	}

	// Resolve semantic version
	resolvedVer, err := utils.MaxSatisfyingVer(convertManifestVersions(manifestVersions), string(constraint))
	ver := types.Version(resolvedVer)
	if err != nil {
		return types.Manifest{}, "", err
	}

	logger.ResolveLog(pkgName, constraint, ver)

	return manifestVersions[ver], ver, nil
}

func convertManifestVersions(manifestVersions map[types.Version]types.Manifest) []string {
	converter := func(a types.Version) string { return string(a) }
	return utils.ConvertSliceType(utils.MapKeysToSlice(manifestVersions), converter)
}

/*
* This function is to check if there is dependency circulation.
*
* If a package is existed in the stack and it satisfy the semantic version,
* it turns out that there is dependency circulation.
 */
func hasCirculation(
	pkgName types.PackageName,
	constraint types.Constraint,
	stack DependencyStack,
) bool {
	for _, depStack := range stack.Items {
		if depStack.Name == pkgName && utils.Satisfies(string(depStack.Version), string(constraint)) {
			return true
		}
	}
	return false
}

/*
* This function is to resolve the installed packages.
 */
func ResolveInstalledPackages(pkgName string) (string, error) {
	constraint := "*"
	if strings.Contains(pkgName, "@") {
		v := pkgName[strings.LastIndex(pkgName, "@")+1:]
		if utils.IsValid(v) {
			constraint = v
			pkgName = pkgName[:strings.LastIndex(pkgName, "@")]
		}
	}
	if constraint == "*" {
		_, resolvedVer, err := ResolvePackage(types.PackageName(pkgName), types.Constraint(constraint))
		if err != nil {
			return "", err
		}
		constraint = string(fmt.Sprintf("^%s", resolvedVer))
	}
	return constraint, nil
}

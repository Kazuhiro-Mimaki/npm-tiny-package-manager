package resolver

import (
	"fmt"

	"npm-tiny-package-manager/lock"
	"npm-tiny-package-manager/logger"
	"npm-tiny-package-manager/npm"
	"npm-tiny-package-manager/types"
	"npm-tiny-package-manager/utils"
)

type Info struct {
	TopLevel map[types.PackageName]TopLevel
}

type TopLevel struct {
	TarballUrl string
	Version    types.Version
}

func ResolveRecursively(pkgName types.PackageName, constraint types.Constraint, installList Info) error {
	/**
	 * Get package manifest from lock
	 */
	manifestVersions, ok := lock.GetItem(pkgName, constraint)

	/**
	 * If the package is not in the lock file, fetch the manifest from npm
	 */
	if !ok {
		manifest, err := npm.FetchPackageManifest(pkgName)
		if err != nil {
			return err
		}
		manifestVersions = manifest.Versions
	}

	/**
	 * Resolve semantic version
	 */
	maxVersion, err := npm.MaxSatisfyingVer(utils.MapKeysToSlice(manifestVersions), constraint)
	if err != nil {
		return err
	}

	logger.ResolveLog(pkgName, constraint, maxVersion)

	matchedManifest := manifestVersions[maxVersion]

	installList.TopLevel[pkgName] = TopLevel{
		TarballUrl: matchedManifest.Dist.Tarball,
		Version:    types.Version(maxVersion),
	}

	lock.UpsertLock(lock.LockKey(fmt.Sprintf("%s@%s", pkgName, constraint)), lock.Lock{
		Version:      maxVersion,
		Shasum:       matchedManifest.Dist.Shasum,
		Url:          matchedManifest.Dist.Tarball,
		Dependencies: matchedManifest.Dependencies,
	})

	if len(matchedManifest.Dependencies) > 0 {
		for depName, depConstraint := range matchedManifest.Dependencies {
			err = ResolveRecursively(depName, depConstraint, installList)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

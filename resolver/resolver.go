package resolver

import (
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

func ResolveRecursively(pkgName types.PackageName, constraint types.Version, installList Info) error {
	manifest, err := npm.FetchPackageManifest(pkgName)
	if err != nil {
		return err
	}

	maxVersion, err := npm.MaxSatisfyingVer(utils.MapKeysToSlice(manifest.Versions), string(constraint))
	if err != nil {
		return err
	}

	logger.ResolveLog(pkgName, constraint, maxVersion)

	matchedManifest := manifest.Versions[maxVersion]

	installList.TopLevel[pkgName] = TopLevel{
		TarballUrl: matchedManifest.Dist.Tarball,
		Version:    types.Version(maxVersion),
	}

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

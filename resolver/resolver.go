package resolver

import (
	"fmt"

	"npm-tiny-package-manager/npm"
	"npm-tiny-package-manager/utils"
)

type Info struct {
	TopLevel map[npm.PackageName]TopLevel
}

type TopLevel struct {
	TarballUrl string
	Version    npm.Version
}

func ResolveRecursively(pkgName npm.PackageName, constraint npm.Version, installList Info) error {
	manifest, err := npm.FetchPackageManifest(pkgName)
	if err != nil {
		return err
	}

	maxVersion, err := npm.MaxSatisfyingVer(utils.MapKeysToSlice(manifest.Versions), string(constraint))
	if err != nil {
		return err
	}

	fmt.Printf("Resolving %s@%s => %s\n", pkgName, constraint, maxVersion)

	matchedManifest := manifest.Versions[maxVersion]

	installList.TopLevel[pkgName] = TopLevel{
		TarballUrl: matchedManifest.Dist.Tarball,
		Version:    npm.Version(maxVersion),
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

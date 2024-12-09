package npm

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"npm-tiny-package-manager/logger"
	"npm-tiny-package-manager/types"

	"github.com/Masterminds/semver/v3"
)

const REGISTRY = "https://registry.npmjs.org"

var CACHE = make(types.NpmManifestCache)

/*
 * Fetch the package manifest from the npm registry
 */
func FetchPackageManifest(name types.PackageName) (types.NpmManifest, error) {
	if v, ok := CACHE[name]; ok {
		return v, nil
	}
	resp, err := http.Get(fmt.Sprintf("%s/%s", REGISTRY, name))
	if err != nil {
		return types.NpmManifest{}, err
	}
	defer resp.Body.Close()
	var nm types.NpmManifest
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return types.NpmManifest{}, err
	}
	if err := json.Unmarshal(body, &nm); err != nil {
		return types.NpmManifest{}, err
	}
	CACHE[name] = nm
	return nm, nil
}

/*
 * Install the package tarball
 */
func InstallTarball(pkgName types.PackageName, version types.Version, tarballUrl, location string) error {
	var err error

	resp, err := http.Get(tarballUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	gzf, err := gzip.NewReader(resp.Body)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(gzf)

	/*
	 * Extract the tarball to the node_modules directory
	 */
	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		path := fmt.Sprintf("%s/node_modules/%s/%s", location, pkgName, strings.Replace(header.Name, "package/", "", 1))

		err = install(tarReader, path)
		if err != nil {
			return err
		}
	}

	logger.InstalledLog(pkgName, version)

	return nil
}

/*
 * Install the package tarball
 */
func install(tarReader *tar.Reader, path string) error {
	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return err
	}

	out, err := os.Create(path)
	out.Chmod(os.ModePerm)
	if err != nil {
		return err
	}
	defer out.Close()

	io.Copy(out, tarReader)

	return nil
}

/*
 * Get the latest version of a package
 */
func MaxSatisfyingVer(versions []types.Version, constraint types.Constraint) (types.Version, error) {
	c, err := semver.NewConstraint(string(constraint))
	if err != nil {
		return "", err
	}

	var maxVersion *semver.Version

	for i := 0; i < len(versions); i++ {
		v, err := semver.NewVersion(string(versions[i]))
		if err != nil {
			return "", err
		}
		if c.Check(v) {
			if maxVersion == nil || v.GreaterThan(maxVersion) {
				maxVersion = v
			}
		}
	}

	return types.Version(maxVersion.String()), nil
}

/*
 * Check if a version satisfies a constraint
 */
func Satisfies(version types.Version, constraint types.Constraint) bool {
	c, err := semver.NewConstraint(string(constraint))
	if err != nil {
		return false
	}

	v, err := semver.NewVersion(string(version))
	if err != nil {
		return false
	}

	return c.Check(v)
}

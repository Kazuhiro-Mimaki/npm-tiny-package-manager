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

	"github.com/Masterminds/semver/v3"
)

type (
	PackageName string
	Version     string
)

type (
	Dependencies map[PackageName]Version
)

type PackageJson struct {
	Dependencies    Dependencies
	DevDependencies Dependencies
}

type Dist struct {
	Shasum  string
	Tarball string
}

type Manifest struct {
	Dependencies Dependencies
	Dist         Dist
}

type NpmManifest struct {
	Name     string
	Versions map[Version]Manifest
}

const REGISTRY = "https://registry.npmjs.org"

func FetchPackageManifest(name PackageName) (NpmManifest, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%s", REGISTRY, name))
	if err != nil {
		return NpmManifest{}, err
	}
	defer resp.Body.Close()
	var nm NpmManifest
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return NpmManifest{}, err
	}
	if err := json.Unmarshal(body, &nm); err != nil {
		return NpmManifest{}, err
	}
	return nm, nil
}

func InstallTarball(npmManifest NpmManifest, version Version) error {
	var err error

	url := npmManifest.Versions[version].Dist.Tarball
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	gzf, err := gzip.NewReader(resp.Body)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(gzf)

	/**
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

		path := fmt.Sprintf("./node_modules/%s/%s", npmManifest.Name, strings.Replace(header.Name, "package/", "", 1))

		err = install(tarReader, path)
		if err != nil {
			return err
		}
	}

	return nil
}

func install(tarReader *tar.Reader, path string) error {
	err := os.MkdirAll(filepath.Dir(path), 0o755)
	if err != nil {
		return err
	}

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	io.Copy(out, tarReader)

	return nil
}

func MaxSatisfyingVer(versions []Version, constraint string) (string, error) {
	c, err := semver.NewConstraint(constraint)
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

	return maxVersion.String(), nil
}

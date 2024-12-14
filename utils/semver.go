package utils

import "github.com/Masterminds/semver/v3"

/*
 * Get the latest version of a package
 */
func MaxSatisfyingVer(versions []string, constraint string) (string, error) {
	c, err := semver.NewConstraint(constraint)
	if err != nil {
		return "", err
	}

	var maxVersion *semver.Version

	for i := 0; i < len(versions); i++ {
		v, err := semver.NewVersion(versions[i])
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

/*
 * Check if a version satisfies a constraint
 */
func Satisfies(version string, constraint string) bool {
	c, err := semver.NewConstraint(constraint)
	if err != nil {
		return false
	}

	v, err := semver.NewVersion(version)
	if err != nil {
		return false
	}

	return c.Check(v)
}

/*
 * Check if a version is valid
 */
func IsValid(version string) bool {
	_, err := semver.NewConstraint(version)
	return err == nil
}

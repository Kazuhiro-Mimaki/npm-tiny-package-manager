package package_json

import (
	"encoding/json"
	"io"
	"os"

	"npm-tiny-package-manager/types"
)

type PackageJson struct {
	Name            string
	Version         string
	Description     string
	Dependencies    types.Dependencies
	DevDependencies types.Dependencies
}

const PATH = "package.json"

/*
 * Parse the package.json file
 */
func ParsePackageJson() (PackageJson, error) {
	jsonFile, err := os.Open(PATH)
	if err != nil {
		return PackageJson{}, err
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	pkj := PackageJson{
		Name:            "",
		Version:         "",
		Description:     "",
		Dependencies:    make(map[types.PackageName]types.Constraint),
		DevDependencies: make(map[types.PackageName]types.Constraint),
	}
	json.Unmarshal([]byte(byteValue), &pkj)

	return pkj, nil
}

/*
 * Write the package.json file
 */
func WritePackageJson(pkj PackageJson) error {
	byteValue, err := json.Marshal(pkj)
	if err != nil {
		return err
	}

	err = os.WriteFile(PATH, byteValue, 0o644)
	if err != nil {
		return err
	}

	return nil
}

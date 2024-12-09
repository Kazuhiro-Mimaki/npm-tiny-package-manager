package file

import (
	"encoding/json"
	"io"
	"os"

	"npm-tiny-package-manager/types"
)

const PATH = "package.json"

/*
 * Parse the package.json file
 */
func ParsePackageJson() (types.PackageJson, error) {
	jsonFile, err := os.Open(PATH)
	if err != nil {
		return types.PackageJson{}, err
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var pkj types.PackageJson
	json.Unmarshal([]byte(byteValue), &pkj)

	return pkj, nil
}

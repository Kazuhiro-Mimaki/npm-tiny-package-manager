package file

import (
	"encoding/json"
	"io"
	"os"

	"npm-tiny-package-manager/npm"
)

const PATH = "package.json"

func ParsePackageJson() (npm.PackageJson, error) {
	jsonFile, err := os.Open(PATH)
	if err != nil {
		return npm.PackageJson{}, err
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var pkj npm.PackageJson
	json.Unmarshal([]byte(byteValue), &pkj)

	return pkj, nil
}

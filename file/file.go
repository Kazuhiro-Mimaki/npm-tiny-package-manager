package file

import (
	"encoding/json"
	"io"
	"os"
)

type (
	Dependencies map[string]string
)

type PackageJson struct {
	Dependencies    Dependencies
	DevDependencies Dependencies
}

const path = "package.json"

func ParsePackageJson() (PackageJson, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return PackageJson{}, err
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var pkj PackageJson
	json.Unmarshal([]byte(byteValue), &pkj)

	return pkj, nil
}

package types

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

type NpmManifestCache map[PackageName]NpmManifest

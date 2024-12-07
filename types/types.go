package types

type (
	PackageName string
	Version     string
	Constraint  string
)

type (
	Dependencies map[PackageName]Constraint
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
	Dist         Dist
	Dependencies Dependencies
}

type NpmManifest struct {
	Name     string
	Versions map[Version]Manifest
}

type NpmManifestCache map[PackageName]NpmManifest

package libraryinfo

type PackageInfo struct {
	Name          string `json:"name"`
	Version       string `json:"version"`
	SourcePackage string `json:"source_package"`
}

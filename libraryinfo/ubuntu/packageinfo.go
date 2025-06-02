package ubuntu

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/jimbersoftware/pra_client/libraryinfo"
)

// PackageInfo holds details about an installed package, including its source package

// getSourcePackage resolves the source package name for a given binary package
func getSourcePackage(pkg string) string {
	cmd := exec.Command("apt-cache", "show", pkg)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "" // fallback to empty string if it fails
	}

	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Source:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "Source:"))
		}
	}
	return pkg // fallback to the binary name if Source field not found
}

// GetInstalledPackages retrieves a list of installed packages, their versions, and source packages
func GetInstalledPackages() ([]libraryinfo.PackageInfo, error) {
	cmd := exec.Command("dpkg-query", "-W", "-f=${Package} ${Version}\\n")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(out.String(), "\n")
	var packages []libraryinfo.PackageInfo

	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			name := parts[0]
			version := parts[1]
			source := getSourcePackage(name)

			packages = append(packages, libraryinfo.PackageInfo{
				Name:          name,
				Version:       version,
				SourcePackage: source,
			})
		}
	}

	return packages, nil
}

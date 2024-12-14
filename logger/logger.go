package logger

import (
	"fmt"
	"log/slog"

	"npm-tiny-package-manager/types"
)

/*
 * Log the package resolution
 */
func ResolveLog(pkgName types.PackageName, constraint types.Constraint, maxVersion types.Version) {
	slog.Info(fmt.Sprintf("Resolving %s@%s => %s\n", pkgName, constraint, maxVersion))
}

/*
 * Log the package installation
 */
func InstalledLog(pkgName types.PackageName, version types.Version) {
	slog.Info(fmt.Sprintf("%s@%s Installed\n", pkgName, version))
}

/*
 * Log the package conflict
 */
func ConflictLog(pkgName types.PackageName, constraint types.Constraint, conflictVersion types.Version) {
	slog.Info(fmt.Sprintf("Conflict %s@%s => %s\n", pkgName, constraint, conflictVersion))
}

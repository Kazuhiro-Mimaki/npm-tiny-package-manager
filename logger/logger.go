package logger

import (
	"fmt"
	"log/slog"

	"npm-tiny-package-manager/types"
)

func ResolveLog(pkgName types.PackageName, constraint types.Constraint, maxVersion types.Version) {
	slog.Info(fmt.Sprintf("Resolving %s@%s => %s\n", pkgName, constraint, maxVersion))
}

func InstalledLog(pkgName types.PackageName, version types.Version) {
	slog.Info(fmt.Sprintf("%s@%s Installed\n", pkgName, version))
}

func ConflictLog(pkgName types.PackageName, constraint types.Constraint, conflictVersion types.Version) {
	slog.Info(fmt.Sprintf("Conflict %s@%s => %s\n", pkgName, constraint, conflictVersion))
}

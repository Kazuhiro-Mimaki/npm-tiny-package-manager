package lock

import (
	"encoding/json"
	"fmt"
	"os"

	"npm-tiny-package-manager/types"
)

type Lock struct {
	Version      types.Version
	Url          string
	Shasum       string
	Dependencies types.Dependencies
}

type (
	// `packageName@version` is the key
	LockKey string
	LockMap map[LockKey]Lock
)

const PATH = "package-lock.json"

/*
 * The old lock is only for reading from the lock file,
 * so the old lock should be read only except reading the lock file.
 */
var OldLock LockMap

/*
 * The new lock is only for writing to the lock file,
 * so the new lock should be written only except saving the lock file.
 */
var NewLock LockMap

/*
 * Simply read the lock file.
 * Skip it if we cannot find the lock file.
 */
func ReadLock() error {
	_, err := os.Stat(PATH)
	if os.IsNotExist(err) {
		return nil
	}
	r, err := os.Open(PATH)
	if err != nil {
		return err
	}
	defer r.Close()

	err = json.NewDecoder(r).Decode(&OldLock)
	if err != nil {
		return err
	}
	return nil
}

/*
 * Retrieve the information of a package by name and it's semantic
 * version range.
 *
 * Note that we don't return the data directly.
 * That is, we just do format the data,
 * which make the data structure similar to npm registry.
 *
 * This can let us avoid changing the logic of the `collectDeps`
 * function in the `list` module.
 */
func GetItem(pkgName types.PackageName, constraint types.Constraint) (map[types.Version]types.Manifest, bool) {
	res := map[types.Version]types.Manifest{}
	item, ok := OldLock[LockKey(fmt.Sprintf("%s@%s", pkgName, constraint))]
	if !ok {
		return res, false
	}
	res[item.Version] = types.Manifest{
		Dist: types.Dist{
			Shasum:  item.Shasum,
			Tarball: item.Url,
		},
		Dependencies: item.Dependencies,
	}

	return res, true
}

/*
 * Save the information of a package to the lock.
 * If that information is not existed in the lock, create it.
 * Otherwise, just update it.
 */
func UpsertLock(lockKey LockKey, lock Lock) {
	if NewLock == nil {
		NewLock = LockMap{}
	}
	NewLock[lockKey] = lock
}

/*
 * Simply save the lock file.
 */
func SaveLock() error {
	w, err := json.Marshal(NewLock)
	if err != nil {
		return err
	}
	err = os.WriteFile(PATH, w, 0o644)
	if err != nil {
		return err
	}
	return nil
}

package lecture

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func CheckExistDir(path string) error {
	if d, err := os.Stat(path); os.IsNotExist(err) || d.IsDir() {
		return err
	}
	return nil
}

func GetFileLists(inputPath string) ([]string, error) {
	findList := []string{}

	err := filepath.WalkDir(inputPath, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return errors.Wrap(err, "failed filepath.WalkDir")
		}

		if info.IsDir() {
			return nil
		}

		findList = append(findList, path)
		return nil
	})
	return findList, err
}

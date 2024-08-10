package filesystem

import (
	"os"
	"path/filepath"
)

func DirectorySize(path string) (int64, error) {
	var size int64

	iterator := func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			size += info.Size()
		}

		return err
	}

	err := filepath.Walk(path, iterator)
	if err != nil {
		return 0, err
	}

	return size, nil
}

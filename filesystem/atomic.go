package filesystem

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
)

// Adapted from: https://github.com/containerd/continuity/blob/main/ioutils.go under Apache License

/*
   Copyright The containerd Authors.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

const defaultUmask = 0o077

var Mask = os.FileMode(defaultUmask) //nolint:gochecknoglobals

// WriteFile atomically writes data to a file by first writing to a temp file and calling rename.
func WriteFile(filename string, data []byte, perm os.FileMode) error {
	reader := bytes.NewBuffer(data)
	dataSize := int64(len(data))
	perm = (^Mask) & perm

	tmpFile, err := os.CreateTemp(filepath.Dir(filename), ".tmp-"+filepath.Base(filename))
	if err != nil {
		return err
	}

	err = os.Chmod(tmpFile.Name(), perm)
	if err != nil {
		tmpFile.Close()

		return err
	}

	n, err := io.Copy(tmpFile, reader)
	if err == nil && n < dataSize {
		tmpFile.Close()

		return io.ErrShortWrite
	}

	if err != nil {
		tmpFile.Close()

		return err
	}

	if err := tmpFile.Sync(); err != nil {
		tmpFile.Close()

		return err
	}

	if err := tmpFile.Close(); err != nil {
		return err
	}

	return os.Rename(tmpFile.Name(), filename)
}

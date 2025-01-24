/*
   Copyright Farcloser.

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

package store

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"os"

	"github.com/peterbourgon/diskv/v3"

	"go.farcloser.world/core/filesystem"
)

const transformBlockSize = 64 // grouping of chars per directory depth

func transform(key string) []string {
	var (
		sliceSize = len(key) / transformBlockSize
		pathSlice = make([]string, sliceSize)
	)

	for i := range sliceSize {
		from, to := i*transformBlockSize, (i+1)*transformBlockSize
		pathSlice[i] = key[from:to]
	}

	return pathSlice
}

func hash(s string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
}

type Options struct {
	Path      string
	CacheSize int64
}

func New(options *Options) *Store {
	path := options.Path
	if path == "" {
		path = defaultStoreDir
	}

	var cacheSize uint64
	//nolint:gocritic
	if options.CacheSize == 0 {
		cacheSize = defaultCacheSize
	} else if options.CacheSize < 0 {
		cacheSize = 0
	} else {
		cacheSize = uint64(options.CacheSize)
	}

	return &Store{
		diskv: diskv.New(diskv.Options{
			BasePath:     path,
			Transform:    transform,
			CacheSizeMax: cacheSize,
			PathPerm:     filesystem.DirPermissionsPrivate,
			FilePerm:     filesystem.FilePermissionsPrivate,
			TempDir:      "", // "FIXME TO GET ATOMIC WRITES",
		}),
	}
}

type Store struct {
	diskv *diskv.Diskv
	lock  *os.File
}

func (st *Store) Read(name string) (content []byte, err error) {
	if st.lock == nil {
		err = st.ReadOnlyLock()
		if err != nil {
			return nil, err
		}

		defer func() {
			if unlockErr := st.Unlock(); unlockErr != nil {
				err = errors.Join(err, unlockErr)
			}
		}()
	}

	content, err = st.diskv.Read(hash(name))
	if err != nil {
		err = errors.Join(ErrFileStoreFail, err)
	}

	return content, err
}

func (st *Store) ReadFromKey(key string) (content []byte, err error) {
	if st.lock == nil {
		err = st.ReadOnlyLock()
		if err != nil {
			return nil, err
		}

		defer func() {
			if unlockErr := st.Unlock(); unlockErr != nil {
				err = errors.Join(err, unlockErr)
			}
		}()
	}

	content, err = st.diskv.Read(key)
	if err != nil {
		err = errors.Join(ErrFileStoreFail, err)
	}

	return content, err
}

func (st *Store) Has(name string) (has bool, err error) {
	if st.lock == nil {
		err = st.ReadOnlyLock()
		if err != nil {
			return false, err
		}

		defer func() {
			if unlockErr := st.Unlock(); unlockErr != nil {
				err = errors.Join(err, unlockErr)
			}
		}()
	}

	return st.diskv.Has(hash(name)), nil
}

func (st *Store) Keys() <-chan string {
	return st.diskv.Keys(nil)
}

func (st *Store) Digest(name string) string {
	return hash(name)
}

func (st *Store) Write(name string, value []byte) (err error) {
	if st.lock == nil {
		err = st.WriteLock()
		if err != nil {
			return err
		}

		defer func() {
			if unlockErr := st.Unlock(); unlockErr != nil {
				err = errors.Join(err, unlockErr)
			}
		}()
	}

	err = st.diskv.Write(hash(name), value)
	if err != nil {
		err = errors.Join(ErrFileStoreFail, err)
	}

	return err
}

func (st *Store) Delete(name string) (err error) {
	if st.lock == nil {
		err = st.WriteLock()
		if err != nil {
			return err
		}

		defer func() {
			if unlockErr := st.Unlock(); unlockErr != nil {
				err = errors.Join(err, unlockErr)
			}
		}()
	}

	err = st.diskv.Erase(hash(name))
	if err != nil {
		err = errors.Join(ErrFileStoreFail, err)
	}

	return err
}

func (st *Store) Rename(oldName, newName string) (err error) {
	if st.lock == nil {
		err = st.WriteLock()
		if err != nil {
			return err
		}

		defer func() {
			if unlockErr := st.Unlock(); unlockErr != nil {
				err = errors.Join(err, unlockErr)
			}
		}()
	}

	var content []byte

	content, err = st.Read(oldName)
	if err != nil {
		return err
	}

	err = st.Write(newName, content)
	if err != nil {
		return err
	}

	return st.Delete(oldName)
}

// Lock by default gets an exclusive read/write lock.
func (st *Store) Lock() (err error) {
	return st.WriteLock()
}

func (st *Store) WriteLock() (err error) {
	err = os.MkdirAll(st.diskv.BasePath, filesystem.DirPermissionsPrivate)
	if err != nil {
		return err
	}

	lock, err := filesystem.Lock(st.diskv.BasePath)
	if err != nil {
		err = errors.Join(ErrFileStoreFail, err)
	} else {
		st.lock = lock
	}

	return err
}

func (st *Store) ReadOnlyLock() (err error) {
	err = os.MkdirAll(st.diskv.BasePath, filesystem.DirPermissionsPrivate)
	if err != nil {
		return err
	}

	lock, err := filesystem.ReadOnlyLock(st.diskv.BasePath)
	if err != nil {
		err = errors.Join(ErrFileStoreFail, err)
	} else {
		st.lock = lock
	}

	return err
}

func (st *Store) Unlock() (err error) {
	defer func() {
		if err != nil {
			err = errors.Join(ErrFileStoreFail, err)
		}
	}()

	err = filesystem.Unlock(st.lock)
	if err != nil {
		err = errors.Join(ErrFileStoreFail, err)
	} else {
		st.lock = nil
	}

	return err
}

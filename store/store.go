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

func (st *Store) Read(name string) ([]byte, error) {
	if st.lock == nil {
		err := st.ReadOnlyLock()
		if err != nil {
			return nil, err
		}

		defer func() {
			_ = st.Unlock()
		}()
	}

	content, err := st.diskv.Read(hash(name))
	if err != nil {
		err = errors.Join(ErrFileStoreFail, err)
	}

	return content, err
}

func (st *Store) ReadFromKey(key string) ([]byte, error) {
	if st.lock == nil {
		err := st.ReadOnlyLock()
		if err != nil {
			return nil, err
		}

		defer func() {
			_ = st.Unlock()
		}()
	}

	content, err := st.diskv.Read(key)
	if err != nil {
		err = errors.Join(ErrFileStoreFail, err)
	}

	return content, err
}

func (st *Store) Has(name string) (bool, error) {
	if st.lock == nil {
		err := st.ReadOnlyLock()
		if err != nil {
			return false, err
		}

		defer func() {
			_ = st.Unlock()
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

func (st *Store) Write(name string, value []byte) error {
	if st.lock == nil {
		err := st.WriteLock()
		if err != nil {
			return err
		}

		defer func() {
			_ = st.Unlock()
		}()
	}

	err := st.diskv.Write(hash(name), value)
	if err != nil {
		err = errors.Join(ErrFileStoreFail, err)
	}

	return err
}

func (st *Store) Delete(name string) error {
	if st.lock == nil {
		err := st.WriteLock()
		if err != nil {
			return err
		}

		defer func() {
			_ = st.Unlock()
		}()
	}

	err := st.diskv.Erase(hash(name))
	if err != nil {
		err = errors.Join(ErrFileStoreFail, err)
	}

	return err
}

func (st *Store) Rename(oldName, newName string) error {
	if st.lock == nil {
		err := st.WriteLock()
		if err != nil {
			return err
		}

		defer func() {
			_ = st.Unlock()
		}()
	}

	content, err := st.Read(oldName)
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
func (st *Store) Lock() error {
	return st.WriteLock()
}

func (st *Store) WriteLock() error {
	var err error
	st.lock, err = filesystem.Lock(st.diskv.BasePath)
	if err != nil {
		err = errors.Join(ErrFileStoreFail, err)
	}

	return err
}

func (st *Store) ReadOnlyLock() error {
	var err error

	st.lock, err = filesystem.ReadOnlyLock(st.diskv.BasePath)
	if err != nil {
		err = errors.Join(ErrFileStoreFail, err)
	}

	return err
}

func (st *Store) Unlock() error {
	err := filesystem.Unlock(st.lock)
	if err != nil {
		err = errors.Join(ErrFileStoreFail, err)
	}

	return err
}
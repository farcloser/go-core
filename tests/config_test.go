package tests_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"testing"

	"go.farcloser.world/core/config"
)

func TestConfigLoadTargetDoesNotExist(t *testing.T) {
	dir, _ := os.UserHomeDir()
	conf := config.New(dir, "does", "not", "exist")
	err := config.Load(conf)

	if err == nil || !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("should have returned fs.ErrNotExist: %s", err)
	}
}

func TestConfigLoadTargetIsADirectory(t *testing.T) {
	dir, _ := os.UserHomeDir()
	conf := config.New(dir, ".")
	err := config.Load(conf)
	//	t.Fatalf("should have returned fs.PathError: %s", err)

	var pe *fs.PathError
	if err == nil || !errors.As(err, &pe) {
		t.Fatalf("should have returned fs.PathError: %s", err)
	}
}

func TestConfigLoadTargetUnreadable(t *testing.T) {
	dir, _ := os.UserHomeDir()
	filename := path.Join(dir, "testunreadablefile")

	tmpFile, err := os.CreateTemp(filepath.Dir(filename), ".tmp-"+filepath.Base(filename))
	if err != nil {
		t.Fatalf("unexpected failure! %s", err)
	}

	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			panic(fmt.Sprintf("failed cleaning up %s\n", name))
		}
	}(tmpFile.Name())

	err = os.Chmod(tmpFile.Name(), 0o000)
	if err != nil {
		t.Fatalf("unexpected failure! %s", err)
	}

	err = tmpFile.Close()
	if err != nil {
		t.Fatalf("unexpected failure! %s", err)
	}

	d, f := path.Split(tmpFile.Name())
	conf := config.New(d, f)

	err = config.Load(conf)
	if err == nil || !errors.Is(err, fs.ErrPermission) {
		t.Fatalf("should have returned fs.ErrPermission: %s", err)
	}
}

func TestConfigLoadIsNotJSON(t *testing.T) {
	dir, _ := os.UserHomeDir()
	filename := path.Join(dir, "testnotjson")

	tmpFile, err := os.CreateTemp(filepath.Dir(filename), ".tmp-"+filepath.Base(filename))
	if err != nil {
		t.Fatalf("unexpected failure! %s", err)
	}

	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			panic(fmt.Sprintf("failed cleaning up %s\n", name))
		}
	}(tmpFile.Name())

	err = tmpFile.Close()
	if err != nil {
		t.Fatalf("unexpected failure! %s", err)
	}

	d, f := path.Split(tmpFile.Name())
	conf := config.New(d, f)

	var pe *json.SyntaxError

	err = config.Load(conf)
	if err == nil || !errors.As(err, &pe) {
		t.Fatalf("should have returned json.SyntaxError: %s", err)
	}
}

func TestConfigLoadWrongType(t *testing.T) {
	dir, _ := os.UserHomeDir()
	filename := path.Join(dir, "wrongtype")

	tmpFile, err := os.CreateTemp(filepath.Dir(filename), ".tmp-"+filepath.Base(filename))
	if err != nil {
		t.Fatalf("unexpected failure! %s", err)
	}

	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			panic(fmt.Sprintf("failed cleaning up %s\n", name))
		}
	}(tmpFile.Name())

	_, err = io.Copy(tmpFile, bytes.NewBufferString("{\"umask\": \"foobar\"}"))
	if err != nil {
		t.Fatalf("unexpected failure! %s", err)
	}

	err = tmpFile.Close()
	if err != nil {
		t.Fatalf("unexpected failure! %s", err)
	}

	d, f := path.Split(tmpFile.Name())
	conf := config.New(d, f)

	err = config.Load(conf)

	var pe *json.UnmarshalTypeError
	if err == nil || !errors.As(err, &pe) {
		t.Fatalf("should have returned json.SyntaxError: %s", err)
	}
}

func TestConfigResolve(t *testing.T) {
	conf := config.New("/somewhere", "thing.foo")
	l := conf.Resolve("/", "perdita")
	t.Fatalf("should have returned shit: %s", l)
}

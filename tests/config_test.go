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

	"go.codecomet.dev/core/config"
)

func TestConfigLoadTargetDoesNotExist(t *testing.T) {
	dir, _ := os.UserHomeDir()
	conf := config.New(false, "", path.Join(dir, "does", "not", "exist"))
	err := conf.Load()

	if err == nil || !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("should have returned fs.ErrNotExist: %s", err)
	}
}

func TestConfigLoadTargetIsADirectory(t *testing.T) {
	dir, _ := os.UserHomeDir()
	conf := config.New(false, "", dir)
	err := conf.Load()
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

	conf := config.New(false, "", tmpFile.Name())

	err = conf.Load()
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

	conf := config.New(false, "", tmpFile.Name())

	var pe *json.SyntaxError

	err = conf.Load()
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

	_, err = io.Copy(tmpFile, bytes.NewBuffer([]byte("{\"umask\": \"foobar\"}")))
	if err != nil {
		t.Fatalf("unexpected failure! %s", err)
	}

	err = tmpFile.Close()
	if err != nil {
		t.Fatalf("unexpected failure! %s", err)
	}

	conf := config.New(false, "", tmpFile.Name())

	err = conf.Load()

	var pe *json.UnmarshalTypeError
	if err == nil || !errors.As(err, &pe) {
		t.Fatalf("should have returned json.SyntaxError: %s", err)
	}
}

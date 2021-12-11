package main

// Borrowed from Concourse: https://github.com/concourse/atc/blob/master/atccmd/file_flag.go

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FilesFlag is a flag for passing a path to one or more files on disk.
// The files are expected to be files, not directories, that actually exist.
type FilesFlag struct {
	fs *[]string
}

// Set implements pflag's Value interface.
func (f FilesFlag) Set(value string) error {
	for _, v := range strings.Split(value, ",") {
		stat, err := os.Stat(v)
		if err != nil {
			return err
		} else if stat.IsDir() {
			return fmt.Errorf("path %q is a directory, not a file", v)
		}

		abs, err := filepath.Abs(value)
		if err != nil {
			return err
		}

		*f.fs = append(*f.fs, abs)
	}

	return nil
}

// String implements pflag's Value interface.
func (f FilesFlag) String() string {
	return strings.Join(*f.fs, ",")
}

// Type implements pflag's Value interface.
func (f FilesFlag) Type() string {
	return "files"
}

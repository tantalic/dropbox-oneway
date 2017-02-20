package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

const (
	DROPBOX_TOKEN     = "DROPBOX_TOKEN"
	LOCAL_DIRECTORY   = "LOCAL_DIRECTORY"
	DROPBOX_DIRECTORY = "DROPBOX_DIRECTORY"
)

// options are the runtime configurations which affect
// application behavior.
type options struct {
	DropboxToken     string
	DropboxDirectory string
	LocalDirectory   string
}

// optsFromEnv loads the configuration options for the
// application from the environment or an error if the
// provided options are invalid.
func optsFromEnv() (*options, error) {
	opts := &options{
		DropboxToken:     os.Getenv(DROPBOX_TOKEN),
		DropboxDirectory: os.Getenv(DROPBOX_DIRECTORY),
	}

	if opts.DropboxToken == "" {
		return nil, fmt.Errorf("dropbox token must be provided in the %s environment variable", DROPBOX_TOKEN)
	}

	if opts.DropboxDirectory == "" {
		return nil, fmt.Errorf("dropbox directory must be provided in the %s environment variable", DROPBOX_DIRECTORY)
	}

	var err error
	opts.LocalDirectory, err = filepath.Abs(os.Getenv(LOCAL_DIRECTORY))
	if err != nil {
		return nil, errors.Wrapf(err, "invalid directory in %s environment variable", LOCAL_DIRECTORY)
	}

	err = os.MkdirAll(opts.LocalDirectory, os.ModePerm)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to create directory %s", opts.LocalDirectory)
	}

	return opts, nil
}

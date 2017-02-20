package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"tantalic.com/dropbox"

	"github.com/cenkalti/backoff"
	"github.com/pkg/errors"
)

// download uses the provided dropbox client to download
// the requested file to the given location. If a download
// fails multiple attemps will be made with an exponential
// backoff.
func download(client dropbox.Client, file dropbox.MetaData, location string) {
	path := filepath.Join(location, file.DisplayPath)
	dir := filepath.Dir(path)

	// Create the directory where this will live (if it does not
	// already exist)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Printf("unable to create directory %s: %s", dir, err.Error())
		return
	}

	attempt := func() error {
		reader, err := client.Download(file.Path)
		if err != nil {
			return errors.Wrapf(err, "error downloading %s", file.Path)
		}
		defer reader.Close()

		// Create temporary file for writing
		tempfile, err := ioutil.TempFile(dir, ".dropbox-oneway-dl-"+file.Name)
		if err != nil {
			return errors.Wrapf(err, "unable to open temporary file for %s\n", file.Name)
		}
		defer tempfile.Close()
		defer os.Remove(tempfile.Name())

		// Write to temporary file
		written, err := io.Copy(tempfile, reader)
		if err != nil {
			return errors.Wrapf(err, "error downloading %s", file.Path)
		}

		// Move to final destination
		os.Rename(tempfile.Name(), path)

		log.Printf("downloaded %s (%d bytes)", path, written)
		return nil
	}

	notify := func(err error, sleep time.Duration) {
		log.Printf("Error downloading %s: %s. Retry in %s.\n", file.Path, err.Error(), sleep)
	}

	b := &backoff.ExponentialBackOff{
		InitialInterval:     500 * time.Millisecond,
		RandomizationFactor: .75,
		Multiplier:          2,
		MaxInterval:         5 * time.Minute,
		MaxElapsedTime:      30 * time.Minute,
		Clock:               backoff.SystemClock,
	}
	err = backoff.RetryNotify(attempt, b, notify)
	if err != nil {
		log.Printf("Error downloading %s: %s. Will NOT retry.\n", file.Path, err.Error())
	}
}

package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"tantalic.com/dropbox"
)

const (
	Version = "0.1.1"
)

func main() {
	opts, err := optsFromEnv()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	log.Printf("Starting dropbox-oneway (v%s)\n", Version)

	client := dropbox.Client{
		AuthorizationToken: opts.DropboxToken,
	}

	// Creates a channel for new/modified files and background
	// goroutine that downloads each file sent to the channel.
	fileChan := make(chan dropbox.MetaData)
	go func() {
		for file := range fileChan {
			if file.IsFile() {
				go download(client, file, opts.LocalDirectory)
			}
		}
	}()

	// Watches for new/modified files in the given directory
	// adding them to the channel created above.
	err = client.Watch(time.Second*10, dropbox.ListOptions{
		Path:      opts.DropboxDirectory,
		Recursive: true,
	}, fileChan)

	if err != nil {
		log.Printf("Error: %s\n", err.Error())
	}

	// Wait for exit signal
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		log.Println(sig)
		done <- true
	}()
	<-done
	log.Println("exiting")
}

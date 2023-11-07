package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
)

const TieredLocalityFilePath = "/etc/cm/aa"

func main() {
	var log logr.Logger
	zapLog, err := zap.NewDevelopment()
	if err != nil {
		panic(fmt.Sprintf("new zap log error (%v)?", err))
	}
	log = zapr.NewLogger(zapLog)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Error(err, "new file watcher failed")
	}
	err = watcher.Add(TieredLocalityFilePath)
	if err != nil {
		log.Error(err, "new file watcher failed")
	}

	//
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				// linux vim is a RENAME operation
				// configmap is a soft link, notified only when its target file have been changed.
				// target file (old) removed
				if event.Has(fsnotify.Remove) {
					log.Info("tiered locality config modified, reload it.")
					err = watcher.Add(TieredLocalityFilePath)
					if err != nil {
						log.Error(err, "new file watcher failed")
					}
				}
			case err := <-watcher.Errors:
				log.Error(err, "watch tiered locality config error")
			}
		}
	}()

	for {

	}
}

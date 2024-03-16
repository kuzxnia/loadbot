package agent

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
)

func (a *Agent) WatchConfigFile(configFile string) (err error) {
	// Start listening for events.
	if a.configChange == nil {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}
		a.configChange = watcher
	}

	go func() {
		for {
			if a.state != AgentStateLeader {
				a.stateChange.L.Lock()
				log.Println("waiting for node to be elected as master")
				a.stateChange.Wait() // wait until state change
        log.Println("Started watching config file - new state: Leader")
			}

			select {
			case err, ok := <-a.configChange.Errors:
				if !ok { // channel was closed (i.e. Watcher.Close() was called)
					return
				}
				log.Println("error:", err)
			case event, ok := <-a.configChange.Events:
				if !ok { // channel was closed (i.e. Watcher.Close() was called)
					return
				}

				// Ignore files we're not interested in. Can use a
				if event.Name != configFile {
					continue
				}
				if event.Has(fsnotify.Write) {
					log.Println("modified file:", event.Name, " applying new configuration")
					a.ApplyConfigFromFile(event.Name)
				}
			}
		}
	}()

	st, err := os.Lstat(configFile)
	if err != nil {
		return
	}

	if st.IsDir() {
		return errors.New(configFile + " is a directory, not a file")
	}

	// Watch the directory, not the file itself. This solves various issues where files are frequently
	// renamed, such as editors saving them.
	err = a.configChange.Add(filepath.Dir(configFile))
	if err != nil {
		log.Fatal(err)
		return
	}
	return
}

// agent config should be removed from base config and later
// zarzadzanie workloadem przez

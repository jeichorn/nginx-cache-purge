package nginxcp
//https://github.com/nathany/looper
// https://raw.githubusercontent.com/nathany/looper/5ad8fceca57dd2699d9e4c131c8a05ba24303f1d/watch.go

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"
    "fmt"

	"github.com/go-fsnotify/fsnotify"
)

type RecursiveWatcher struct {
	*fsnotify.Watcher
	Files   chan string
	Folders chan string
}

func NewRecursiveWatcher(path string) (*RecursiveWatcher, error) {
    DebugMessage(fmt.Sprintf("Watching %s", path))
	folders := Subfolders(path)
	if len(folders) == 0 {
		return nil, errors.New("No folders to watch.")
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	rw := &RecursiveWatcher{Watcher: watcher}

	rw.Files = make(chan string, 10)
	rw.Folders = make(chan string, len(folders))

	for _, folder := range folders {
		rw.AddFolder(folder)
	}
	return rw, nil
}

func (watcher *RecursiveWatcher) AddFolder(folder string) {
	err := watcher.Add(folder)
	if err != nil {
		log.Println("Error watching: ", folder, err)
	}
	watcher.Folders <- folder
}

func (watcher *RecursiveWatcher) Run(debug bool) {
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				// create a file/directory
				if event.Op&fsnotify.Create == fsnotify.Create {
					fi, err := os.Stat(event.Name)
					if err != nil {
						// eg. stat .subl513.tmp : no such file or directory
						if debug {
                            DebugMessage(event.Name)
							DebugError(err)
						}
					} else if fi.IsDir() {
						if debug {
							DebugMessage("Detected new directory %s", event.Name)
						}
						if !shouldIgnoreFile(filepath.Base(event.Name)) {
							watcher.AddFolder(event.Name)
						}
					} else {
						if debug {
							DebugMessage("Detected new file %s", event.Name)
						}
						watcher.Files <- event.Name // created a file
					}
				}

				if event.Op&fsnotify.Write == fsnotify.Write {
					// modified a file, assuming that you don't modify folders
					if debug {
						DebugMessage("Detected file modification %s", event.Name)
					}
					watcher.Files <- event.Name
				}

				if event.Op&fsnotify.Remove == fsnotify.Remove {
					// deleted a file
					if debug {
						DebugMessage("Detected file modification %s", event.Name)
					}
					watcher.Files <- event.Name
				}

			case err := <-watcher.Errors:
                if (err == nil) {
                    continue;
                }

				log.Println("watcher error", err)
			}
		}
	}()
}

// Subfolders returns a slice of subfolders (recursive), including the folder provided.
func Subfolders(path string) (paths []string) {
	filepath.Walk(path, func(newPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			name := info.Name()
			// skip folders that begin with a dot
			if shouldIgnoreFile(name) && name != "." && name != ".." {
				return filepath.SkipDir
			}
			paths = append(paths, newPath)
		}
		return nil
	})
	return paths
}

// shouldIgnoreFile determines if a file should be ignored.
// File names that begin with "." or "_" are ignored by the go tool.
func shouldIgnoreFile(name string) bool {
	return strings.HasPrefix(name, ".") || strings.HasPrefix(name, "_")
}

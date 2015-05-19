package nginxcp;

import (
    "log"
)

func EventLoop(path string, keys *CacheKeys, debug bool) {
    watcher, err := NewRecursiveWatcher(path)
    if err != nil {
        log.Fatal(err)
    }

    queue := NewRedisQueue()
    watcher.Run(debug)
    loadInitial(path, keys)
    go queue.Run()
    defer watcher.Close()

    for {
        select {
        case file := <-watcher.Files:
            keys.addEntryFromFile(file)
        case job := <-queue.Jobs:
            keys.removeUsingJob(job)
        }

    }
}



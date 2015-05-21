package nginxcp;

import (
    "log"
    "fmt"
    "os"
)

func EventLoop(path string, keys *CacheKeys, debug int) {
    var pingFile string = fmt.Sprintf("%s/ping", path)
    os.Remove(pingFile)
    os.Remove(fmt.Sprintf("%s/.ping", path))

    watcher, err := NewRecursiveWatcher(path)
    if err != nil {
        log.Fatal(err)
    }

    queue := NewRedisQueue()
    watcher.Run()
    loadInitial(path, keys)
    go queue.Run()
    defer watcher.Close()
    ping := Ping{}

    go ping.Run(pingFile, keys)
    go Info(keys)

    for {
        select {
        case file := <-watcher.Files:
            keys.addEntryFromFile(file)
        case job := <-queue.Jobs:
            keys.removeUsingJob(job)
            queue.completeJob(job)
        }

    }
}



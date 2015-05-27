package nginxcp;

import (
    "fmt"
    "os"
)

func EventLoop(path string, debug int) {
    var pingFile string = fmt.Sprintf("%s/ping", path)
    os.Remove(pingFile)
    os.Remove(fmt.Sprintf("%s/.ping", path))

    queue := NewRedisQueue()
    queue.clearInPurgeList()
    go queue.Run()
    
    purge := NewPurge(path)
    go purge.Run()

    for {
        select {
        case job := <-queue.Jobs:
            purge.Jobs <- job
        }
    }
}

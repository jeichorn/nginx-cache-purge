package main

import (
    "flag"

    "nginxcp"
)
func main() {
    var debug int
    var path string


    flag.IntVar(&debug, "debug", 0, "Enable Debug")
    flag.StringVar(&path, "path", ".", "Path to watch")
    flag.Parse()

    nginxcp.Header()
    if (debug > 0) {
        nginxcp.DebugLevel = debug
        nginxcp.DebugEnabled()
    }

    cachekeys := nginxcp.NewCacheKeys()
    nginxcp.EventLoop(path, cachekeys, debug)
}

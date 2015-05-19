package main

import (
    "flag"

    "nginxcp"
)
func main() {
    var debug bool
    var path string


    flag.BoolVar(&debug, "debug", true, "Enable Debug")
    flag.StringVar(&path, "path", ".", "Path to watch")
    flag.Parse()

    nginxcp.Header()
    if (debug) {
        nginxcp.DebugEnabled()
    }

    cachekeys := nginxcp.NewCacheKeys()
    nginxcp.EventLoop(path, cachekeys, debug)
}

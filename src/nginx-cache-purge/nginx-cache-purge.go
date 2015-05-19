package main

import (
    "flag"
    "path/filepath"

    "nginxcp"
)
func main() {
    var debug int
    var path string


    flag.IntVar(&debug, "debug", 0, "Enable Debug")
    flag.StringVar(&path, "path", ".", "Path to watch")
    flag.Parse()

    path = filepath.Clean(path)
    path, _ = filepath.EvalSymlinks(path)
    path, _ = filepath.Abs(path)

    nginxcp.Header()
    if (debug > 0) {
        nginxcp.DebugLevel = debug
        nginxcp.DebugEnabled()
    }

    cachekeys := nginxcp.NewCacheKeys()
    nginxcp.EventLoop(path, cachekeys, debug)
}

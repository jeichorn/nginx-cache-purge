package main

import (
    "flag"
    "path/filepath"

    "nginxcp"
)
func main() {
    var debug int
    var path string
    var ansi bool


    flag.IntVar(&debug, "debug", 0, "Enable Debug")
    flag.StringVar(&path, "path", ".", "Path to watch")
    flag.BoolVar(&ansi, "ansi", true, "Colored Ansible Output")
    flag.Parse()

    path = filepath.Clean(path)
    path, _ = filepath.EvalSymlinks(path)
    path, _ = filepath.Abs(path)

    nginxcp.AnsiOutput = ansi
    nginxcp.Header()
    if (debug > 0) {
        nginxcp.DebugLevel = debug
        nginxcp.DebugEnabled()
    }

    nginxcp.EventLoop(path, debug)
}

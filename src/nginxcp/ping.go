package nginxcp

import (
    "os"
    "fmt"
    "time"
    "errors"
    "path/filepath"
)

type Ping struct {
    last int64
}

func (ping *Ping) Run(filename string, keys *CacheKeys) {
    for {
        if (ping.last != 0) {
            info := keys.getFile(filename)
            var test string = fmt.Sprintf("ping--ping/%d", ping.last)
            if (info.key != test) {
                PrintError(errors.New(fmt.Sprintf("Inotify test failed %s != %s", info.key, test)))
                os.Exit(3)
            } else {
                keys.removeEntry(filename, true)
                PrintTrace2("Ping Successful")
            }
        }

        tmpfile := fmt.Sprintf("%s/.ping", filepath.Dir(filename))
        file, err := os.Create(tmpfile)

        if (err != nil) {
            PrintError(err)
            os.Exit(2)
        }

        ping.last = time.Now().UnixNano()
        file.WriteString(fmt.Sprintf("KEY: ping--ping/%d\n", ping.last))
        file.Close()

        err_r := os.Rename(tmpfile, filename)
        if (err_r != nil) {
            PrintError(err_r)
            os.Exit(3)
        }

        time.Sleep(10 * time.Second)
    }
}

func Info(keys *CacheKeys) {
    for {
        time.Sleep(60 * 5 * time.Second)
        keys.printKeyCounts()
    }
}

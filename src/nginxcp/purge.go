package nginxcp

import (
    "path/filepath"
    "fmt"
    "regexp"
    "log"
    "strings"
    "os"
    "errors"
    "time"

    "github.com/jeichorn/go-caching"
    "github.com/jeichorn/go-caching/container/concurrent"
    "github.com/jeichorn/go-caching/container/memory"
    _ "github.com/jeichorn/go-caching/container/memory/arc"
)

var jobSplitter = regexp.MustCompile(`^([^:]+)::(.+)$`)

type Purge struct {
    Path string
    Jobs chan string
    Cache *caching.Cache
}

func NewPurge(path string) *Purge {
    arc := memory.ARC.New(100000)
    concurrent := concurrent.New(arc)
    cache := caching.NewCache(concurrent)
    return &Purge{path, make(chan string, 10), cache}
}

func (purge *Purge) Purge(job string) {
    PrintInfo("Purge job %s", job)
    matched := jobSplitter.FindAllStringSubmatch(job, -1)
    var host string
    var regex string
    if (len(matched) == 1 && len(matched[0]) == 3) {
        host = string(matched[0][1])
        regex = string(matched[0][2])
    } else {
        PrintError(errors.New(fmt.Sprintf("Bad Job: %s", job)))
        return
    }

    regex = strings.Replace(regexp.QuoteMeta(regex), "\\(\\.\\*\\)", "(.*)", -1)
    regexString := fmt.Sprintf(`^([^-]+--)?(https?)?%s%s(\?.*)?$`, host, regex)

    tester, err := regexp.Compile(regexString)

    if (err != nil) {
        log.Println("Bad regex", err)
    }


    var count int = 0
    var deleted int = 0;
    filepath.Walk(purge.Path, func(path string, info os.FileInfo, err error) error {
        if (err != nil || info.IsDir()) {
            return nil
        }

        item, err := purge.Cache.Get(path)
        if (err != nil) {
            PrintError(err)
        }
        var key string = "BADKEY"
        if (item != nil) {
            PrintInfo("Got %#v from cache", item)
            if str, ok := item.(string); ok {
                key = str
            }
        } else {
            PrintTrace2("Cache Miss: %s", path)
            info := keyFromFile(path)
            key = info.key
            purge.Cache.Set(path, key, caching.NewExpiration(time.Now().Add(time.Minute * 30), time.Minute * 30))
        }

        if (tester.MatchString(key)) {
            PrintTrace1("Found a match: %s", key)
            os.Remove(path)
            deleted++
        } else {
            PrintTrace3("Miss: %s", key)
        }
        count++

        return nil
    })

    PrintInfo("Tested %d files deleted %d, %s", count, deleted, job)
}

func (purge *Purge) Run() {
    for {
        select {
        case job := <-purge.Jobs:
            purge.Purge(job)
        }
    }
}

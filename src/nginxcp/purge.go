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
    Jobs chan JobBag
    Cache *caching.Cache
}

func NewPurge(path string) *Purge {
    arc := memory.ARC.New(300000)
    concurrent := concurrent.New(arc)
    cache := caching.NewCache(concurrent)
    return &Purge{path, make(chan JobBag, 10), cache}
}

func (purge *Purge) Purge(jobs JobBag) {
    var list map[string]map[string]map[string]int = make(map[string]map[string]map[string]int)
    filepath.Walk(purge.Path, func(path string, info os.FileInfo, err error) error {
        if (err != nil || info.IsDir()) {
            return nil
        }

        item, err := purge.Cache.Get(path)
        if (err != nil) {
            PrintError(err)
        }
        var key string = "BADKEY"
        var domain string = "BADDOMAIN"
        var altdomain string = "BADDOMAIN"
        if (item != nil) {
            PrintTrace3("Got %#v from cache", item)
            if info, ok := item.(*CacheFileInfo); ok {
                key = info.key
                domain = info.domain
                altdomain = info.altdomain
            }
        } else {
            PrintTrace2("Cache Miss: %s", path)
            info := keyFromFile(path)
            key = info.key
            domain = info.domain
            altdomain = info.altdomain
            purge.Cache.Set(path, info, caching.NewExpiration(time.Now().Add(time.Minute * 30), time.Minute * 30))
        }
        if _, ok := list[domain]; !ok {
            list[domain] = make(map[string]map[string]int)
        }
        if _, ok := list[domain][key]; !ok {
            list[domain][key] = make(map[string]int)
        }
        list[domain][key][path] = 1

        if domain != altdomain {
            PrintTrace3("Have an altdomain %s, key: %s", altdomain, key)
            if _, ok := list[altdomain]; !ok {
                list[altdomain] = make(map[string]map[string]int)
            }
            if _, ok := list[altdomain][key]; !ok {
                list[altdomain][key] = make(map[string]int)
            }
            list[altdomain][key][path] = 1
        }

        return nil
    })

    for _, job := range jobs.Bag {
        matched := jobSplitter.FindAllStringSubmatch(job, -1)
        var host string
        var regex string
        if (len(matched) == 1 && len(matched[0]) == 3) {
            host = string(matched[0][1])
            regex = string(matched[0][2])
        } else {
            PrintError(errors.New(fmt.Sprintf("Bad Job: %s", job)))
            continue
        }

        PrintInfo("Purge job %s %s", host, regex)

        // we quote out most regex, but we all
        // (.*)
        // (/?)
        regex = strings.Replace(strings.Replace(regexp.QuoteMeta(regex), "\\(\\.\\*\\)", "(.*)", -1), "\\(/\\?\\)", "(/?)", -1)

        regexString := fmt.Sprintf(`^([^-]+--)?(https?)?%s%s(\?.*)?$`, host, regex)

        tester, err := regexp.Compile(regexString)

        if (err != nil) {
            log.Println("Bad regex", err)
        }


        var count int = 0
        var deleted int = 0;


        // check if the host exists at all, if it doesn't we can bail right away
        if _, ok := list[host]; !ok {
            PrintTrace1("No keys for: %s", host)
        } else {
            for domain, keys := range list {
                if (domain != host) {
                    continue;
                }

                for key, files := range keys {
                    if (tester.MatchString(key)) {
                        PrintTrace1("Found a match: %s", key)
                        for file, _ := range files {
                            os.Remove(file)
                            deleted++
                        }
                    } else {
                        PrintTrace3("Miss: %s", key)
                    }
                    count+=len(files)
                }
            }
        }

        PrintInfo("Tested %d files deleted %d, %s", count, deleted, job)
    }
}

func (purge *Purge) Run() {
    for {
        select {
        case job := <-purge.Jobs:
            purge.Purge(job)
        }
    }
}

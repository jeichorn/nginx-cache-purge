package nginxcp

import (
    "path/filepath"
    "fmt"
    "regexp"
    "log"
    "strings"
    "os"
    "errors"
)

var jobSplitter = regexp.MustCompile(`^([^:]+)::(.+)$`)

type Purge struct {
    Path string
    Jobs chan string
}

func NewPurge(path string) *Purge {
    return &Purge{path, make(chan string, 10)}
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

        key := keyFromFile(path)

        if (tester.MatchString(key.key)) {
            PrintTrace1("Found a match: %s", key.key)
            os.Remove(path)
            deleted++
        } else {
            PrintTrace2("Miss: %s", key.key)
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

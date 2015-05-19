package nginxcp

import (
    "path/filepath"
    "os"
    "io"
    "bufio"
    "fmt"
    "log"
    "errors"
    "regexp"
)

type CacheFileInfo struct {
    key string
    domain string
    deleted bool
    successful bool
}

// standard--woocommerce.bluga.info/grid-3-column/?
// standard--httpswww.bit9.com/forms/free-antivirus-plus-endpoint-protection/?campaign=70180000000f5DK
var domainFromKey = regexp.MustCompile(`^[^-]+-[^-]*-(?:https?)?([^/?]+)`)

func loadInitial(cachePath string, keys *CacheKeys) {

    filepath.Walk(cachePath, func(path string, info os.FileInfo, err error) error {
        if (err != nil || info.IsDir()) {
            return nil
        }

        keys.addEntryFromFile(fmt.Sprintf("%s/%s", cachePath, info.Name()))

        return nil
    })
}

func keyFromFile(file string) *CacheFileInfo {
    info := new(CacheFileInfo)
    f, err := os.Open(file)
    if err != nil {
        info.deleted = true
        return info
    }
    defer f.Close()
    r := bufio.NewReaderSize(f, 4*1024)

    line, isPrefix, err := r.ReadLine()

    for err == nil && !isPrefix {
        s := string(line)
        line, isPrefix, err = r.ReadLine()

        if (len(s) > 5 && s[0:4] == "KEY:") {
            info.key = s[5:len(s)]
            info.successful = true

            matched := domainFromKey.FindAllStringSubmatch(info.key, -1)
            if (len(matched) == 1 && len(matched[0]) == 2) {
                info.domain = string(matched[0][1])
            }

            return info
        }
    }
    if isPrefix {
        log.Println(errors.New(fmt.Sprintf("buffer size to small: %s", file)))
        return info
    }
    if err != io.EOF {
        log.Println(err)
        return info
    }

    return info
}

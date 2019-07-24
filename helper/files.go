package helper

import (
    "encoding/json"
    "io/ioutil"
    "log"
    "os"
    "path/filepath"
)

func ReadJsonFile(fPath string, res interface{}) {
    jsonFile, err := os.Open(fPath)
    if err != nil {
        log.Println(err)
        return
    }
    defer jsonFile.Close()

    byteValue, _ := ioutil.ReadAll(jsonFile)
    json.Unmarshal([]byte(byteValue), res)
    return
}

func WriteFile(data []byte, fpath, fname string) (e error) {
    if IsFileNotExist(fpath) {
        e = os.MkdirAll(fpath, 0777)
        if e != nil {
            return
        }
    }
    e = ioutil.WriteFile(filepath.Join(fpath, fname), data, 0777)
    return
}

func IsFileNotExist(path string) bool {
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return true
    }

    return false
}

func IsFileExist(path string) bool {
    if _, err := os.Stat(path); err == nil {
        return true
    }

    return false
}

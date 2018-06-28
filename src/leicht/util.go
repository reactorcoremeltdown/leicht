package main

import (
    "os"
    "log"
    "io/ioutil"
    "encoding/json"
)

type Config struct {
    Token string
    Socket string
    SocketMode os.FileMode
    Script string
    Logging bool
    Debug bool
    LogDirectory string
    DoNotLogBlacklisted bool
    WhitelistEnabled bool
    Whitelist []string
}

func LoadConfig(cfgpath string) (c *Config, err error) {
    var bfile []byte
    if bfile, err = ioutil.ReadFile(cfgpath); err != nil {
        log.Fatalf("Error reading config file: %s\n", err.Error())
    }
    c = new(Config)
    err = json.Unmarshal(bfile, &c)

    if c.Socket == "" {
      c.Socket = "/var/run/leicht/leicht.socket"
    }

    return
}

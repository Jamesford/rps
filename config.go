package main

import (
    "os"
    "log"
    "github.com/BurntSushi/toml"
)

// Config struct
type Config struct {
    Port string
    VerifyToken string
    Scope string
    ClientID string
    ClientSecret string
    RedirectURI string
    SuccessRedirect string
    FailureRedirect string
}

// ReadConfig reads info from config file
func ReadConfig() Config {
    configfile := "config.toml"
    _, err := os.Stat(configfile)
    if err != nil {
        log.Fatal("Config file is missing: ", configfile)
    }

    var config Config
    if _, err := toml.DecodeFile(configfile, &config); err != nil {
        log.Fatal(err)
    }

    return config
}
package main

import (
    "fmt"
    "github.com/cemalkilic/shorten-backend/config"
    "net/http"
    "os"
)

func main() {
    cfg, _ := config.LoadConfig(".")

    res, err := http.Get(fmt.Sprintf("http://127.0.0.1:%s/healthcheckz", cfg.ServerPort))

    if res != nil && res.StatusCode != 200 {
        os.Exit(1)
    }

    if err != nil {
        os.Exit(2)
    }
}

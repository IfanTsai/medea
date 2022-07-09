package main

import (
	"log"
	"{{ .GoModulePath }}/web"
)

func main() {
    ginServer, err := web.NewGinServer("0.0.0.0:8080")
    if err != nil {
        log.Fatal("cannot new gin server:", err)
    }

    web.Run(ginServer)
}

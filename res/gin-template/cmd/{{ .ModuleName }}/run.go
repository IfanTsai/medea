package main

import (
	"fmt"
	"log"
	"{{ .GoModulePath }}/web"

	"github.com/IfanTsai/go-lib/config"
)

func main() {
    config.Init()

    ip := config.GetIPWithDefault("0.0.0.0")
    port := config.GetPortWithDefault(8080)
    addr := fmt.Sprintf("%s:%d", ip, port)

    ginServer, err := web.NewGinServer(addr)
    if err != nil {
        log.Fatal("cannot new gin server:", err)
    }

    web.Run(ginServer)
}

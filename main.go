package main

import (
	"fmt"
	"github.com/ASparkOfFire/ignis/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/js", utils.WASIWrapper("/home/asparkoffire/jstest/index.wasm"))
	r.Any("/api/v1/*any", utils.WASIWrapper("/home/asparkoffire/GolandProjects/testproject/example/example.wasm"))

	fmt.Println("Listening on 6969")
	r.Run(":6969")
}

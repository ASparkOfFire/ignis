package main

import (
	"fmt"
	"github.com/ASparkOfFire/ignis/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/", utils.WASIWrapper("/home/asparkoffire/GolandProjects/testproject/example/example.wasm"))
	r.GET("/user/:id", utils.WASIWrapper("/home/asparkoffire/GolandProjects/testproject/example/example.wasm"))
	r.GET("/user", utils.WASIWrapper("/home/asparkoffire/GolandProjects/testproject/example/example.wasm"))

	fmt.Println("Listening on 6969")
	r.Run(":6969")
}

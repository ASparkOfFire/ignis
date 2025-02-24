package main

import (
	"fmt"
	"github.com/ASparkOfFire/ignis/internal/cache"
	"github.com/ASparkOfFire/ignis/internal/runtime"
	"github.com/ASparkOfFire/ignis/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func main() {
	r := gin.Default()
	modCache := cache.NewModCache[uuid.UUID]()
	id := uuid.MustParse("006e267f-b578-43ba-a844-7c34aa2bf00a")
	idJs := uuid.MustParse("006e267f-b578-43ba-a844-7c34aa2bf00b")

	r.Any("/api/v1/*any", utils.WASIWrapper("./example/go/example.wasm", modCache, runtime.RuntimeEngineWASM, id))
	r.Any("/js", utils.WASIWrapper("./example/js/dist/example.js", modCache, runtime.RuntimeEngineJS, idJs))

	fmt.Println("Listening on 6969")
	r.Run(":6969")
}

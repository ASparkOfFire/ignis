package main

import (
	"github.com/ASparkOfFire/ignis/sdk"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func HandleRoot(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusOK, gin.H{
		"message": "Hello from Ignis",
	})
}

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}
type CreateUserRequest struct {
	Username string `json:"username"`
}

var (
	Users = map[int]User{
		1: {ID: 1, Username: "FastTiger123", CreatedAt: time.Date(2023, 8, 14, 15, 30, 0, 0, time.UTC)},
		2: {ID: 2, Username: "CleverEagle99", CreatedAt: time.Date(2023, 11, 7, 10, 15, 0, 0, time.UTC)},
		3: {ID: 3, Username: "BraveWolf42", CreatedAt: time.Date(2024, 2, 19, 20, 45, 0, 0, time.UTC)},
		4: {ID: 4, Username: "CoolDragon77", CreatedAt: time.Date(2024, 5, 3, 8, 10, 0, 0, time.UTC)},
		5: {ID: 5, Username: "SharpHawk88", CreatedAt: time.Date(2024, 7, 1, 14, 20, 0, 0, time.UTC)},
	}
)

func HandleGetUser(c *gin.Context) {
	idParam := c.Params.ByName("id")
	if idParam == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "id is required",
		})
		return
	}
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "id is invalid",
		})
		return
	}
	user, ok := Users[id]
	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "user not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
	return
}

func HandleGetUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"user": Users,
	})
	return
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.GET("/", HandleRoot)
	router.GET("/user/:id", HandleGetUser)
	router.GET("/user", HandleGetUsers)

	sdk.Handle(router)
}

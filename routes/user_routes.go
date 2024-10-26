package routes

import (
	controllers "github.com/ayush/mongo-kibana/controller"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterUserRoutes(r *gin.Engine, client *mongo.Client) {
	controllers.InitUserCollection(client)
	r.GET("/users", controllers.GetUsers)
	r.GET("/users/:id", controllers.GetUserByID)
	r.POST("/users", controllers.CreateUser)
	r.PUT("/users/:id", controllers.UpdateUser)
	r.PATCH("/users/:id", controllers.PatchUser)
	r.DELETE("/users/:id", controllers.DeleteUser)
}

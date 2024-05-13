package main

import (
	"github.com/gin-gonic/gin"
	L "nagarjuna2323/books_api/internal/middlewares/logger"
	userroutes "nagarjuna2323/books_api/internal/routes/userRoutes"
)

func main() {
	port := ":3000"
	userRouter := gin.Default()
	userroutes.UserRoutes(userRouter)
	err := userRouter.Run(port)
	if err != nil {
		L.BKSLog("E", "Unable to Serve Api on "+L.PrintStruct(port), nil)
		return
	}

}

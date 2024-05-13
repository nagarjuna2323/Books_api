package userroutes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"nagarjuna2323/books_api/internal/controllers"
	"nagarjuna2323/books_api/internal/middlewares/authorization"
	"nagarjuna2323/books_api/internal/routes"
)

func UserRoutes(userRoute *gin.Engine) {
	//enable cors
	userRoute.Use(cors.Default())
	userRoute.POST(routes.UserLogin, func(ctx *gin.Context) {
		controllers.LogInService(ctx)
	})
	userRoute.POST(routes.UserSignup, func(context *gin.Context) {
		controllers.SignUpService(context)
	})
	userRoute.GET(routes.GetBooks, authorization.AuthorizeRequest(), func(ctx *gin.Context) {
		controllers.Home(ctx)
	})
	userRoute.PUT(routes.AddingNewBooks, authorization.AuthorizeRequest(), func(ctx *gin.Context) {
		controllers.AddBook(ctx)
	})
	userRoute.DELETE(routes.DeletingBooks, authorization.AuthorizeRequest(), func(ctx *gin.Context) {
		controllers.DeleteBook(ctx)
	})
}

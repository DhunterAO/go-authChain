package server

import (
	"github.com/kataras/iris"
)

func addUrls(app *iris.Application) {
	// post info
	app.Post("/operation", postOperation)
}

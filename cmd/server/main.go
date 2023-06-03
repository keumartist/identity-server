package main

import "art-sso/internal/bootstrap"

func main() {
	app := bootstrap.InitApp()

	app.Listen(":3000")
}

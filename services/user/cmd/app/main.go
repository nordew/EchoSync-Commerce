package main

import "userService/internal/app"

func main() {
	if err := app.Run(); err != nil {
		panic(err)
	}
}

package main

import "marketService/internal/app"

func main() {
	if err := app.Run(); err != nil {
		panic(err)
	}
}

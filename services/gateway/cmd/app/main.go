package main

import "gateway/internal/app"

func main() {
	if err := app.Run(); err != nil {
		panic(err)
	}
}

package main

import (
	"backend/src/app"
)

func main() {
	a := &app.App{}
	a.Init()
	a.Run()
}

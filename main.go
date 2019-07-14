package main

import (
	"controller"
)

func main()  {
	a := controller.App{}
	a.Initialize()
	a.Run(":9000")
}
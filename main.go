package main

import (
	"github.com/lingdor/logwaiter/bootstrap"
	"github.com/lingdor/logwaiter/common"
)

func main() {

	defer func() {
		defer common.AppRecover()
	}()

	var app = bootstrap.NewApp()
	app.Start()

}

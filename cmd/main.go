package main

import (
	"github.com/flexstack.ai/membots-be/conf"
	_ "github.com/flexstack.ai/membots-be/docs"
	app "github.com/flexstack.ai/membots-be/internal"
)

func main() {
	config := conf.GetConfiguration()
	app.RunApp(config)
}

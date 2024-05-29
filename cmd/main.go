package main

import (
	"github.com/flexstack.ai/membots-be/conf"
	app "github.com/flexstack.ai/membots-be/internal"
)

func main() {
	config := conf.GetConfiguration()
	app.RunApp(config)
}

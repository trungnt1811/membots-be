package main

import (
	"github.com/astraprotocol/membots-be/conf"
	_ "github.com/astraprotocol/membots-be/docs"
	app "github.com/astraprotocol/membots-be/internal"
)

func main() {
	config := conf.GetConfiguration()
	app.RunApp(config)
}

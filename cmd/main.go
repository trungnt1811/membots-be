package main

import (
	"github.com/astraprotocol/affiliate-system/conf"
	_ "github.com/astraprotocol/affiliate-system/docs"
	app "github.com/astraprotocol/affiliate-system/internal"
)

func main() {
	config := conf.GetConfiguration()
	app.RunApp(config)
}

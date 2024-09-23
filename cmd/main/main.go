package main

import (
	"api_chat/api"
)

func main() {
	init := api.NewInitializerAPI("./config/config")
	init.InitConfig()
	init.Startup()
}

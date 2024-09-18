package main

import (
	"api_chat/api"
)

func main() {
	//init := api.NewInitializerAPI("./config/config")
	init := api.NewInitializerAPI("./config/config.dev")
	init.InitConfig()
	init.Startup()
}

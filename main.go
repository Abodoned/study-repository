package main

import (
	httpcore "MinerGame/http_core"
	repositor "MinerGame/repositor"
	"fmt"
)

func main() {

	factory := repositor.NewFactory()
	handlers := httpcore.NewHTTPHandlers(factory)
	server := httpcore.NewHTTPServer(handlers)
	if err := server.StartServer(); err != nil {
		fmt.Println("Ошибка запуска сервера!")
	}

}

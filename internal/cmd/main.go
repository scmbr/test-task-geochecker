package main

import "github.com/scmbr/test-task-geochecker/internal/app"

const configsDir = "configs"

func main() {
	app.Run(configsDir)
}

package main

import (
	"os"

	std_widgets "github.com/therecipe/qt/widgets"
)

func main() {
	app := std_widgets.NewQApplication(len(os.Args), os.Args)

	app.Exec()
}

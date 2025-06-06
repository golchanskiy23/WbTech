package main

import "fmt"

func main() {
	app := NewApplication()
	editor := &Editor{Text: "Hello, World!", Selection: "World"}
	app.ActiveEditor = editor

	fmt.Println("Исходный текст:", editor.Text)

	app.ExecuteCommand(&CutCommand{BaseCommand{App: app, Editor: editor}})
	fmt.Println("После вырезания:", editor.Text)

	app.Undo()
	fmt.Println("После Undo:", editor.Text)
}

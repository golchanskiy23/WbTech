package main

type Editor struct {
	Text      string
	Selection string
	CursorPos int
}

func (e *Editor) GetSelection() string {
	return e.Selection
}

func (e *Editor) DeleteSelection() {
	e.Text = ""
}

func (e *Editor) ReplaceSelection(text string) {
	e.Text = text
}

type Command interface {
	Execute() bool
	Undo()
}

type BaseCommand struct {
	App    *Application
	Editor *Editor
	backup string
}

func (c *BaseCommand) SaveBackup() {
	c.backup = c.Editor.Text
}

func (c *BaseCommand) Undo() {
	c.Editor.Text = c.backup
}

type CopyCommand struct {
	BaseCommand
}

func (c *CopyCommand) Execute() bool {
	c.App.Clipboard = c.Editor.GetSelection()
	return false
}

func (c *CopyCommand) Undo() {}

type CutCommand struct {
	BaseCommand
}

func (c *CutCommand) Execute() bool {
	c.SaveBackup()
	c.App.Clipboard = c.Editor.GetSelection()
	c.Editor.DeleteSelection()
	return true
}

type PasteCommand struct {
	BaseCommand
}

func (c *PasteCommand) Execute() bool {
	c.SaveBackup()
	c.Editor.ReplaceSelection(c.App.Clipboard)
	return true
}

type UndoCommand struct {
	BaseCommand
}

func (c *UndoCommand) Execute() bool {
	c.App.Undo()
	return false
}

func (c *UndoCommand) Undo() {}

type CommandHistory struct {
	history []Command
}

func (h *CommandHistory) Push(c Command) {
	h.history = append(h.history, c)
}

func (h *CommandHistory) Pop() Command {
	if len(h.history) == 0 {
		return nil
	}
	last := h.history[len(h.history)-1]
	h.history = h.history[:len(h.history)-1]
	return last
}

type Application struct {
	Clipboard    string
	Editors      []*Editor
	ActiveEditor *Editor
	History      *CommandHistory
}

func (app *Application) ExecuteCommand(command Command) {
	if command.Execute() {
		app.History.Push(command)
	}
}

func (app *Application) Undo() {
	cmd := app.History.Pop()
	if cmd != nil {
		cmd.Undo()
	}
}

func NewApplication() *Application {
	return &Application{
		History: &CommandHistory{},
	}
}

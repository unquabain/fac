package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/Unquabain/fac/display"
	"github.com/Unquabain/fac/task"
	"github.com/jroimartin/gocui"
	yaml "gopkg.in/yaml.v2"
)

func printUsage() {
	fmt.Printf("Usage: %s taskfile.yaml\n", os.Args[0])
	fmt.Println(``)
	fmt.Println(`Options:`)
	fmt.Println(`  taskfile.yaml A YAML file with the things to do`)
	fmt.Println(``)
	fmt.Println(`Example Taskfile:`)
	fmt.Println(`---
Clear Logs:
  command: zsh
  args:
    - rm
    - logs/development.txt
    - logs/test.txt
  expectedReturnCode: 7
Update Bundler:
  command: bin/bundle
  args:
    - update
  dependencies:
    - Clear Logs
  environment:
    RAILS_ENV: development
  expectedStdOutRegex: Bundle complete!
  `)
	fmt.Println(``)
	fmt.Println(`(expectedStdErrRegex is supported as well)`)
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(-1)
	}
	yamlFile := os.Args[1]
	buff, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		log.Printf(`No can do, Compadre. %v`, err)
		printUsage()
		os.Exit(-2)
	}
	list := make(task.TaskList)
	err = yaml.Unmarshal(buff, &list)
	if err != nil {
		log.Printf(`No love here. %v`, err)
		printUsage()
		os.Exit(-3)
	}
	manager := &display.TaskLayoutManager{TaskList: list}

	g, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		log.Printf(`I could show you, but I'd have to charge: %v`, err)
		os.Exit(-4)
	}
	defer g.Close()
	g.SetManager(manager)

	handler := func(s *task.Task) {
		g.Update(func(gg *gocui.Gui) error {
			return manager.Layout(gg)
		})
	}
	go func() {
		err := list.RunAll(handler)
		if err != nil {
			log.Printf(`Ouch!: %v`, err)
			printUsage()
			os.Exit(-5)
		}
	}()

	err = g.SetKeybinding(
		"",
		gocui.KeyCtrlC,
		gocui.ModNone,
		func(_ *gocui.Gui, _ *gocui.View) error { return gocui.ErrQuit },
	)
	if err != nil {
		log.Printf(`No keybindings for you: %v`, err)
		os.Exit(-6)
	}

	if err = g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Printf(`I die. %v`, err)
		os.Exit(-7)
	}
}

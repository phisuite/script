package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	defaultHost = "localhost"
	defaultPort = 50051
)

type instructor struct {
	cmd *flag.FlagSet
	editor *string
	inspector *string
}

func main() {
	reader := initInspectCmd()
	writer := initUpdateCmd()

	if len(os.Args) < 2 {
		reader.cmd.Usage()
		writer.cmd.Usage()
		os.Exit(1)
	}
	var err error
	switch os.Args[1] {
	case "inspect":
		err = reader.inspect()
	case "update":
		err = writer.update()
	default:
		reader.cmd.Usage()
		writer.cmd.Usage()
		os.Exit(1)
	}
	if err != nil {
		log.Fatal(err)
	}
}

func initUpdateCmd() *instructor {
	i := &instructor{}
	defaultAddr := fmt.Sprintf("%s:%d", defaultHost, defaultPort)
	i.cmd = flag.NewFlagSet("update", flag.ExitOnError)
	i.editor = i.cmd.String("editor", defaultAddr, "Address of the Phi Suite Schema Editor")
	i.inspector = i.cmd.String("inspector", "", "Address of the Phi Suite Schema Inspector")
	return i
}

func initInspectCmd() *instructor {
	i := &instructor{}
	defaultAddr := fmt.Sprintf("%s:%d", defaultHost, defaultPort)
	i.cmd = flag.NewFlagSet("inspect", flag.ExitOnError)
	i.inspector = i.cmd.String("inspector", defaultAddr, "Address of the Phi Suite Schema Inspector")
	return i
}

func (i *instructor) inspect() (err error) {
	err = i.cmd.Parse(os.Args[2:])
	if err != nil {
		return
	}
	args := i.cmd.Args()
	filename := ""
	if len(args) > 0 {
		filename = args[0]
	}
	if filename != "" && !strings.HasSuffix(filename, ".phi") {
		filename += ".phi"
	}
	inspector := &inspector{}
	err = inspector.inspect(*i.inspector)
	if err != nil {
		return
	}
	fmt.Println(inspector.kernel)
	return
}

func (i *instructor) update() (err error) {
	err = i.cmd.Parse(os.Args[2:])
	if err != nil {
		return
	}
	args := i.cmd.Args()
	filename := "main.phi"
	if len(args) > 0 {
		filename = args[0]
	}
	if !strings.HasSuffix(filename, ".phi") {
		filename += ".phi"
	}
	updater := &updater{}
	err = updater.update(filename, *i.editor, *i.inspector)
	if err != nil {
		return
	}
	fmt.Println(updater.kernel)
	return
}

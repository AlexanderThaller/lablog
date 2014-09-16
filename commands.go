package main

import (
	"fmt"
	"os"

	"github.com/AlexanderThaller/logger"
	"github.com/juju/errgo"
)

type Command struct {
	Type CommandType
	Args []string
	Config
}

type CommandType uint

const (
	CommandList CommandType = iota
	CommandNote
	CommandTrack
	CommandTodo
)

const (
	CommandListString  = "list"
	CommandNoteString  = "note"
	CommandTrackString = "track"
	CommandTodoString  = "todo"
)

func (typ CommandType) String() string {
	switch typ {
	case CommandList:
		return CommandListString
	case CommandNote:
		return CommandNoteString
	case CommandTrack:
		return CommandTrackString
	case CommandTodo:
		return CommandTodoString
	default:
		return "Unkown"
	}
}

func NewCommand(typ CommandType, args []string) Command {
	command := new(Command)
	command.Type = typ
	command.Args = args

	return *command
}

func parseCommand(args []string) (Command, error) {
	if len(args) == 1 {
		command := NewCommand(CommandList, args)
		return command, nil
	}

	var command Command
	var err error

	switch args[1] {
	case CommandListString:
		command = NewCommand(CommandList, args)
	case CommandNoteString:
		command = NewCommand(CommandNote, args)
	case CommandTrackString:
		command = NewCommand(CommandTrack, args)
	case CommandTodoString:
		command = NewCommand(CommandTodo, args)
	default:
		err = errgo.New("do not know the command " + args[1])
	}

	return command, err
}

func (com Command) Run() error {
	var err error

	switch com.Type {
	case CommandList:
		err = com.runList()
	case CommandNote:
		err = com.runNote()
	default:
		err = errgo.New("do not implement the command " + com.Type.String())
	}

	return err
}

func (com Command) runList() error {
	l := logger.New(Name, "Command", "runList")
	l.Debug("Args: ", com.Args)
	l.Debug("Args Len: ", len(com.Args))

	switch len(com.Args) {
	case 1, 2:
		return com.runListProjects()
	default:
		return errgo.New("do not know a list command with this parameter count")
	}
}

func (com Command) runListProjects() error {
	projects, err := GetProjects(com.Config.DataPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	for _, d := range projects {
		fmt.Println(d)
	}

	return nil
}

func (com Command) runNote() error {
	l := logger.New(Name, "Command", "runNote")
	l.Debug("Args: ", com.Args)
	l.Debug("Args Len: ", len(com.Args))

	switch len(com.Args) {
	case 4:
		project := com.Args[2]
		note := com.Args[3]
		return WriteProjectNote(com.DataPath, project, note)
	default:
		return errgo.New("do not know a note command with this parameter count")
	}
}

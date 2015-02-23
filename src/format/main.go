package format

import (
	"bytes"
	"io"
	"os/exec"
	"sort"
	"time"

	"github.com/AlexanderThaller/lablog/src/data"
	"github.com/AlexanderThaller/lablog/src/helper"
	"github.com/juju/errgo"
)

const (
	Name = "format"
)

const AsciiDocSettings = `:toc: right
:toclevels: 2
:sectanchors:
:sectlink:
:icons: font
:linkattrs:
:numbered:
:idprefix:
:idseparator: -
:doctype: book
:source-highlighter: pygments
:listing-caption: Listing`

func ProjectsEntries(writer io.Writer, projects []data.Project, start, end time.Time) error {
	io.WriteString(writer, AsciiDocSettings+"\n\n")
	io.WriteString(writer, "= Entries \n\n")

	for _, project := range projects {
		notes, err := helper.FilteredNotesByStartEnd(project, start, end)
		if err != nil {
			return errgo.Notef(err, "can not get filtered notes")
		}

		todos, err := helper.FilteredTodosByStartEnd(project, start, end)
		if err != nil {
			return errgo.Notef(err, "can not get filtered notes")
		}
		todos = data.FilterTodosLatest(todos)
		todos = data.FilterTodosAreNotDone(todos)

		project.Format(writer, 1)
		if len(todos) != 0 {
			Todos(writer, todos)
			io.WriteString(writer, "\n")
		}

		if len(notes) != 0 {
			Notes(writer, notes)
		}
	}

	return nil
}

func ProjectsNotes(writer io.Writer, projects []data.Project, start, end time.Time) error {
	io.WriteString(writer, AsciiDocSettings+"\n\n")
	io.WriteString(writer, "= Notes \n\n")

	for _, project := range projects {
		notes, err := helper.FilteredNotesByStartEnd(project, start, end)
		if err != nil {
			return errgo.Notef(err, "can not get filtered notes")
		}

		project.Format(writer, 1)
		Notes(writer, notes)
	}

	return nil
}

func ProjectsTodos(writer io.Writer, projects []data.Project, start, end time.Time) error {
	io.WriteString(writer, AsciiDocSettings+"\n\n")
	io.WriteString(writer, "= Todos \n\n")

	for _, project := range projects {
		todos, err := helper.FilteredTodosByStartEnd(project, start, end)
		if err != nil {
			return errgo.Notef(err, "can not get filtered notes")
		}
		todos = data.FilterTodosLatest(todos)
		todos = data.FilterTodosAreNotDone(todos)

		if len(todos) == 0 {
			continue
		}

		project.Format(writer, 1)
		Todos(writer, todos)
		io.WriteString(writer, "\n")
	}

	return nil
}

func ProjectsDates(writer io.Writer, projects []data.Project, start, end time.Time) error {
	io.WriteString(writer, AsciiDocSettings+"\n\n")
	io.WriteString(writer, "= Todos \n\n")

	dates, err := helper.ProjectsDates(projects, start, end)
	if err != nil {
		return errgo.Notef(err, "can not get dates for projects")
	}

	sort.Strings(dates)

	for _, date := range dates {
		io.WriteString(writer, "* "+date+"\n")
	}

	return nil
}

func Todos(writer io.Writer, todos []data.Todo) {
	io.WriteString(writer, "=== Todos\n\n")

	sort.Sort(data.TodosByName(todos))
	for _, todo := range todos {
		todo.Format(writer, 2)
	}
}

func Notes(writer io.Writer, notes []data.Note) {
	io.WriteString(writer, "=== Notes\n\n")

	notes = data.FilterNotesNotEmpty(notes)

	sort.Sort(data.NotesByTimeStamp(notes))
	for _, note := range notes {
		note.Format(writer, 2)
	}
}

func AsciiDoctor(reader io.Reader, writer io.Writer) error {
	stderr := new(bytes.Buffer)

	command := exec.Command("asciidoctor", "-")
	command.Stdin = reader
	command.Stdout = writer
	command.Stderr = stderr

	err := command.Run()
	if err != nil {
		return errgo.Notef(errgo.Notef(err, "can not run asciidoctor"),
			stderr.String())
	}

	return nil
}

func Log(writer io.Writer, projects []data.Project, start, end time.Time) error {
	return errgo.New("not implemented")
}
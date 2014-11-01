package main

import (
	"testing"

	"github.com/AlexanderThaller/logger"
)

func Test_RunNoDataPath(t *testing.T) {
	l := logger.New(Name, "Test", "Command", "Run", "NoDataPath")

	command := new(Command)
	err := command.Run()

	got := err.Error()
	expected := "the datapath can not be empty"

	testerr_output(t, l, err, got, expected)
}

func Test_RunDates(t *testing.T) {
	l := logger.New(Name, "Test", "Command", "Run", "Dates")

	action := ActionDates

	expected := `2014-10-31
2014-11-02
`

	testCommandRunOutput(t, l, action, expected)
}

func Test_RunList(t *testing.T) {
	l := logger.New(Name, "Test", "Command", "Run", "List")

	action := ActionList
	expected := `TestNotes
TestTodos
TestTracks
`

	testCommandRunOutput(t, l, action, expected)
}

func Test_RunNotes(t *testing.T) {
	l := logger.New(Name, "Test", "Command", "Run", "Notes")
	l.SetLevel(logger.Info)

	action := ActionNotes
	expected := "= Lablog -- Notes\n"
	expected += AsciiDocSettings + "\n\n"
	expected += `== TestNotes

=== 2014-10-31T21:36:31.49146148+01:00
Test1

=== 2014-10-31T21:36:33.49871531+01:00
Test2

=== 2014-10-31T21:36:35.138412374+01:00
Test3

=== 2014-10-31T21:36:36.810478305+01:00
Test4

=== 2014-10-31T21:36:38.450479686+01:00
Test5

`

	testCommandRunOutput(t, l, action, expected)
}

func Test_RunProjects(t *testing.T) {
	l := logger.New(Name, "Test", "Command", "Run", "Projects")

	action := ActionProjects
	expected := `TestNotes
TestTodos
TestTracks
`

	testCommandRunOutput(t, l, action, expected)
}

func Test_RunTodos(t *testing.T) {
	l := logger.New(Name, "Test", "Command", "Run", "Todos")
	l.SetLevel(logger.Info)

	action := ActionTodos
	expected := "= Lablog -- Todos\n"
	expected += AsciiDocSettings + "\n\n"
	expected += `== TestTodos

* Test1
* Test2
* Test3
* Test4
* Test5
* Test7

`

	testCommandRunOutput(t, l, action, expected)
}

func Test_RunTracks(t *testing.T) {
	l := logger.New(Name, "Test", "Command", "Run", "Tracks")
	l.SetLevel(logger.Info)

	action := ActionTracks
	expected := "= Lablog -- Tracks\n"
	expected += AsciiDocSettings + "\n\n"
	expected += `== TestTracks

* 2014-11-02T00:46:27.250010094+01:00 - Test1
* 2014-11-02T00:46:31.186052306+01:00 - Test2
* 2014-11-02T00:46:32.794131714+01:00 - Test3
* 2014-11-02T00:46:34.322047221+01:00 - Test4
* 2014-11-02T00:46:35.658221386+01:00 - Test5
* 2014-11-02T00:46:38.385921461+01:00
* 2014-11-02T00:46:51.833861322+01:00 - Test6
* 2014-11-02T00:46:53.713763133+01:00 - Test6
* 2014-11-02T00:46:55.609688812+01:00 - Test7
* 2014-11-02T00:46:56.713668447+01:00 - Test7
* 2014-11-02T00:46:57.953493565+01:00 - Test7

`

	testCommandRunOutput(t, l, action, expected)
}

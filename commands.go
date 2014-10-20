package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/AlexanderThaller/logger"
	"github.com/jinzhu/now"
	"github.com/juju/errgo"
)

type Command struct {
	Action        string
	Args          []string
	DataPath      string
	EndTime       time.Time
	Project       string
	SCM           string
	SCMAutoCommit bool
	SCMAutoPush   bool
	StartTime     time.Time
	TimeStamp     time.Time
	Value         string
}

const (
	CommitMessageTimeStampFormat = RecordTimeStampFormat
	DateFormat                   = "2006-01-02"
)

const (
	ActionDone         = "done"
	ActionList         = "list"
	ActionListDates    = "listdates"
	ActionListNotes    = "listnotes"
	ActionListProjects = "listprojects"
	ActionListTodos    = "listtodos"
	ActionListTracks   = "listtracks"
	ActionNote         = "note"
	ActionRename       = "rename"
	ActionTodo         = "todo"
	ActionMerge        = "merge"
	ActionTrack        = "track"
)

func NewCommand() *Command {
	return new(Command)
}

func (com *Command) Run() error {
	if com.DataPath == "" {
		return errgo.New("the datapath can not be empty")
	}

	switch com.Action {
	case ActionDone:
		return com.runDone()
	case ActionNote:
		return com.runNote()
	case ActionListDates:
		return com.runListDates()
	case ActionList:
		return com.runList()
	case ActionListNotes:
		return com.runListNotes()
	case ActionListProjects:
		return com.runListProjects()
	case ActionListTodos:
		return com.runListTodos()
	case ActionListTracks:
		return com.runListTracks()
	case ActionTodo:
		return com.runTodo()
	case ActionRename:
		return com.runRename()
	case ActionMerge:
		return com.runMerge()
	case ActionTrack:
		return com.runTrack()
	default:
		return errgo.New("Do not recognize the action: " + com.Action)
	}
}

func (com *Command) runTrack() error {
	if com.Project == "" {
		return errgo.New("track command needs an project")
	}

	track := new(Track)
	track.Project = com.Project
	track.TimeStamp = com.TimeStamp
	track.Value = com.Value

	return com.Write(track)
}

func (com *Command) runMerge() error {
	if com.Project == "" {
		return errgo.New("Project name can not be empty")
	}
	srcproject := com.Project
	dstproject := com.Value

	if !com.checkProjectExists(srcproject) {
		return errgo.New("no project with the name " + srcproject)
	}

	if !com.checkProjectExists(dstproject) {
		return errgo.New("the project " + dstproject + " already exists")
	}

	srcpath := path.Join(com.DataPath, srcproject+".csv")
	dstpath := path.Join(com.DataPath, dstproject+".csv")

	err := com.runMergeFiles(srcpath, dstpath)
	if err != nil {
		return err
	}

	srcfile := srcproject + ".csv"
	err = scmRemove(com.SCM, srcfile, com.DataPath)
	if err != nil {
		return err
	}

	dstfile := dstproject + ".csv"
	err = scmAdd(com.SCM, com.DataPath, dstfile)
	if err != nil {
		return err
	}

	message := srcproject + " - merged - " + dstproject
	err = scmCommit(com.SCM, com.DataPath, message)
	if err != nil {
		return err
	}

	return nil
}

func (com *Command) runMergeFiles(srcpath, dstpath string) error {
	srcdata, err := ioutil.ReadFile(srcpath)
	if err != nil {
		return errgo.New(err.Error())
	}

	dstfile, err := os.OpenFile(dstpath, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return errgo.New(err.Error())
	}
	defer dstfile.Close()

	_, err = dstfile.Write(srcdata)
	if err != nil {
		return errgo.New(err.Error())
	}

	return nil
}

func (com *Command) runRename() error {
	if com.Project == "" {
		return errgo.New("Project name can not be empty")
	}
	oldproject := com.Project
	newproject := com.Value

	if !com.checkProjectExists(oldproject) {
		return errgo.New("no project with the name " + oldproject)
	}

	if com.checkProjectExists(newproject) {
		return errgo.New("the project " + newproject + " already exists")
	}

	oldpath := oldproject + ".csv"
	newpath := newproject + ".csv"

	err := scmRename(com.SCM, oldpath, newpath, com.DataPath)
	if err != nil {
		return err
	}

	message := oldproject + " - renamed - " + newproject
	err = scmCommit(com.SCM, com.DataPath, message)
	if err != nil {
		return err
	}

	return nil
}

func (com *Command) runDone() error {
	l := logger.New(Name, "Command", "run", "Done")

	l.Trace("Args length: ", len(com.Args))
	if com.Value == "" {
		return errgo.New("todo command needs a value")
	}
	l.Trace("Project: ", com.Project)
	if com.Project == "" {
		return errgo.New("todo command needs an project")
	}

	done := new(Todo)
	done.Project = com.Project
	done.TimeStamp = com.TimeStamp
	done.Value = com.Value
	done.Done = true
	l.Trace("Done: ", fmt.Sprintf("%+v", done))

	return com.Write(done)
}

func (com *Command) runNote() error {
	l := logger.New(Name, "Command", "run", "Note")

	l.Trace("Args length: ", len(com.Args))
	if com.Value == "" {
		return errgo.New("note command needs a value")
	}
	l.Trace("Project: ", com.Project)
	if com.Project == "" {
		return errgo.New("note command needs an project")
	}

	note := new(Note)
	note.Project = com.Project
	note.TimeStamp = com.TimeStamp
	note.Value = com.Value
	l.Trace("Note: ", fmt.Sprintf("%+v", note))

	return com.Write(note)
}

func (com *Command) runTodo() error {
	l := logger.New(Name, "Command", "run", "Todo")

	l.Trace("Args length: ", len(com.Args))
	if com.Value == "" {
		return errgo.New("todo command needs a value")
	}
	l.Trace("Project: ", com.Project)
	if com.Project == "" {
		return errgo.New("todo command needs an project")
	}

	todo := new(Todo)
	todo.Project = com.Project
	todo.TimeStamp = com.TimeStamp
	todo.Value = com.Value
	todo.Done = false
	l.Trace("Todo: ", fmt.Sprintf("%+v", todo))

	return com.Write(todo)
}

func (com *Command) runList() error {
	if com.Project == "" {
		return com.runListProjects()
	}

	if !com.checkProjectExists(com.Project) {
		return errgo.New("project " + com.Project + " does not exist")
	}

	notes, err := com.getProjectNotes(com.Project)
	if err != nil {
		return err
	}
	if len(notes) != 0 {
		return com.runListProjectNotes(com.Project)
	}

	return com.runListProjectTodos(com.Project)
}

func (com *Command) runListNotes() error {
	if com.Project != "" {
		return com.runListProjectNotes(com.Project)
	}

	projects, err := com.getProjects()
	if err != nil {
		return err
	}

	for _, project := range projects {
		err := com.runListProjectNotes(project)
		if err != nil {
			return err
		}
	}

	return nil
}

func (com *Command) runListTodos() error {
	if com.Project != "" {
		return com.runListProjectTodos(com.Project)
	}

	projects, err := com.getProjects()
	if err != nil {
		return err
	}

	for _, project := range projects {
		err := com.runListProjectTodos(project)
		if err != nil {
			return err
		}
	}

	return nil
}

func (com *Command) runListTracks() error {
	if com.Project != "" {
		return com.runListProjectTracks(com.Project)
	}

	projects, err := com.getProjects()
	if err != nil {
		return err
	}

	for _, project := range projects {
		err := com.runListProjectTracks(project)
		if err != nil {
			return err
		}
	}

	return nil
}

func (com *Command) runListProjects() error {
	projects, err := com.getProjects()
	if err != nil {
		return err
	}

	out := make(map[string]struct{})
	for _, project := range projects {
		notes, err := com.getProjectNotes(project)
		if err != nil {
			return err
		}

		if len(notes) == 0 {
			continue
		}

		out[project] = struct{}{}
	}

	for _, project := range projects {
		todos, err := com.getProjectTodos(project)
		if err != nil {
			return err
		}

		if len(todos) == 0 {
			continue
		}

		out[project] = struct{}{}
	}

	var outsort []string
	for project := range out {
		outsort = append(outsort, project)
	}
	sort.Strings(outsort)

	for _, project := range outsort {
		fmt.Println(project)
	}

	return nil
}

func (com *Command) runListDates() error {
	l := logger.New(Name, "Command", "run", "ListDates")

	var dates []string
	var err error

	if com.Project == "" {
		dates, err = com.getDates()
	} else {
		dates, err = com.getProjectDates(com.Project)
	}

	if err != nil {
		return err
	}

	sort.Strings(dates)
	for _, date := range dates {
		timestamp, err := now.Parse(date)
		if err != nil {
			l.Warning("Can not parse timestamp: ", errgo.Details(err))
			continue
		}

		if timestamp.Before(com.StartTime) {
			continue
		}

		if timestamp.After(com.EndTime) {
			continue
		}

		fmt.Println(date)
	}

	return nil
}

func (com *Command) getDates() ([]string, error) {
	projects, err := com.getProjects()
	if err != nil {
		return nil, err
	}

	datemap := make(map[string]struct{})
	for _, project := range projects {
		dates, err := com.getProjectDates(project)
		if err != nil {
			return nil, err
		}

		for _, date := range dates {
			datemap[date] = struct{}{}
		}
	}

	var out []string
	for date := range datemap {
		out = append(out, date)
	}

	return out, nil
}

func (com *Command) getProjectDates(project string) ([]string, error) {
	if com.DataPath == "" {
		return nil, errgo.New("datapath can not be empty")
	}
	if project == "" {
		return nil, errgo.New("project name can not be empty")
	}
	if !com.checkProjectExists(project) {
		return nil, errgo.New("project does not exist")
	}

	var out []string

	notes, err := com.getProjectNotes(project)
	if err != nil {
		return nil, err
	}

	todos, err := com.getProjectTodos(project)
	if err != nil {
		return nil, err
	}
	todos = com.filterTodos(todos)

	datemap := make(map[string]struct{})

	for _, note := range notes {
		date, err := time.Parse(RecordTimeStampFormat, note.GetTimeStamp())
		if err != nil {
			return nil, err
		}

		datemap[date.Format(DateFormat)] = struct{}{}
	}

	for _, todo := range todos {
		datemap[todo.TimeStamp.Format(DateFormat)] = struct{}{}
	}

	for date := range datemap {
		out = append(out, date)
	}

	return out, nil
}

func (com *Command) getProjects() ([]string, error) {
	dir, err := ioutil.ReadDir(com.DataPath)
	if err != nil {
		return nil, err
	}

	var out []string
	for _, d := range dir {
		file := d.Name()

		// Skip dotfiles
		if strings.HasPrefix(file, ".") {
			continue
		}

		ext := filepath.Ext(file)
		name := file[0 : len(file)-len(ext)]

		out = append(out, name)
	}

	sort.Strings(out)
	return out, nil
}

func (com *Command) runListProjectTodos(project string) error {
	if project == "" {
		return errgo.New("project name can not be empty")
	}
	if !com.checkProjectExists(project) {
		return errgo.New("the project does not exist")
	}

	todos, err := com.getProjectTodos(project)
	if err != nil {
		return err
	}
	todos = com.filterTodos(todos)

	if len(todos) == 0 {
		return nil
	}

	fmt.Println("#", project, "- Todos")

	var out []string
	for _, todo := range todos {
		out = append(out, "  * "+todo.GetValue())
	}
	sort.Strings(out)

	for _, todo := range out {
		fmt.Println(todo)
	}
	fmt.Println("")

	return nil
}

func (com *Command) runListProjectTracks(project string) error {
	if project == "" {
		return errgo.New("project name can not be empty")
	}
	if !com.checkProjectExists(project) {
		return errgo.New("the project does not exist")
	}

	tracks, err := com.getProjectTracks(project)
	if err != nil {
		return err
	}

	if len(tracks) == 0 {
		return nil
	}

	fmt.Println("#", project, "- Tracks")

	for _, track := range tracks {
		fmt.Println("  *", track.TimeStamp, "-", track.Value)
	}
	fmt.Println("")

	return nil
}

func (com *Command) filterTodos(todos []Todo) []Todo {
	filter := make(map[string]Todo)

	sort.Sort(TodoByDate(todos))
	for _, todo := range todos {
		filter[todo.Value] = todo
	}

	var out []Todo
	for _, todo := range filter {
		if todo.Done {
			continue
		}

		out = append(out, todo)
	}

	return out
}

func (com *Command) runListProjectNotes(project string) error {
	if project == "" {
		return errgo.New("project name can not be empty")
	}
	if !com.checkProjectExists(project) {
		return errgo.New("project" + project + " does not exist")
	}

	notes, err := com.getProjectNotes(project)
	if err != nil {
		return err
	}

	if len(notes) == 0 {
		return nil
	}

	fmt.Println("#", project)
	sort.Sort(NotesByDate(notes))

	reg, err := regexp.Compile("(?m)^#")
	if err != nil {
		return err
	}

	for _, note := range notes {
		fmt.Println("##", note.GetTimeStamp())

		out := reg.ReplaceAllString(note.GetValue(), "###")
		fmt.Println(out)
		fmt.Println("")
	}

	return nil
}

func (com *Command) checkProjectExists(project string) bool {
	projects, err := com.getProjects()
	if err != nil {
		return false
	}

	for _, d := range projects {
		if d == project {
			return true
		}
	}

	return false
}

func (com *Command) getProjectNotes(project string) ([]Note, error) {
	l := logger.New(Name, "Command", "get", "ProjectRecords")
	l.SetLevel(logger.Debug)

	if com.DataPath == "" {
		return nil, errgo.New("datapath can not be empty")
	}
	if project == "" {
		return nil, errgo.New("project name can not be empty")
	}
	if !com.checkProjectExists(project) {
		return nil, errgo.New("project does not exist")
	}

	filepath := filepath.Join(com.DataPath, project+".csv")
	file, err := os.OpenFile(filepath, os.O_RDONLY, 0640)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 3

	var out []Note
	for {
		csv, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			continue
		}

		note, err := NoteFromCSV(csv)
		if err != nil {
			continue
		}
		note.SetProject(project)

		if note.TimeStamp.Before(com.StartTime) {
			continue
		}

		if note.TimeStamp.After(com.EndTime) {
			continue
		}

		out = append(out, note)
	}

	return out, err
}

func (com *Command) getProjectTodos(project string) ([]Todo, error) {
	l := logger.New(Name, "Command", "get", "ProjectRecords")
	l.SetLevel(logger.Debug)

	if com.DataPath == "" {
		return nil, errgo.New("datapath can not be empty")
	}
	if project == "" {
		return nil, errgo.New("project name can not be empty")
	}
	if !com.checkProjectExists(project) {
		return nil, errgo.New("project does not exist")
	}

	filepath := filepath.Join(com.DataPath, project+".csv")
	file, err := os.OpenFile(filepath, os.O_RDONLY, 0640)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 4

	var out []Todo
	for {
		csv, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			continue
		}

		todo, err := TodoFromCSV(csv)
		if err != nil {
			continue
		}

		out = append(out, todo)
	}

	return out, err
}

func (com *Command) getProjectTracks(project string) ([]Track, error) {
	l := logger.New(Name, "Command", "get", "ProjectTracks")
	l.SetLevel(logger.Debug)

	if com.DataPath == "" {
		return nil, errgo.New("datapath can not be empty")
	}
	if project == "" {
		return nil, errgo.New("project name can not be empty")
	}
	if !com.checkProjectExists(project) {
		return nil, errgo.New("project does not exist")
	}

	filepath := filepath.Join(com.DataPath, project+".csv")
	file, err := os.OpenFile(filepath, os.O_RDONLY, 0640)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 3

	var out []Track
	for {
		csv, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			continue
		}

		track, err := TrackFromCSV(csv)
		if err != nil {
			continue
		}

		out = append(out, track)
	}

	return out, err
}

func (com *Command) Write(record Record) error {
	if com.DataPath == "" {
		return errgo.New("datapath can not be empty")
	}

	if com.Project == "" {
		return errgo.New("project name can not be empty")
	}

	path := com.DataPath
	project := com.Project

	err := os.MkdirAll(path, 0750)
	if err != nil {
		return err
	}

	filepath := filepath.Join(path, project+".csv")
	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0640)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	err = writer.Write(record.CSV())
	if err != nil {
		return err
	}
	writer.Flush()

	err = com.Commit(record)
	if err != nil {
		return err
	}

	return nil
}

func (com *Command) Commit(record Record) error {
	if !com.SCMAutoCommit {
		return nil
	}

	if com.SCM == "" {
		return errgo.New("Can not use an empty scm for commiting")
	}

	filename := record.GetProject() + ".csv"
	err := scmAdd(com.SCM, com.DataPath, filename)
	if err != nil {
		return err
	}

	message := com.Project + " - "
	message += record.GetAction() + " - "
	message += com.TimeStamp.Format(CommitMessageTimeStampFormat)
	err = scmCommit(com.SCM, com.DataPath, message)
	if err != nil {
		return err
	}

	return nil
}

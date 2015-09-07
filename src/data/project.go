package data

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/juju/errgo"
)

// Projects returns all projects which can be read from the given datadir.
func GetProjects(datadir string) ([]Project, error) {
	names, err := projectNames(datadir)
	if err != nil {
		return nil, errgo.Notef(err, "can not get project files")
	}

	var projects []Project
	for _, name := range names {
		project := Project{
			Name:    name,
			Datadir: datadir,
		}

		projects = append(projects, project)
	}

	return projects, nil
}

// Project represents a project.
type Project struct {
	Name    string
	Datadir string
}

func (proj Project) IsEmpty() bool {
	filepath := proj.FilePath()
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return true
	}

	return false
}

func (proj Project) FilePath() string {
	return path.Join(proj.Datadir, proj.Name+".csv")
}

// ProjectsByName allows sorting project slices by name.
type ProjectsByName []Project

func (by ProjectsByName) Len() int {
	return len(by)
}

func (by ProjectsByName) Swap(i, j int) {
	by[i], by[j] = by[j], by[i]
}

func (by ProjectsByName) Less(i, j int) bool {
	return by[i].Name < by[j].Name
}

const (
	projectFileExtention = ".csv"
)

// projectNames read the files in the datadir and returns every filename
// (without the extention) that ends with projectFileExtention and does not
// start with a dot.
func projectNames(datadir string) ([]string, error) {
	dir, err := ioutil.ReadDir(datadir)
	if err != nil {
		return nil, errgo.Notef(err, "can not read files from directory")
	}

	var names []string
	for _, file := range dir {
		filename := file.Name()

		// Skip dotfiles
		if strings.HasPrefix(filename, ".") {
			continue
		}

		// Skip file not ending with right extention
		if !strings.HasSuffix(filename, projectFileExtention) {
			continue
		}

		ext := filepath.Ext(filename)

		// Strip extention from filename to get only the name
		name := filename[0 : len(filename)-len(ext)]

		names = append(names, name)
	}

	return names, nil
}

func (project Project) Notes() ([]Note, error) {
	file, err := project.File()
	if err != nil {
		return nil, errgo.Notef(err, "can not open project file")
	}
	defer file.Close()

	parser := csv.NewReader(file)
	parser.FieldsPerRecord = 3

	var out []Note

	for {
		csv, err := parser.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			// an error would mean that the csv line is not a note so we can skip it
			continue
		}

		if csv[1] != "note" {
			continue
		}

		note, err := ParseNote(project, csv)
		if err != nil {
			return nil, errgo.Notef(err, "can not parse note from csv")
		}

		out = append(out, note)
	}

	return out, nil
}

func (project Project) Todos() ([]Todo, error) {
	file, err := project.File()
	if err != nil {
		return nil, errgo.Notef(err, "can not open project file")
	}
	defer file.Close()

	parser := csv.NewReader(file)
	parser.FieldsPerRecord = 4

	var out []Todo

	for {
		csv, err := parser.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			// an error would mean that the csv line is not a todo so we can skip it
			continue
		}

		if csv[1] != "todo" {
			continue
		}

		todo, err := ParseTodo(project, csv)
		if err != nil {
			return nil, errgo.Notef(err, "can not parse todo from csv")
		}

		out = append(out, todo)
	}

	return out, nil
}

func (project Project) Tracks() ([]Track, error) {
	file, err := project.File()
	if err != nil {
		return nil, errgo.Notef(err, "can not open project file")
	}
	defer file.Close()

	parser := csv.NewReader(file)
	parser.FieldsPerRecord = 3

	var out []Track

	for {
		csv, err := parser.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			// an error would mean that the csv line is not a track so we can skip it
			continue
		}

		if csv[1] != "track" {
			continue
		}

		track, err := ParseTrack(project, csv)
		if err != nil {
			return nil, errgo.Notef(err, "can not parse todo from csv")
		}

		out = append(out, track)
	}

	return out, nil
}

func (project Project) File() (*os.File, error) {
	if project.Datadir == "" {
		return nil, errgo.New("path to datadir can not be emtpy")
	}

	filepath := filepath.Join(project.Datadir, project.Name+".csv")
	file, err := os.OpenFile(filepath, os.O_RDONLY, 0640)
	if err != nil {
		return nil, errgo.Notef(err, "can not open file")
	}

	return file, nil
}

func (project Project) Format(writer io.Writer, indent uint) {
	indentchar := strings.Repeat("=", int(indent))
	io.WriteString(writer, indentchar+"= "+project.Name+"\n")
}

func ProjectsOrArgs(args []string, datadir string) (Projects, error) {
	var projects []Project

	if len(args) == 0 {
		var err error
		projects, err = GetProjects(datadir)
		if err != nil {
			return nil, errgo.Notef(err, "can not get projects")
		}
	} else {
		for _, arg := range args {
			project := Project{
				Name:    arg,
				Datadir: datadir,
			}

			projects = append(projects, project)
		}
	}

	return Projects(projects), nil
}

func (project Project) IsActive() bool {
	tracks, err := project.Tracks()
	if err != nil {
		return false
	}

	sort.Sort(TracksByTimeStamp(tracks))

	active := tracks[len(tracks)-1].Active

	return active
}

func (project Project) IsSubproject(subproject Project) bool {
	if project.Name == subproject.Name {
		return false
	}

	return strings.HasPrefix(subproject.Name, project.Name)
}

func FilterProjectsSubprojects(allprojects []Project, projects []Project) []Project {
	out := FilterEmptyProjects(projects)

	for _, project := range projects {
		for _, subproject := range allprojects {
			if project.IsSubproject(subproject) {
				out = append(out, subproject)
			}
		}
	}

	return out
}

func Subprojects(args []string, datadir string, subprojects bool) ([]Project, error) {
	if !subprojects {
		projects, err := ProjectsOrArgs(args, datadir)
		if err != nil {
			return nil, errgo.Notef(err, "can not get projects")
		}

		return projects, nil
	}

	projects, err := ProjectsOrArgs(args, datadir)
	if err != nil {
		return nil, errgo.Notef(err, "can not get main projects")
	}

	allprojects, err := GetProjects(datadir)
	if err != nil {
		return nil, errgo.Notef(err, "can not get list of all projects")
	}

	projects = FilterProjectsSubprojects(allprojects, projects)
	projects = FilterDoubleProjects(projects)

	return projects, nil
}

func FilterDoubleProjects(projects []Project) []Project {
	filter := make(map[string]Project)

	for _, project := range projects {
		filter[project.Name] = project
	}

	var out []Project
	for _, project := range filter {
		out = append(out, project)
	}

	return out
}

func FilterEmptyProjects(projects []Project) []Project {
	var out []Project

	for _, project := range projects {
		if project.IsEmpty() {
			continue
		}

		out = append(out, project)
	}

	return out
}

type Projects []Project

func (project Project) Search(searchstr string) ([]string, error) {
	var out []string
	notes, err := project.Notes()
	if err != nil {
		return nil, errgo.Notef(err, "can not get notes for project "+project.Name)
	}

	for _, note := range notes {
		buff := bytes.NewBuffer([]byte(note.Text))

		scanner := bufio.NewScanner(buff)
		for scanner.Scan() {
			text := scanner.Text()
			if strings.Contains(text, searchstr) {
				out = append(out, text)
			}
		}
	}

	return out, nil
}

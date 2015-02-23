package commands

import (
	"os"
	"strings"
	"time"

	"github.com/AlexanderThaller/cobra"
	"github.com/AlexanderThaller/lablog/src/data"
	"github.com/AlexanderThaller/lablog/src/helper"
	"github.com/AlexanderThaller/logger"
)

var cmdNote = &cobra.Command{
	Use:     "note [project] [text]",
	Short:   "Create a new note for the project.",
	Long:    `Create a note which will record the current timestamp and the given text for the given project.`,
	Run:     runNote,
	PostRun: finished,
}

var flagNoteTimeStamp time.Time
var flagNoteTimeStampRaw string

func init() {
	flagNoteTimeStamp = time.Now()

	cmdNote.Flags().StringVarP(&flagNoteTimeStampRaw, "timestamp", "t",
		flagNoteTimeStamp.String(), "The timestamp for which to record the note.")
}

func runNote(cmd *cobra.Command, args []string) {
	l := logger.New("commands", "note")

	if len(args) < 2 {
		l.Alert("need at least two arguments to run")
		os.Exit(1)
	}

	project := args[0]
	text := strings.Join(args[1:], " ")

	timestamp, err := helper.DefaultOrRawTimestamp(flagNoteTimeStamp, flagNoteTimeStampRaw)
	errexit(l, err, "can not get timestamp")

	note := data.Note{
		Project:   data.Project{Name: project},
		TimeStamp: timestamp,
		Text:      text,
	}

	l.Trace("Note: ", note)
	recordAndCommit(l, flagLablogDataDir, note)
}
package scm

import (
	"os/exec"
	"time"

	"github.com/juju/errgo"
)

func Commit(scm, datapath, message string) error {
	switch scm {
	case "git":
		return gitCommit(datapath, message)
	default:
		return errgo.New("do not know the scm " + scm)
	}
}

func gitCommit(datapath, message string) error {
	command := exec.Command("git", "commit", "-m", message)
	command.Dir = datapath

	output, err := command.CombinedOutput()
	if err != nil {
		err = errgo.New("problem when commiting to git: " + err.Error() + " - " +
			string(output))

		return err
	}

	// Give git time to commit everything and remove the lockfile.
	time.Sleep(5 * time.Millisecond)
	return nil
}

// TODO: Change this to scm, filename, datapath
func Add(scm, datapath, filename string) error {
	switch scm {
	case "git":
		return gitAdd(datapath, filename)
	default:
		return errgo.New("do not know the scm " + scm)
	}
}

func gitAdd(datapath, filename string) error {
	command := exec.Command("git", "add", filename)
	command.Dir = datapath

	output, err := command.CombinedOutput()
	if err != nil {
		err = errgo.New("problem when adding to git: " + err.Error() + " - " +
			string(output))

		return err
	}

	// Give git time to commit everything and remove the lockfile.
	time.Sleep(5 * time.Millisecond)
	return nil
}

func Push(scm, datapath string) error {
	switch scm {
	case "git":
		return gitPush(datapath)
	default:
		return errgo.New("do not know the scm " + scm)
	}
}

func gitPush(datapath string) error {
	command := exec.Command("git", "push")
	command.Dir = datapath

	output, err := command.CombinedOutput()
	if err != nil {
		err = errgo.New("problem when pushing to git: " + err.Error() + " - " +
			string(output))

		return err
	}

	// Give git time to commit everything and remove the lockfile.
	time.Sleep(5 * time.Millisecond)
	return nil
}

func Rename(scm, oldpath, newpath, datapath string) error {
	switch scm {
	case "git":
		return gitRename(oldpath, newpath, datapath)
	default:
		return errgo.New("do not know the scm " + scm)
	}
}

func gitRename(oldpath, newpath, datapath string) error {
	command := exec.Command("git", "mv", oldpath, newpath)
	command.Dir = datapath

	output, err := command.CombinedOutput()
	if err != nil {
		err = errgo.New("problem when moving file in git: " + err.Error() + " - " +
			string(output))

		return err
	}

	// Give git time
	time.Sleep(5 * time.Millisecond)
	return nil
}

func Remove(scm, filename, datapath string) error {
	switch scm {
	case "git":
		return gitRemove(filename, datapath)
	default:
		return errgo.New("do not know the scm " + scm)
	}
}

func gitRemove(filename, datapath string) error {
	command := exec.Command("git", "rm", filename)
	command.Dir = datapath

	output, err := command.CombinedOutput()
	if err != nil {
		err = errgo.New("problem when moving file in git: " + err.Error() + " - " +
			string(output))

		return err
	}

	// Give git time
	time.Sleep(5 * time.Millisecond)
	return nil
}

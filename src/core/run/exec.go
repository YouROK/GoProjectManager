package run

import (
	"io"
	"os"
	"os/exec"
)

type Exec struct {
	StdOut io.Writer
	StdErr io.Writer
	StdIn  io.Reader
	cmd    *exec.Cmd
}

func NewExec(cmd string, args []string, envs []string) *Exec {
	e := &Exec{}
	e.cmd = exec.Command(cmd, args...)
	e.cmd.Env = envs
	e.cmd.Env = append(e.cmd.Env, "PATH="+os.Getenv("PATH"))
	return e
}

func (e *Exec) Run() error {
	e.cmd.Stdout = e.StdOut
	e.cmd.Stderr = e.StdErr
	e.cmd.Stdin = e.StdIn
	return e.cmd.Run()
}

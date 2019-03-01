package shell

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	Shell       = []string{"/bin/sh", "-c"}
	Panic       = true
	Trace       = false
	TracePrefix = "+"

	exit = os.Exit
)

func Path(parts ...string) string {
	return filepath.Join(parts...)
}

func PathTemplate(parts ...string) func(...interface{}) string {
	return func(values ...interface{}) string {
		return fmt.Sprintf(Path(parts...), values...)
	}
}

func Quote(arg string) string {
	return fmt.Sprintf("'%s'", strings.Replace(arg, "'", "'\\''", -1))
}

func ErrExit() {
	if p, ok := recover().(*Process); p != nil {
		if !ok {
			fmt.Fprintf(os.Stderr, "Unexpected panic: %v\n", p)
			exit(1)
		}
		fmt.Fprintf(os.Stderr, "%s\n", p.Error())
		exit(p.ExitStatus)
	}
}

func Cmd(cmd ...interface{}) *Command {
	c := new(Command)
	c.addArgs(cmd...)
	return c
}

func Run(cmd ...interface{}) *Process {
	return Cmd(cmd...).Run()
}

func Start(cmd ...interface{}) *Process {
	return Cmd(cmd...).Start()
}

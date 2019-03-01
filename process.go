package shell

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"syscall"
)

type Process struct {
	cmd    *exec.Cmd
	killed bool

	Stdout     *bytes.Buffer
	Stderr     *bytes.Buffer
	Stdin      io.WriteCloser
	ExitStatus int
}

func (p *Process) Wait() error {
	err := p.cmd.Wait()
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if stat, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				p.ExitStatus = int(stat.ExitStatus())
				fmt.Printf("%v\n", err)
				if Panic && !p.killed {
					panic(p)
				}
			}
		}
	}
	return err
}

func (p *Process) Kill() error {
	p.killed = true
	err := p.cmd.Process.Kill()
	if err != nil {
		return fmt.Errorf("killed error: %s", err)
	}
	if err := p.Wait(); err == nil {
		if !strings.Contains(err.Error(), "signal: killed") {
			return err
		}
	}
	return nil
}

func (p *Process) String() string {
	return strings.Trim(p.Stdout.String(), "\n")
}

func (p *Process) Bytes() []byte {
	return p.Stdout.Bytes()
}

func (p *Process) Error() error {
	errlines := strings.Split(p.Stderr.String(), "\n")
	s := len(errlines)
	if s > 1 {
		s -= 1
	} else {
		s = 0
	}
	return fmt.Errorf("[%v] %s\n", p.ExitStatus, errlines[s])
}

func (p *Process) Read(b []byte) (int, error) {
	return p.Stdout.Read(b)
}

func (p *Process) Write(b []byte) (int, error) {
	return p.Stdin.Write(b)
}

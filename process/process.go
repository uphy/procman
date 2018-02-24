package process

import (
	"bufio"
	"fmt"
	"os/exec"
	"sync"
	"syscall"
)

type (
	Process struct {
		command    []string
		cmd        *exec.Cmd
		state      State
		mutex      *sync.Mutex
		stopSignal chan struct{}
	}
	State string
)

const (
	StateStarted State = "started"
	StateError   State = "error"
	StateStopped State = "stopped"
)

func New(command []string) *Process {
	return &Process{
		command: command,
		state:   StateStopped,
		mutex:   new(sync.Mutex),
	}
}

func (p *Process) State() State {
	return p.state
}

func (p *Process) createProcess() {
	p.cmd = exec.Command(p.command[0], p.command[1:]...)
	p.cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

func (p *Process) Start() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if p.state == StateStarted {
		return nil
	}
	p.state = StateStarted
	p.createProcess()
	stdout, err := p.cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := p.cmd.StderrPipe()
	if err != nil {
		return err
	}
	if err := p.cmd.Start(); err != nil {
		return err
	}
	go func() {
		stdoutScanner := bufio.NewScanner(stdout)
		stderrScanner := bufio.NewScanner(stderr)
		go func() {
			for stdoutScanner.Scan() {
				fmt.Println(stdoutScanner.Text())
			}
		}()
		go func() {
			for stderrScanner.Scan() {
				fmt.Println(stderrScanner.Text())
			}
		}()
		p.cmd.Wait()

		p.mutex.Lock()
		defer p.mutex.Unlock()
		p.state = StateStopped
	}()
	return nil
}

func (p *Process) Stop() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if p.state != StateStarted {
		return nil
	}
	pgid, err := syscall.Getpgid(p.cmd.Process.Pid)
	if err != nil {
		return err
	}
	return syscall.Kill(-pgid, syscall.SIGTERM)
}

func (p *Process) Restart() error {
	if err := p.Stop(); err != nil {
		return err
	}
	return p.Start()
}

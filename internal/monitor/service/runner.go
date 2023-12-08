package service

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"os/exec"
	"strconv"
)

const (
	golangExecutable = "/bin/go.exe"
	serverEntryPoint = "/cmd/server/main.go"
)

var (
	goRoot     = os.Getenv("GOROOT")
	workDir, _ = os.Getwd()
)

type Runner interface {
	RunInstance() error
	Shutdown() error
	Reboot(serviceName string) error
}

type runner struct {
	logger    *zerolog.Logger
	processes map[string]*os.Process
}

func NewRunner(logger *zerolog.Logger) Runner {
	return &runner{
		logger:    logger,
		processes: make(map[string]*os.Process),
	}
}

func (r *runner) RunInstance() error {
	cmd := exec.Command(goRoot+golangExecutable, "run", workDir+serverEntryPoint)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("cmd.Start: %w", err)
	}

	if err := cmd.Err; err != nil {
		return fmt.Errorf("cmd.Err: %w", err)
	}

	process := cmd.Process
	if process != nil {
		r.logger.Debug().Msgf("Add process (%d) to tracked", process.Pid)
		r.processes["server"] = process
	}

	return nil
}

func (r *runner) Reboot(serviceName string) error {
	process, ok := r.processes[serviceName]
	if !ok {
		return fmt.Errorf("%s process not found in running", serviceName)
	}

	if err := r.killProcess(process.Pid); err != nil {
		return fmt.Errorf("r.killProcess: %w", err)
	}
	r.logger.Debug().Msgf("Successfully kill process (%d)", process.Pid)

	if err := r.RunInstance(); err != nil {
		return fmt.Errorf("r.RunInstance: %w", err)
	}
	r.logger.Debug().Msg("Runner successfully start new server instance")

	return nil
}

func (r *runner) Shutdown() error {
	for key, process := range r.processes {
		if err := r.killProcess(process.Pid); err != nil {
			return fmt.Errorf("r.killProcess: %w", err)
		}
		delete(r.processes, key)
		r.logger.Debug().Msgf("Successfully kill process (%d)", process.Pid)
	}

	return nil
}

func (r *runner) killProcess(pid int) error {
	kill := exec.Command("taskkill", "/T", "/F", "/PID", strconv.Itoa(pid))

	if err := kill.Run(); err != nil {
		return fmt.Errorf("kill.Run: %w", err)
	}

	return nil
}

package rollover

import (
	"github.com/pkg/errors"
	"log"
	"os"
	"strconv"
	"syscall"
	"time"
)

const EnvKeyParentPID = "ROLLOVER_PARENT_PID"

// Wait if there is any parent process, kill and wait parent exit
// program should always perform this safe check before executing anything else
// to avoid lock, permission, or other issues
func Wait() (parent *os.Process, err error) {
	// check if there is a parent process
	if parent, ok := os.LookupEnv(EnvKeyParentPID); ok {
		pid, err := strconv.Atoi(parent)
		if err != nil {
			return nil, errors.New("error reading rollover parent pid")
		}
		proc, err := os.FindProcess(pid)
		if err != nil {
			return nil, errors.Errorf("can't find the parent process by pid %d", pid)
		}

		// send kill single
		err = proc.Signal(os.Interrupt)
		if err != nil {
			return nil, errors.WithMessagef(err, "failed to send INT to parent pid %d", pid)
		}

		// wait till parent exit
		ticker := time.NewTicker(500 * time.Millisecond)
		for {
			_ = <-ticker.C
			// only way to check process exist is https://stackoverflow.com/questions/15204162/check-if-a-process-exists-in-go-way
			err := proc.Signal(syscall.Signal(0))
			if err != nil {
				// no such process
				return proc, nil
			}
		}
	}
	return nil, nil
}

// Restart the current executable binary
func Restart() (child *os.Process, err error) {
	ex, err := os.Executable()
	if err != nil {
		log.Fatalln(err)
	}
	// execPath := filepath.Dir(ex)

	workdir, err := os.Getwd()
	if nil != err {
		log.Fatalln(err)
	}

	pid := os.Getpid()

	export := EnvKeyParentPID+"="+strconv.Itoa(pid)
        envs := os.Environ()
	replaced := false
	for k, v := range envs {
		if strings.HasPrefix(v, EnvKeyParentPID) {
			envs[k] = export
			replaced = true
			break
		}
	}
	if !replaced {
		envs = append(os.Environ(), export)
	}

	files := make([]*os.File, 3)
	files[syscall.Stdin] = nil
	files[syscall.Stdout] = os.Stdout
	files[syscall.Stderr] = os.Stderr

	// start child with parent process id in environment variable
	return os.StartProcess(ex, os.Args, &os.ProcAttr{
		Dir: workdir,
		Env: envs,
		Files: files,
		Sys: &syscall.SysProcAttr{
			Foreground: false,
			Setsid:     true,
		},
	})
}

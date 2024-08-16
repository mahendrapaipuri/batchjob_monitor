// Package osexec implements subprocess execution functions
package osexec

import (
	"context"
	"math"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

const (
	sudoCmd = "sudo"
)

// Execute command and return stdout/stderr.
func Execute(cmd string, args []string, env []string, logger log.Logger) ([]byte, error) {
	level.Debug(logger).Log("msg", "Executing", "command", cmd, "args", strings.Join(args, " "))

	execCmd := exec.Command(cmd, args...)

	// If env is not nil pointer, add env vars into subprocess cmd
	if env != nil {
		execCmd.Env = append(os.Environ(), env...)
	}

	// According to setpgid docs (https://man7.org/linux/man-pages/man2/setpgid.2.html)
	// we cannot use setpgid and setsid at the same time
	if cmd == sudoCmd {
		// Attach a separate terminal less session to the subprocess
		// This is to avoid prompting for password when we run command with sudo
		// Ref: https://stackoverflow.com/questions/13432947/exec-external-program-script-and-detect-if-it-requests-user-input
		execCmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	} else {
		// Start child process in its own process group so that interrupt signal will
		// not stop the command
		execCmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	}

	// Execute command
	out, err := execCmd.CombinedOutput()
	if err != nil {
		level.Error(logger).
			Log("msg", "Error executing command", "command", cmd, "args", strings.Join(args, " "), "err", err)
	}

	return out, err
}

// ExecuteAs executes a command as a given UID and GID and return stdout/stderr.
func ExecuteAs(cmd string, args []string, uid int, gid int, env []string, logger log.Logger) ([]byte, error) {
	level.Debug(logger).
		Log("msg", "Executing as user", "command", cmd, "args", strings.Join(args, " "), "uid", uid, "gid", gid)

	execCmd := exec.Command(cmd, args...)

	// Check bounds on uid and gid before converting into int32
	var uidInt32, gidInt32 uint32
	if uid > 0 && uid <= math.MaxInt32 {
		uidInt32 = uint32(uid)
	}

	if gid > 0 && gid <= math.MaxInt32 {
		gidInt32 = uint32(gid)
	}

	// According to setpgid docs (https://man7.org/linux/man-pages/man2/setpgid.2.html)
	// we cannot use setpgid and setsid at the same time
	if cmd == sudoCmd {
		// Attach a separate terminal less session to the subprocess
		// This is to avoid prompting for password when we run command with sudo
		// Ref: https://stackoverflow.com/questions/13432947/exec-external-program-script-and-detect-if-it-requests-user-input
		execCmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	} else {
		// Start child process in its own process group so that interrupt signal will
		// not stop the command
		execCmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	}

	// Set uid and gid for process
	execCmd.SysProcAttr.Credential = &syscall.Credential{Uid: uidInt32, Gid: gidInt32}

	// If env is not nil pointer, add env vars into subprocess cmd
	if env != nil {
		execCmd.Env = append(os.Environ(), env...)
	}

	// Execute command
	out, err := execCmd.CombinedOutput()
	if err != nil {
		level.Error(logger).
			Log("msg", "Error executing command as user", "command", cmd, "args", strings.Join(args, " "), "uid", uid, "gid", gid, "err", err)
	}

	return out, err
}

// ExecuteContext executes a command with context and return stdout/stderr.
func ExecuteContext(ctx context.Context, cmd string, args []string, env []string, logger log.Logger) ([]byte, error) {
	level.Debug(logger).Log("msg", "Executing with context", "command", cmd, "args", strings.Join(args, " "))

	execCmd := exec.CommandContext(ctx, cmd, args...)

	// If env is not nil pointer, add env vars into subprocess cmd
	if env != nil {
		execCmd.Env = append(os.Environ(), env...)
	}

	// According to setpgid docs (https://man7.org/linux/man-pages/man2/setpgid.2.html)
	// we cannot use setpgid and setsid at the same time
	if cmd == sudoCmd {
		// Attach a separate terminal less session to the subprocess
		// This is to avoid prompting for password when we run command with sudo
		// Ref: https://stackoverflow.com/questions/13432947/exec-external-program-script-and-detect-if-it-requests-user-input
		execCmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	} else {
		// Start child process in its own process group so that interrupt signal will
		// not stop the command
		execCmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	}

	// Execute command
	out, err := execCmd.CombinedOutput()
	if err != nil {
		level.Error(logger).
			Log("msg", "Error executing command", "command", cmd, "args", strings.Join(args, " "), "err", err)
	}

	return out, err
}

// ExecuteAsContext executes a command as a given UID and GID with context and return stdout/stderr.
func ExecuteAsContext(ctx context.Context, cmd string, args []string, uid int, gid int, env []string, logger log.Logger) ([]byte, error) {
	level.Debug(logger).
		Log("msg", "Executing as user with context", "command", cmd, "args", strings.Join(args, " "), "uid", uid, "gid", gid)

	execCmd := exec.CommandContext(ctx, cmd, args...)

	// Check bounds on uid and gid before converting into int32
	var uidInt32, gidInt32 uint32
	if uid > 0 && uid <= math.MaxInt32 {
		uidInt32 = uint32(uid)
	}

	if gid > 0 && gid <= math.MaxInt32 {
		gidInt32 = uint32(gid)
	}

	// According to setpgid docs (https://man7.org/linux/man-pages/man2/setpgid.2.html)
	// we cannot use setpgid and setsid at the same time
	if cmd == sudoCmd {
		// Attach a separate terminal less session to the subprocess
		// This is to avoid prompting for password when we run command with sudo
		// Ref: https://stackoverflow.com/questions/13432947/exec-external-program-script-and-detect-if-it-requests-user-input
		execCmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	} else {
		// Start child process in its own process group so that interrupt signal will
		// not stop the command
		execCmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	}

	// Set uid and gid for process
	execCmd.SysProcAttr.Credential = &syscall.Credential{Uid: uidInt32, Gid: gidInt32}

	// If env is not nil pointer, add env vars into subprocess cmd
	if env != nil {
		execCmd.Env = append(os.Environ(), env...)
	}

	// Execute command
	out, err := execCmd.CombinedOutput()
	if err != nil {
		level.Error(logger).
			Log("msg", "Error executing command as user", "command", cmd, "args", strings.Join(args, " "), "uid", uid, "gid", gid, "err", err)
	}

	return out, err
}

// ExecuteWithTimeout exwecutes a command with timeout and return stdout/stderr.
func ExecuteWithTimeout(cmd string, args []string, timeout int, env []string, logger log.Logger) ([]byte, error) {
	level.Debug(logger).
		Log("msg", "Executing with timeout", "command", cmd, "args", strings.Join(args, " "), "timeout", timeout)

	ctx := context.Background()

	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
		defer cancel()
	}

	execCmd := exec.CommandContext(ctx, cmd, args...)

	// If env is not nil pointer, add env vars into subprocess cmd
	if env != nil {
		execCmd.Env = append(os.Environ(), env...)
	}

	// According to setpgid docs (https://man7.org/linux/man-pages/man2/setpgid.2.html)
	// we cannot use setpgid and setsid at the same time
	if cmd == sudoCmd {
		// Attach a separate terminal less session to the subprocess
		// This is to avoid prompting for password when we run command with sudo
		// Ref: https://stackoverflow.com/questions/13432947/exec-external-program-script-and-detect-if-it-requests-user-input
		execCmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	} else {
		// Start child process in its own process group so that interrupt signal will
		// not stop the command
		execCmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	}

	// The signal to send to the children when parent receives a kill signal
	// execCmd.SysProcAttr = &syscall.SysProcAttr{Pdeathsig: syscall.SIGTERM}

	// Execute command
	out, err := execCmd.CombinedOutput()
	if err != nil {
		level.Error(logger).
			Log("msg", "Error executing command", "command", cmd, "args", strings.Join(args, " "), "err", err)
	}

	return out, err
}

// ExecuteAsWithTimeout executes a command with timeout as a given UID and GID and return stdout/stderr.
func ExecuteAsWithTimeout(
	cmd string,
	args []string,
	uid int,
	gid int,
	timeout int,
	env []string,
	logger log.Logger,
) ([]byte, error) {
	level.Debug(logger).
		Log("msg", "Executing with timeout as user", "command", cmd, "args", strings.Join(args, " "), "uid", uid, "gid", gid, "timout", timeout)

	ctx := context.Background()

	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
		defer cancel()
	}

	execCmd := exec.CommandContext(ctx, cmd, args...)

	// If env is not nil pointer, add env vars into subprocess cmd
	if env != nil {
		execCmd.Env = append(os.Environ(), env...)
	}

	// Check bounds on uid and gid before converting into int32
	var uidInt32, gidInt32 uint32
	if uid > 0 && uid <= math.MaxInt32 {
		uidInt32 = uint32(uid)
	}

	if gid > 0 && gid <= math.MaxInt32 {
		gidInt32 = uint32(gid)
	}

	// According to setpgid docs (https://man7.org/linux/man-pages/man2/setpgid.2.html)
	// we cannot use setpgid and setsid at the same time
	if cmd == sudoCmd {
		// Attach a separate terminal less session to the subprocess
		// This is to avoid prompting for password when we run command with sudo
		// Ref: https://stackoverflow.com/questions/13432947/exec-external-program-script-and-detect-if-it-requests-user-input
		execCmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	} else {
		// Start child process in its own process group so that interrupt signal will
		// not stop the command
		execCmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	}

	// Set uid and gid for process
	execCmd.SysProcAttr.Credential = &syscall.Credential{Uid: uidInt32, Gid: gidInt32}

	// Execute command
	out, err := execCmd.CombinedOutput()
	if err != nil {
		level.Error(logger).
			Log("msg", "Error executing command as user", "command", cmd, "args", strings.Join(args, " "), "uid", uid, "gid", gid, "err", err)
	}

	return out, err
}

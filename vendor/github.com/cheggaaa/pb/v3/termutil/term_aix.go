//go:build aix && !appengine
// +build aix,!appengine

package termutil

import (
	"fmt"
	"os"
	"syscall"

	"golang.org/x/sys/unix"
)

var (
	tty           = os.Stdin
	unlockSignals = []os.Signal{
		os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGSTOP, syscall.SIGINT, syscall.SIGHUP,
	}
)

// TerminalWidth returns width of the terminal.
func TerminalWidth() (int, error) {
	_, c, err := TerminalSize()
	return c, err
}

// TerminalSize returns size of the terminal.
func TerminalSize() (rows, cols int, err error) {
	ws, err := unix.IoctlGetWinsize(int(os.Stdout.Fd()), unix.TIOCGWINSZ)
	if err != nil {
		return 0, 0, err
	}
	return int(ws.Col), int(ws.Row), nil
}

func lockEcho() error {
	fd := tty.Fd()

	termios, err := unix.IoctlGetTermios(int(fd), unix.TCGETA)
	if err != nil {
		return fmt.Errorf("failed to get terminal attributes: %w", err)
	}

	termios.Lflag &^= unix.ECHO // Turn off the ECHO bit to disable echoing.

	err = unix.IoctlSetTermios(int(fd), unix.TCSETA, termios)
	if err != nil {
		return fmt.Errorf("failed to set terminal attributes: %w", err)
	}

	return nil
}

func unlockEcho() error {
	fd := tty.Fd()
	termios, err := unix.IoctlGetTermios(int(fd), unix.TCGETA)
	if err != nil {
		return fmt.Errorf("failed to get terminal attributes: %w", err)
	}

	termios.Lflag |= unix.ECHO // Turn on the ECHO bit to enable echoing.

	err = unix.IoctlSetTermios(int(fd), unix.TCSETA, termios)
	if err != nil {
		return fmt.Errorf("failed to set terminal attributes: %w", err)
	}

	return nil
}

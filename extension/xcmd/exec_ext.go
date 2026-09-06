// Package xcmd
// @author: fengyi
// @date: 2024/5/25
// @note:
package xcmd

import (
	"fmt"
	"general-agent/extension/logz"
	"os/exec"
)

func Exec(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	logz.InfoNoCtx(fmt.Sprintf("exec cmd: %v", cmd.Args))
	return cmd
}

func NamespaceExec(name string, args ...string) *exec.Cmd {
	a := []string{"netns", "exec", "business", name}
	args = append(a, args...)
	return Exec("ip", args...)
}

func NamespaceExecIf(b bool, name string, args ...string) *exec.Cmd {
	if b {
		return NamespaceExec(name, args...)
	} else {
		return Exec(name, args...)
	}
}

func ExecOutput(name string, args ...string) ([]byte, error) {
	output, err := Exec(name, args...).Output()
	logz.InfoNoCtx(fmt.Sprintf("exec out: %s", string(output)))
	return output, err
}

func NamespaceExecOutput(name string, args ...string) ([]byte, error) {
	output, err := NamespaceExec(name, args...).Output()
	logz.InfoNoCtx(fmt.Sprintf("exec out: %s", string(output)))
	return output, err
}

func NamespaceExecOutputIf(b bool, name string, args ...string) ([]byte, error) {
	if b {
		return NamespaceExecOutput(name, args...)
	} else {
		return ExecOutput(name, args...)
	}
}

func ExecCombinedOutput(name string, args ...string) ([]byte, error) {
	output, err := Exec(name, args...).CombinedOutput()
	logz.InfoNoCtx(fmt.Sprintf("exec out: %s", string(output)))
	return output, err
}
func NamespaceExecCombinedOutput(name string, args ...string) ([]byte, error) {
	output, err := NamespaceExec(name, args...).CombinedOutput()
	logz.InfoNoCtx(fmt.Sprintf("exec out: %s", string(output)))
	return output, err
}

func NamespaceExecCombinedOutputIf(b bool, name string, args ...string) ([]byte, error) {
	if b {
		return NamespaceExecCombinedOutput(name, args...)
	} else {
		return ExecCombinedOutput(name, args...)
	}
}

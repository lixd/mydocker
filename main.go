package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

// 注: 运行时需要 root 权限。
func main() {
	namespace()
}

// namespace 如何在Go中新建Namespace
func namespace() {
	cmd := exec.Command("bash")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
			syscall.CLONE_NEWUSER | syscall.CLONE_NEWNET,
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalln(err)
	}
}

// cGroups cGroups初体验
func cGroups() {

}

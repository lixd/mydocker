package main

import (
	"fmt"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

func commitContainer(imageName string) {
	mntPath := "/root/merged"
	imageTar := "/root/" + imageName + ".tar"
	fmt.Println("commitContainer imageTar:", imageTar)
	if _, err := exec.Command("tar", "-czf", imageTar, "-C", mntPath, ".").CombinedOutput(); err != nil {
		log.Errorf("tar folder %s error %v", mntPath, err)
	}
}

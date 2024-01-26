package main

import (
	"fmt"
	"io"
	"os"

	"mydocker/container"

	log "github.com/sirupsen/logrus"
)

func logContainer(containerId string) {
	logFileLocation := fmt.Sprintf(container.InfoLocFormat, containerId) + container.GetLogfile(containerId)
	file, err := os.Open(logFileLocation)
	defer file.Close()
	if err != nil {
		log.Errorf("Log container open file %s error %v", logFileLocation, err)
		return
	}
	content, err := io.ReadAll(file)
	if err != nil {
		log.Errorf("Log container read file %s error %v", logFileLocation, err)
		return
	}
	_, err = fmt.Fprint(os.Stdout, string(content))
	if err != nil {
		log.Errorf("Log container Fprint  error %v", err)
		return
	}
}

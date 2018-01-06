package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/vulcand/vulcand/vctl/command"
	"github.com/yieldr/vulcand/registry"
)

func main() {
	r, err := registry.GetRegistry()
	if err != nil {
		log.Errorf("Error: %s\n", err)
		return
	}

	cmd := command.NewCommand(r)
	if err := cmd.Run(os.Args); err != nil {
		log.Errorf("Error: %s\n", err)
	}
}

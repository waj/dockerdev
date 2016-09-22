package main

import (
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
	"log"
)

func main() {
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Fatal(err)
	}

	eventReader, err := cli.Events(context.Background(), types.EventsOptions{})
	if err != nil {
		log.Fatal(err)
	}

	dec := json.NewDecoder(eventReader)
	for {
		var event events.Message
		dec.Decode(&event)

		if event.Type == events.ContainerEventType && event.Action == "create" {
			project, hasProject := event.Actor.Attributes["com.docker.compose.project"]
			service, hasService := event.Actor.Attributes["com.docker.compose.service"]
			oneoff := event.Actor.Attributes["com.docker.compose.oneoff"]

			if hasProject && hasService {
				config := network.EndpointSettings{}
				if oneoff == "False" {
					config.Aliases = []string{service + "." + project + ".dev"}
				}

				err := cli.NetworkConnect(context.Background(), "shared", event.Actor.ID, &config)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}

}

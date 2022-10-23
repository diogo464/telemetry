package main

import (
	"fmt"

	"github.com/diogo464/telemetry"
	"github.com/urfave/cli/v2"
)

var CommandStreams = &cli.Command{
	Name:        "streams",
	Description: "Show available streams",
	Action:      actionStreams,
}

func actionStreams(c *cli.Context) error {
	client, err := clientFromContext(c)
	if err != nil {
		return err
	}
	defer client.Close()

	streams, err := client.AvailableStreams(c.Context)
	if err != nil {
		return err
	}

	for _, stream := range streams {
		fmt.Println(stream.Name)
		fmt.Println("\tEncoding:", telemetry.ReadableEncoding(stream.Encoding))
	}

	return nil
}

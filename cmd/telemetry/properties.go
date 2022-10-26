package main

import (
	"fmt"

	"github.com/diogo464/telemetry"
	"github.com/urfave/cli/v2"
)

var CommandProperties = &cli.Command{
	Name:        "properties",
	Description: "Show available properties",
	Action:      actionProperties,
}

func actionProperties(c *cli.Context) error {
	client, err := clientFromContext(c)
	if err != nil {
		return err
	}
	defer client.Close()

	properties, err := client.GetAvailableProperties(c.Context)
	if err != nil {
		return err
	}

	for _, prop := range properties {
		fmt.Println(prop.Name)
		fmt.Println("\tEncoding:", telemetry.ReadableEncoding(prop.Encoding))
	}

	return nil
}

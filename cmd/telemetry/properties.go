package main

import (
	"fmt"

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

	properties, err := client.GetProperties(c.Context)
	if err != nil {
		return err
	}

	for _, prop := range properties {
		fmt.Println(prop.Name, "=", prop.Value)
		fmt.Println("\t", prop.Description)
	}

	return nil
}

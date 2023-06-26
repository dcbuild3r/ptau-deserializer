package main

import (
	"deserializer"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "ptau-deserialize",
		Usage: "Deserialize .ptau files into gnark's .ph1 format",
		Action: func(*cli.Context) error {
			fmt.Println("Hello friend!")
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "convert",
				Aliases: []string{"a"},
				Usage:   "Deserialize .ptau file into .ph1 format",
				Action: func(cCtx *cli.Context) error {
					fmt.Println(".ptau file path ", cCtx.Args().First())
					ptauFilePath := cCtx.Args().First()

					ptau, err := deserializer.ReadPtau(ptauFilePath)

					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "ptau-file",
						Aliases: []string{"c"},
						Usage:   "Load .ptau `FILE`",
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("%s", err)
	}
}

package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "zkey-deserialize",
		Usage: "Desereialize .zkey files into gnark's .ph1 format",
		Action: func(*cli.Context) error {
			fmt.Println("Hello friend!")
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "input",
				Aliases: []string{"c"},
				Usage:   "Load .zkey `FILE`",
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("%s", err)
	}
}

package main

import (
	"fmt"
	"os"

	"github.com/bnb-chain/zkbnb-setup/phase2"
	deserializer "github.com/dcbuild3r/ptau-deserializer/deserialize"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "ptau-deserialize",
		Usage: "Deserialize .ptau files into gnark's .ph1 format",
		Action: func(*cli.Context) error {
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "convert",
				Aliases: []string{"c"},
				Usage:   "Deserialize .ptau file into .ph1 format and write to `OUTPUT`",
				Action: func(cCtx *cli.Context) error {
					ptauFilePath := cCtx.String("input")
					outputFilePath := cCtx.String("output")

					ptau, err := deserializer.ReadPtau(ptauFilePath)

					if err != nil {
						return err
					}

					phase1, err := deserializer.ConvertPtauToPhase1(ptau)

					if err != nil {
						return err
					}

					// Write phase1 to file
					err = deserializer.WritePhase1(phase1, uint8(ptau.Header.Power), outputFilePath)

					if err != nil {
						return err
					}

					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "input",
						Aliases:  []string{"i"},
						Usage:    "Load `FILE`.ptau to convert to .ph1",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "output",
						Aliases:  []string{"o"},
						Usage:    "File output for the phase 1 conversion (`FILE`.ph1)",
						Required: true,
					},
				},
			},
			{
				Name:    "initialize",
				Aliases: []string{"i"},
				Usage:   "Initialize phase 2 from phase 1 and r1cs files. Output to `OUTPUT`.ph2",
				Action: func(cCtx *cli.Context) error {
					phase1FilePath := cCtx.String("input")
					r1csFilePath := cCtx.String("r1cs")
					outputFilePath := cCtx.String("output")

					err := phase2.Initialize(phase1FilePath, r1csFilePath, outputFilePath)

					if err != nil {
						return err
					}

					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "input",
						Aliases:  []string{"i"},
						Usage:    "Load `FILE`.ph1 to initialize phase 2",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "r1cs",
						Aliases:  []string{"r"},
						Usage:    "Load `FILE`.r1cs to initialize phase 2",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "output",
						Aliases:  []string{"o"},
						Usage:    "File output for the phase 2 initialization (`FILE`.ph2)",
						Required: true,
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("%s", err)
	}
}

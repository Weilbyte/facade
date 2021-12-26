package main

import (
	"fmt"
	"log"
	"os"

	cli "github.com/urfave/cli/v2"
	"github.com/weilbyte/facade/lib"
)

func main() {
	app := &cli.App{
		Name:  "facade",
		Usage: "Generates a DLL proxy project",
		Commands: []*cli.Command{
			{
				Name:  "generate",
				Usage: "Generates the DLL proxy project",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "embed",
						Value: false,
						Usage: "Embed the target .dll into the generated project",
					},
					&cli.StringFlag{
						Name:  "out",
						Value: "./out",
						Usage: "Output directory",
					},
				},
				Action: func(c *cli.Context) error {
					if c.NArg() != 1 {
						return cli.Exit("Expecting path to target .dll as argument", 1)
					}

					targetDll := c.Args().First()

					pe, err := lib.GetAndValidate(targetDll)
					if err != nil {
						return cli.Exit(err.Error(), 1)
					}

					lib.GenerateProject(pe, targetDll, c.Bool("embed"), c.String("out"))

					fmt.Printf("Successfully generated the proxy DLL project in %s, with %d forwarded export(s)\n\n", c.String("out"), len(pe.Export.Functions))
					fmt.Println("Edit the main.cpp file to add your payload, then build the project with CMake using a Visual Studio generator.")
					if c.Bool("embed") {
						fmt.Println("The built DLL contains the original and will link with it at runtime.")
					} else {
						fmt.Println("The built DLL will need to have the original in the same folder with \"_o\" appended to its name.")
					}
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

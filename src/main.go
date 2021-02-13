package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/fatih/color"
	. "github.com/nogen-app/create-mf-app/src/flagvalues"
	"github.com/nogen-app/create-mf-app/src/io"
	cli "github.com/urfave/cli/v2"
)

const (
	gitTemplate = "https://github.com/nogen-app/create-mf-app-ts-react-template"
	version     = "v0.1.0"
)

func exitWithError(e *error) {
	color.Set(color.FgRed)
	fmt.Println(*e)
	color.Unset()
	os.Exit(1)
}

func readStdoutLines(scanner *bufio.Scanner, outputChan chan string, finishedChan chan bool) {
	for scanner.Scan() {
		outputChan <- scanner.Text()
	}
	finishedChan <- true
}

func npmClientInstall(npmClient string, finishedChan chan bool, outputChan chan string) {
	cmdNpmClientInstall := exec.Command(npmClient, "install")

	cmdReader, err := cmdNpmClientInstall.StdoutPipe()
	if err != nil {
		exitWithError(&err)
	}

	scriptOutputChan := make(chan string)
	scriptfinishecChan := make(chan bool)

	scanner := bufio.NewScanner(cmdReader)
	go readStdoutLines(scanner, scriptOutputChan, scriptfinishecChan)

	if err := cmdNpmClientInstall.Start(); err != nil {
		exitWithError(&err)
	}
	scriptIsRunning := true

	for scriptIsRunning {
		select {
		case output := <-scriptOutputChan:
			outputChan <- output
		case status := <-scriptfinishecChan:
			scriptIsRunning = !status
		default:
		}
	}

	if err := cmdNpmClientInstall.Wait(); err != nil {
		exitWithError(&err)
	}

	finishedChan <- true
}

func main() {
	app := &cli.App{
		Name:     "create-mf-app",
		Usage:    "Used to scaffold module federation app with typescript and react",
		Version:  version,
		Compiled: time.Now(),
		Authors: []*cli.Author{
			{
				Name:  "nogen I/S",
				Email: "contact@nogen.app",
			},
		},
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Value:   3000,
				Usage:   "The port to use",
			},
			&cli.GenericFlag{
				Name:    "npmclient",
				Aliases: []string{"c"},
				Value: &EnumValue{
					Enum:    []string{"yarn", "npm", "npmm"},
					Default: "yarn",
				},
			},
			&cli.StringFlag{
				Name:        "template",
				Aliases:     []string{"t"},
				Value:       "",
				Usage:       "The template version to use",
				DefaultText: "latest version",
			},
		},
		Action: func(c *cli.Context) error {
			var projectName string
			if c.NArg() == 0 {
				cli.ShowAppHelp(c)
				return fmt.Errorf("projectName not specified")
			}

			projectName = c.Args().First()
			port := c.Int("port")
			npmClient := c.String("npmclient")
			templateVersion := c.String("template")

			if c.NArg() > 1 {
				cli.ShowAppHelp(c)
				return fmt.Errorf("%s %s specified as args, should be option like create-mf-app %s %s %s", c.Args().Get(1), c.Args().Get(2), c.Args().Get(1), c.Args().Get(2), projectName)
			}

			var requiredPrograms []string
			requiredPrograms = append(requiredPrograms, "git")
			requiredPrograms = append(requiredPrograms, npmClient)

			if err := io.CheckRequiredPrograms(requiredPrograms); err != nil {
				exitWithError(&err)
			}

			if _, err := os.Stat("./" + projectName); !os.IsNotExist(err) {
				err = fmt.Errorf("Folder %s already exists", projectName)
				exitWithError(&err)
			}

			if err := io.CloneAndCdGitRepo(gitTemplate, projectName); err != nil {
				exitWithError(&err)
			}

			if err := io.GetTemplateVersion(templateVersion, version); err != nil {
				exitWithError(&err)
			}

			if err := io.CleanGitRepo(); err != nil {
				exitWithError(&err)
			}

			var placeHolders []io.Placeholder
			placeHolders = append(placeHolders, io.Placeholder{Key: "NAME_PLACEHOLDER", Value: projectName})
			placeHolders = append(placeHolders, io.Placeholder{Key: "PORT_PLACEHOLDER", Value: fmt.Sprint(port)})

			placeholderReplacer := io.PlaceholderReplacer{
				Placeholders: placeHolders,
			}

			if err := placeholderReplacer.ReplacePlaceholders("."); err != nil {
				exitWithError(&err)
			}

			npmClientFinished := make(chan bool)
			npmClientOutput := make(chan string)
			go npmClientInstall(npmClient, npmClientFinished, npmClientOutput)

			idx := 0
			var loadingSymbol string
			npmClientIsNotRunning := false
			for !npmClientIsNotRunning {
				switch idx {
				case 0:
					loadingSymbol = "/"
				case 1:
					loadingSymbol = "-"
				case 2:
					loadingSymbol = "\\"
				case 3:
					loadingSymbol = "|"
				}
				fmt.Print(fmt.Sprintf("\r %s running", loadingSymbol))
				if idx != 3 {
					idx++
				} else {
					idx = 0
				}

				select {
				case npmClientStatus := <-npmClientFinished:
					npmClientIsNotRunning = npmClientStatus
				case logLine := <-npmClientOutput:
					fmt.Printf("\r\n %s \033[F", logLine)
				default:
				}
			}

			fmt.Print("\r")
			io.PrintGreenCheckmark("Finished                \r\n")
			io.PrintGreenCheckmark(fmt.Sprintf("project: %s was succesfully created with %s serving on port %d", projectName, npmClient, port))

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		color.Set(color.FgRed)
		exitWithError(&err)
		color.Unset()
	}
}

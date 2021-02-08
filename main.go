package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
	cli "github.com/urfave/cli/v2"
)

type symbol string

const (
	gitTemplate = "https://github.com/nogen-app/create-mf-app-ts-react-template"
	version     = "v0.0.1"
)

const (
	checkmark symbol = "✓"
	cross     symbol = "˟"
)

type placeHolder struct {
	Key   string
	Value string
}

var placeHolders []placeHolder

type enumValue struct {
	Enum     []string
	Default  string
	Selected string
}

func (e *enumValue) Set(value string) error {
	for _, enum := range e.Enum {
		if enum == value {
			e.Selected = value
			return nil
		}
	}

	return fmt.Errorf("Allowed values are %s", strings.Join(e.Enum, ", "))
}

func (e enumValue) String() string {
	if e.Selected == "" {
		return e.Default
	}
	return e.Selected
}

func exitWithError(e *error) {
	color.Set(color.FgRed)
	fmt.Println(*e)
	color.Unset()
	os.Exit(1)
}

func checkRequiredPrograms(npmClient string) error {
	var errorHasOccuredFlag bool
	var programError error

	if err := checkProgram("git"); err != nil {
		errorHasOccuredFlag = true
		programError = err
	}

	if err := checkProgram(npmClient); err != nil {
		errorHasOccuredFlag = true
		programError = err
	}

	if errorHasOccuredFlag {
		return programError
	}

	return nil
}

func checkProgram(program string) error {
	cmdProgramVersion := exec.Command(program, "--version")

	programVersion, err := cmdProgramVersion.Output()
	if err != nil {
		red := color.New(color.FgRed).PrintfFunc()
		red("[%s] ", cross)
		fmt.Printf("%s is not installed or could not be found in path", program)
		fmt.Println()
		return err
	}

	green := color.New(color.FgGreen).PrintfFunc()

	green("[%s] ", checkmark)
	fmt.Printf("%s is installed with version: %s", program, programVersion)

	return nil
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func writeLines(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}

func replacePlaceholder(path string, info os.FileInfo, err error) error {
	if info.IsDir() {
		for _, placeholder := range placeHolders {
			if info.Name() == placeholder.Key {
				newPath := strings.ReplaceAll(path, info.Name(), placeholder.Value)
				os.Rename(path, newPath)
			}
		}

	} else {
		lines, err := readLines(path)
		if err != nil {
			return err
		}

		var newLines []string

		for _, line := range lines {
			newLine := line
			for _, placeholder := range placeHolders {
				newLine = strings.ReplaceAll(newLine, placeholder.Key, placeholder.Value)
			}
			newLines = append(newLines, newLine)
		}

		if err := writeLines(newLines, path); err != nil {
			return err
		}
	}

	return nil
}

func replacePlaceholders() error {
	if err := filepath.Walk(".", replacePlaceholder); err != nil {
		return err
	}

	return nil
}

func getTemplateVersion(userSelectedVersion string) error {
	cmdGitFetch := exec.Command("git", "fetch", "--all", "--tags", "--prune")
	if err := cmdGitFetch.Run(); err != nil {
		return err
	}

	selectedVersion := userSelectedVersion

	if selectedVersion == "" {
		// Get latest compatible version
		cmdGitTags := exec.Command("git", "tag")
		gitTags, err := cmdGitTags.Output()
		if err != nil {
			return err
		} else {
			tags := strings.Split(string(gitTags), "\n")

			var acceptableTags []string

			for _, tag := range tags {
				majorTemplateVersion := strings.Split(tag, ".")[0]
				majorScriptVersion := strings.Split(version, ".")[0]

				if majorTemplateVersion == majorScriptVersion {
					acceptableTags = append(acceptableTags, tag)
				}
			}

			sort.Strings(acceptableTags)
			selectedVersion = acceptableTags[0]
		}
	}

	// Get specific version
	cmdGitCheckout := exec.Command("git", "checkout", "tags/"+selectedVersion)
	if err := cmdGitCheckout.Run(); err != nil {
		return err
	}

	return nil
}

func cleanGitRepo() error {
	if err := os.RemoveAll(".git"); err != nil {
		return err
	}

	if err := os.RemoveAll(".github"); err != nil {
		return err
	}

	return nil
}

func readStdoutLines(scanner *bufio.Scanner, outputChan chan string, finishedChan chan bool) {
	for scanner.Scan() {
		outputChan <- scanner.Text()
	}
	finishedChan <- true
}

func npmClientInstall(npmClient string, finishedChan chan bool, outputChan chan string) {
	cmdNpmClientInstall := exec.Command(npmClient, "install")

	// out, err := cmdNpmClientInstall.Output()
	// if err != nil {
	// 	fmt.Print(string(out))
	// 	exitWithError(&err)
	// }

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
				Value: &enumValue{
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

			if err := checkRequiredPrograms(npmClient); err != nil {
				exitWithError(&err)
			}

			placeHolders = append(placeHolders, placeHolder{Key: "NAME_PLACEHOLDER", Value: projectName})
			placeHolders = append(placeHolders, placeHolder{Key: "PORT_PLACEHOLDER", Value: fmt.Sprint(port)})

			if _, err := os.Stat("./" + projectName); !os.IsNotExist(err) {
				err = fmt.Errorf("Folder %s already exists", projectName)
				exitWithError(&err)
			}

			// clone git repo to specific folder
			cmdGitClone := exec.Command("git", "clone", gitTemplate, projectName)
			cmdGitClone.Run()

			// // go to git repo
			if err := os.Chdir("./" + projectName); err != nil {
				exitWithError(&err)
			}

			if err := getTemplateVersion(templateVersion); err != nil {
				exitWithError(&err)
			}

			if err := cleanGitRepo(); err != nil {
				exitWithError(&err)
			}

			if err := replacePlaceholders(); err != nil {
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
			green := color.New(color.FgGreen).PrintfFunc()

			green("\r[%s] ", checkmark)
			fmt.Print("Finished        \r\n")

			green("[%s] ", checkmark)
			fmt.Printf("project: %s was succesfully created with %s serving on port %d", projectName, npmClient, port)

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

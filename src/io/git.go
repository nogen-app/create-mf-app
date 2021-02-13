package io

import (
	"os"
	"os/exec"
	"sort"
	"strings"
)

//GetTemplateVersion checkouts either the latest compatible version or the user specified version
//Used to make sure the current versionof the cli can always use a newer compatible template version, without begin updated
func GetTemplateVersion(userSelectedVersion string, cliVersion string) error {
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
				majorScriptVersion := strings.Split(cliVersion, ".")[0]

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

//CleanGitRepo removes .git folder and .github folder in current path
//Used to remove the git refrences from the template repo
func CleanGitRepo() error {
	if err := os.RemoveAll(".git"); err != nil {
		return err
	}

	if err := os.RemoveAll(".github"); err != nil {
		return err
	}

	return nil
}

//CloneAndCdGitRepo clones a given repo to the specified project name, and cd's into it
//Used to get the selected template
func CloneAndCdGitRepo(gitTemplate string, projectName string) error {
	cmdGitClone := exec.Command("git", "clone", gitTemplate, projectName)
	cmdGitClone.Run()
	if err := os.Chdir("./" + projectName); err != nil {
		return err
	}

	return nil
}

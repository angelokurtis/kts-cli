package aws

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os/user"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"
)

func ListProfiles() ([]string, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}
	filename := usr.HomeDir + "/.aws/config"
	out, err := ioutil.ReadFile(filename)
	scanner := bufio.NewScanner(bytes.NewReader(out))
	profiles := make([]string, 0, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			profile := line[1:]
			profile = profile[:len(profile)-1]
			profile = strings.Replace(profile, "profile ", "", 1)
			profiles = append(profiles, profile)
		}
	}
	return profiles, nil
}

func SelectProfiles() ([]string, error) {
	profiles, err := ListProfiles()
	if err != nil {
		return nil, err
	}

	var selects []string
	prompt := &survey.MultiSelect{
		Message: "Select the AWS profiles:",
		Options: profiles,
	}

	err = survey.AskOne(prompt, &selects, survey.WithPageSize(10), survey.WithKeepFilter(true))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return selects, nil
}

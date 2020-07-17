package sensedia

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func SelectGateway() (string, error) {
	gateways := []string{"https://manager-testing.sensedia-eng.com/api-manager"}

	if len(gateways) == 1 {
		return gateways[0], nil
	}

	var selected string
	prompt := &survey.Select{
		Message: "Select the Gateway address:",
		Options: gateways,
	}

	err := survey.AskOne(prompt, &selected, survey.WithPageSize(10))
	if err != nil {
		return "", errors.WithStack(err)
	}
	return selected, nil
}

func SelectGatewayUser() (string, string, error) {
	users := map[string]string{
		"deyverson": "Deyverson@16",
		"root":      "manager",
	}

	names := make([]string, 0, 0)
	for name, _ := range users {
		names = append(names, name)
	}

	if len(names) == 1 {
		return names[0], users[names[0]], nil
	}

	var selected string
	prompt := &survey.Select{
		Message: "Witch user would you like to use?",
		Options: names,
	}

	err := survey.AskOne(prompt, &selected, survey.WithPageSize(10))
	if err != nil {
		return "", "", errors.WithStack(err)
	}
	return selected, users[selected], nil
}

func Login(gateway string, login string, password string) (string, string, error) {

	data := url.Values{}
	data.Set("login", login)
	data.Set("password", password)

	gateway = strings.TrimSuffix(gateway, "/") + "/"
	req, _ := http.NewRequest("POST", gateway+"api/v3/login", strings.NewReader(data.Encode())) // URL-encoded payload
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", errors.WithStack(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		return "", "", errors.New("gateway answered with code " + resp.Status)
	}

	sc := resp.Header["Set-Cookie"][0]
	sensediaAuth := ""
	for _, cookie := range strings.Split(sc, ";") {
		if strings.Contains(cookie, "Sensedia-Auth=") {
			sensediaAuth = strings.Split(strings.TrimSpace(cookie), "=")[1]
		}
	}
	xsrf := resp.Header["Xsrf-Token"][0]

	return sensediaAuth, xsrf, nil
}

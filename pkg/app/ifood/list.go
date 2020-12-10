package ifood

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
	"os"
)

var client = &http.Client{}

const url = "https://marketplace.ifood.com.br/v4/customers/me/orders?page=0&size=10000"

func List() (Orders, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	accountId := os.Getenv("IFOOD_ACCOUNT_ID")
	token := os.Getenv("IFOOD_TOKEN")
	req.Header.Add("account_id", accountId)
	req.Header.Add("authorization", "Bearer "+token)

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer res.Body.Close()

	var target []*Order
	if err := json.NewDecoder(res.Body).Decode(&target); err != nil {
		return nil, errors.WithStack(err)
	}

	return target, nil
}

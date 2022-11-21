package lorenz

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/SimonRichardson/colossus/blist/selectors"
	cli "github.com/SimonRichardson/colossus/blist/services/lorenz/client"
)

type versionBody struct {
	Records map[string][]string `json:"records"`
}

func version(client cli.Client) (selectors.Version, error) {
	response, err := client.Get("version", cli.Unversioned, func(http.Header) {})
	if err != nil {
		return selectors.Version(""), err
	}

	if err := readError(response, http.StatusOK); err != nil {
		return selectors.Version(""), err
	}

	var record versionBody
	if err := json.Unmarshal(response.Bytes, &record); err != nil {
		return selectors.Version(""), err
	}

	version := ""
	if x, ok := record.Records["lorenz"]; ok {
		version = strings.Join(x, ",")
	}

	return selectors.Version(version), nil
}

package util

import (
	"fmt"
	"github.com/ghetzel/cli"
	"gopkg.in/unrolled/render.v1"
	"net/http"
	"time"
)

const ApplicationName = `corona-api`
const ApplicationSummary = `a REST API for building web-based graphical user session applications`
const ApplicationVersion = `0.0.5`

var StartedAt = time.Now()
var ApiRenderer = render.New()

func RegisterSubcommands() []cli.Command {
	return []cli.Command{
		{
			Name:  `version`,
			Usage: `Output only the version string and exit`,
			Action: func(c *cli.Context) {
				fmt.Println(ApplicationVersion)
			},
		},
	}
}

func Respond(w http.ResponseWriter, code int, payload interface{}, err error) {
	response := make(map[string]interface{})
	response[`responded_at`] = time.Now().Format(time.RFC3339)
	response[`payload`] = payload

	if code >= http.StatusBadRequest {
		response[`success`] = false

		if err != nil {
			response[`error`] = err.Error()
		}
	} else {
		response[`success`] = true
	}

	ApiRenderer.JSON(w, code, response)
}

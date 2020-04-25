package main

import (
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
	"strings"
	"sync"
)

type Plugin struct {
	plugin.MattermostPlugin
	configurationLock sync.RWMutex
	configuration *configuration
}

func getCommand() *model.Command {
    return &model.Command{
        Trigger: "crm",
        Description: "Look up contact info in vTiger crm.",
        DisplayName: "vTiger",
        AutoComplete: true,
        AutoCompleteDesc: "Search vTiger crm.",
        AutoCompleteHint: "[keyword]",
    }
}

func getCommandResponse(responseType, text string) *model.CommandResponse {
    return &model.CommandResponse{
		ResponseType: responseType,
		Text:         text,
		Username:     "vTiger",
		Type:         model.POST_DEFAULT,
	}
}

func (p *Plugin) OnActivate() error {
    p.API.RegisterCommand(getCommand())
    return nil
}

func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	var results string
	split := strings.Fields(args.Command)
	if (len(split) < 2) {
		results = "Keyword is missing. Usage: /crm [keyword]."
	} else {
		results = p.search(split[1])
	}
    resp := getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, results)
    return resp, nil
}

func (p *Plugin) search(keyword string) string {
	config := vTigerAccess{
		 p.configuration.UserName,
		 p.configuration.AccessKey,
		 p.configuration.BaseUrl,
	}
	return search(config, keyword)
}

package mattermost

import (
	"net/http"

	"github.com/mattermost/mattermost-app-servicenow/app"
	"github.com/mattermost/mattermost-app-servicenow/clients/mattermostclient"
	"github.com/mattermost/mattermost-app-servicenow/constants"
	"github.com/mattermost/mattermost-app-servicenow/utils"
	"github.com/mattermost/mattermost-plugin-apps/apps"
)

func fBindings(w http.ResponseWriter, r *http.Request, c *apps.Call) {
	commands := &apps.Binding{
		Location: apps.LocationCommand,
		Bindings: []*apps.Binding{},
	}

	connectionCommand := getConnectBinding()

	if app.IsUserConnected(c.Context.BotAccessToken, c.Context.MattermostSiteURL, c.Context.ActingUserID) {
		connectionCommand = getDisconnectBinding()
	}

	commands.Bindings = append(commands.Bindings, connectionCommand)
	client := mattermostclient.NewMMClient(c.Context.BotUserID, c.Context.BotAccessToken, c.Context.MattermostSiteURL)

	user, err := client.GetUser(c.Context.ActingUserID)
	if err == nil && user.IsSystemAdmin() {
		commands.Bindings = append(commands.Bindings, getSysAdminCommandBindings())
	}

	out := []*apps.Binding{}

	if app.IsUserConnected(c.Context.BotAccessToken, c.Context.MattermostSiteURL, c.Context.ActingUserID) {
		postBindings, commandBindings, headerBindings := app.GetTablesBindings()
		if postBindings != nil {
			out = append(out, &apps.Binding{
				Location: apps.LocationPostMenu,
				Bindings: []*apps.Binding{generateTableBindingsCalls(postBindings)},
			})
		}

		if commandBindings != nil {
			commands.Bindings = append(commands.Bindings, generateTableBindingsCalls(commandBindings))
		}

		if headerBindings != nil {
			out = append(out, &apps.Binding{
				Location: apps.LocationChannelHeader,
				Bindings: []*apps.Binding{generateTableBindingsCalls(headerBindings)},
			})
		}
	}

	out = append(out, commands)

	utils.WriteBindings(w, out)
}

func generateTableBindingsCalls(b *apps.Binding) *apps.Binding {
	if len(b.Bindings) == 0 {
		b.Call = getCreateTicketCall(b.Call.Path, formActionOpen)
	}

	for _, subBinding := range b.Bindings {
		subBinding.Call = getCreateTicketCall(subBinding.Call.Path, formActionOpen)
	}

	return b
}

func getSysAdminCommandBindings() *apps.Binding {
	return &apps.Binding{
		Location:    constants.LocationConfigure,
		Label:       "config",
		Icon:        "",
		Hint:        "",
		Description: "Configure the plugin",
		Bindings: []*apps.Binding{
			{
				Location:    constants.LocationConfigureOAuth,
				Label:       "oauth",
				Icon:        "",
				Hint:        "",
				Description: "Configure OAuth options",
				Call:        getConfigureOAuthCall(formActionOpen),
			},
		},
	}
}
func getConnectBinding() *apps.Binding {
	return &apps.Binding{
		Location:    constants.LocationConnect,
		Label:       "connect",
		Icon:        "",
		Hint:        "",
		Description: "Connect your ServiceNow account",
		Form:        &apps.Form{},
		Call:        getConnectCall(),
	}
}

func getDisconnectBinding() *apps.Binding {
	return &apps.Binding{
		Location:    constants.LocationDisconnect,
		Label:       "disconnect",
		Icon:        "",
		Hint:        "",
		Description: "Disconnect from ServiceNow",
		Form:        &apps.Form{},
		Call:        getDisconnectCall(),
	}
}

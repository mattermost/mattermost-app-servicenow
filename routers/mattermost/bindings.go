package mattermost

import (
	"net/http"

	"github.com/mattermost/mattermost-app-servicenow/app"
	"github.com/mattermost/mattermost-app-servicenow/clients/mattermostclient"
	"github.com/mattermost/mattermost-app-servicenow/constants"
	"github.com/mattermost/mattermost-app-servicenow/utils"
	"github.com/mattermost/mattermost-plugin-apps/server/api"
)

func fBindings(w http.ResponseWriter, r *http.Request, claims *api.JWTClaims, c *api.Call) {
	commands := &api.Binding{
		Location: api.LocationCommand,
		Bindings: []*api.Binding{},
	}

	connectionCommand := getConnectBinding()

	if app.IsUserConnected(c.Context.BotAccessToken, c.Context.MattermostSiteURL, claims.ActingUserID) {
		connectionCommand = getDisconnectBinding()
	}

	commands.Bindings = append(commands.Bindings, connectionCommand)
	client := mattermostclient.NewMMClient(c.Context.BotUserID, c.Context.BotAccessToken, c.Context.MattermostSiteURL)

	user, err := client.GetUser(claims.ActingUserID)
	if err == nil && user.IsSystemAdmin() {
		commands.Bindings = append(commands.Bindings, getSysAdminCommandBindings())
	}

	out := []*api.Binding{}

	if app.IsUserConnected(c.Context.BotAccessToken, c.Context.MattermostSiteURL, claims.ActingUserID) {
		postBindings, commandBindings, headerBindings := app.GetTablesBindings()
		if postBindings != nil {
			out = append(out, &api.Binding{
				Location: api.LocationPostMenu,
				Bindings: []*api.Binding{generateTableBindingsCalls(postBindings)},
			})
		}

		if commandBindings != nil {
			commands.Bindings = append(commands.Bindings, generateTableBindingsCalls(commandBindings))
		}

		if headerBindings != nil {
			out = append(out, &api.Binding{
				Location: api.LocationChannelHeader,
				Bindings: []*api.Binding{generateTableBindingsCalls(headerBindings)},
			})
		}
	}

	out = append(out, commands)

	utils.WriteBindings(w, out)
}

func generateTableBindingsCalls(b *api.Binding) *api.Binding {
	if len(b.Bindings) == 0 {
		b.Call = getCreateTicketCall(b.Call.URL, formActionOpen)
	}

	for _, subBinding := range b.Bindings {
		subBinding.Call = getCreateTicketCall(subBinding.Call.URL, formActionOpen)
	}

	return b
}

func getSysAdminCommandBindings() *api.Binding {
	return &api.Binding{
		Location:    constants.LocationConfigure,
		Label:       "config",
		Icon:        "",
		Hint:        "",
		Description: "Configure the plugin",
		Bindings: []*api.Binding{
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
func getConnectBinding() *api.Binding {
	return &api.Binding{
		Location:    constants.LocationConnect,
		Label:       "connect",
		Icon:        "",
		Hint:        "",
		Description: "Connect your ServiceNow account",
		Form:        &api.Form{},
		Call:        getConnectCall(),
	}
}

func getDisconnectBinding() *api.Binding {
	return &api.Binding{
		Location:    constants.LocationDisconnect,
		Label:       "disconnect",
		Icon:        "",
		Hint:        "",
		Description: "Disconnect from ServiceNow",
		Form:        &api.Form{},
		Call:        getDisconnectCall(),
	}
}

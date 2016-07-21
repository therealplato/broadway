package cfg

// ServerCfg is the configuration object for the server
var ServerCfg ServerCfgType

// ServerCfgType declares what server config looks like
type ServerCfgType struct {
	AuthBearerToken string // a global token required for all requests except GET/POST command/
	SlackToken      string // the expected Slack custom command token.
	ServerHost      string // passed to gin and configures the listen address of the server
	SlackWebhook    string // your team's slack incoming message webhook URL
	ManifestsPath   string // the absolute path where manifest files are read from
	PlaybooksPath   string // the absolute path where playbook files are read from
}

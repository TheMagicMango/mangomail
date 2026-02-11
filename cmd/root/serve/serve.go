// (c) Magic Mango and individual authors
// SPDX-License-Identifier: Apache-2.0

package serve

import (
	"context"

	"github.com/TheMagicMango/mangomail/configs"
	"github.com/TheMagicMango/mangomail/internal/infra/web"
	"github.com/TheMagicMango/mangomail/pkg/events"
	"github.com/TheMagicMango/mangomail/pkg/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const serviceName = "mangomail"

var (
	port       string
	apiKey     string
	apiKeyFile string
	cfg        *configs.MangomailConfig
)

var Cmd = &cobra.Command{
	Use:     "serve [--port PORT]",
	Short:   "Start the MangoMail HTTP API server",
	Long:    `Starts an HTTP server that exposes a REST API for sending emails programmatically.`,
	Example: examples,
	Args:    cobra.NoArgs,
	RunE:    run,
}

const examples = `# Start HTTP API server:
mangomail serve

# Start on a specific port:
mangomail serve --port 3000`

// @title           MangoMail API
// @version         1.0
// @description     HTTP API for sending emails programmatically via Resend
//
// @contact.name   Magic Mango
// @contact.url    https://github.com/TheMagicMango/mangomail
//
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
//
// @host      mangomail.up.railway.app
// @BasePath  /
//
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key

func init() {
	Cmd.Flags().StringVar(&port, "port", "", "HTTP server port (default: 8081)")
	cobra.CheckErr(viper.BindPFlag("MANGOMAIL_HTTP_ADDRESS", Cmd.Flags().Lookup("port")))

	Cmd.Flags().StringVar(&apiKey, "api-key", "", "API key for HTTP API authentication")
	cobra.CheckErr(viper.BindPFlag("MANGOMAIL_API_KEY", Cmd.Flags().Lookup("api-key")))

	Cmd.Flags().StringVar(&apiKeyFile, "api-key-file", "", "Path to file containing API key")
	cobra.CheckErr(viper.BindPFlag("MANGOMAIL_API_KEY_FILE", Cmd.Flags().Lookup("api-key-file")))

	Cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		var err error
		cfg, err = configs.LoadMangomailConfig()
		if err != nil {
			return err
		}
		return nil
	}
}

func run(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.MaxStartupTime)

	defer cancel()

	createInfo := &web.CreateInfo{
		CreateInfo: service.CreateInfo{
			Name:                 serviceName,
			LogLevel:             cfg.LogLevel,
			LogColor:             cfg.LogColor,
			EnableSignalHandling: true,
			TelemetryCreate:      true,
			TelemetryAddress:     cfg.TelemetryAddress,
		},
		Config: *cfg,
	}

	createInfo.EventDispatcher = events.NewEventDispatcher()

	service, err := web.Create(ctx, createInfo)
	cobra.CheckErr(err)

	// Start service
	return service.Serve()
}

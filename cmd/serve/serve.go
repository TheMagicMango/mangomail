package serve

import (
	"fmt"

	"github.com/TheMagicMango/mangomail/cmd/root"
	"github.com/TheMagicMango/mangomail/configs"
	"github.com/TheMagicMango/mangomail/internal/domain/event"
	"github.com/TheMagicMango/mangomail/internal/domain/event/handler"
	"github.com/TheMagicMango/mangomail/internal/infra/http/controller"
	"github.com/TheMagicMango/mangomail/internal/infra/http/router"
	"github.com/TheMagicMango/mangomail/internal/infra/http/server"
	"github.com/TheMagicMango/mangomail/internal/infra/reader/file"
	"github.com/TheMagicMango/mangomail/internal/usecase"
	"github.com/TheMagicMango/mangomail/pkg/events"
	"github.com/resend/resend-go/v2"
	"github.com/spf13/cobra"
)

var (
	port         int
	from         string
	templatePath string
	cfg          *configs.MangomailConfig
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the MangoMail HTTP API server",
	Long:  `Starts an HTTP server that exposes a REST API for sending emails programmatically.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		cfg, err = configs.LoadMangomailConfig()
		if err != nil {
			return err
		}
		return nil
	},
	RunE: run,
}

func init() {
	serveCmd.Flags().IntVar(&port, "port", 8081, "Port to listen on")
	serveCmd.Flags().StringVar(&from, "from", "", "Sender email address (required)")
	cobra.CheckErr(serveCmd.MarkFlagRequired("from"))
	serveCmd.Flags().StringVar(&templatePath, "template", "examples/template-email.html", "Path to HTML email template")

	root.Cmd.AddCommand(serveCmd)
}

func run(cmd *cobra.Command, args []string) error {
	fileReader := file.NewFileReader()
	htmlTemplate, err := fileReader.LoadHTML(templatePath)
	if err != nil {
		return fmt.Errorf("failed to load HTML template: %w", err)
	}

	resendClient := resend.NewClient(cfg.MangomailResendApiKey.Value)

	eventDispatcher := events.NewEventDispatcher()
	emailHandler := handler.NewEmailSentHandler(resendClient)

	if err := eventDispatcher.Register(event.NewEmailSent().GetName(), emailHandler); err != nil {
		return fmt.Errorf("failed to register email handler: %w", err)
	}

	sendEmailUC := usecase.NewSendEmailUseCase(eventDispatcher)
	emailController := controller.NewEmailController(sendEmailUC, from, htmlTemplate)
	r := router.NewRouter(emailController, cfg)

	return server.Start(r, port)
}

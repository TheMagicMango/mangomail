package root

import (
	"fmt"

	"github.com/TheMagicMango/mangomail/configs"
	"github.com/TheMagicMango/mangomail/internal/domain/event"
	"github.com/TheMagicMango/mangomail/internal/domain/event/handler"
	"github.com/TheMagicMango/mangomail/internal/infra/reader/file"
	"github.com/TheMagicMango/mangomail/internal/infra/version"
	"github.com/TheMagicMango/mangomail/internal/usecase"
	"github.com/TheMagicMango/mangomail/pkg/events"
	"github.com/resend/resend-go/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const serviceName = "mangomail"

var (
	htmlPath         string
	samplePath       string
	from             string
	subject          string
	replyTo          string
	attachments      []string
	resendApiKey     string
	resendApiKeyFile string
	logLevel         string
	rateLimit        uint64
	cfg              *configs.MangomailConfig
)

var Cmd = &cobra.Command{
	Use:     serviceName,
	Short:   "Send bulk emails using HTML templates and CSV data",
	Long:    `MangoMail is a CLI tool for sending bulk emails using HTML templates and CSV sample data with Resend.`,
	Version: version.BuildVersion,
	Args:    cobra.ExactArgs(1),
	Example: `mangomail my-campaign --html templates/welcome.html --sample data/sample.csv --from "sender@example.com" --subject "Welcome {{name}}! --resend-api-key-file ~/.mangomail/secrets/resend_api_key"`,
	RunE: run,
}

func init() {
	configs.SetDefaults()

	Cmd.Flags().StringVar(&htmlPath, "html", "", "Path to HTML template file (required)")
	cobra.CheckErr(Cmd.MarkFlagRequired("html"))

	Cmd.Flags().StringVar(&samplePath, "sample", "", "Path to CSV sample file (required)")
	cobra.CheckErr(Cmd.MarkFlagRequired("sample"))

	Cmd.Flags().StringVar(&from, "from", "", "Sender email address (required)")
	cobra.CheckErr(Cmd.MarkFlagRequired("from"))

	Cmd.Flags().StringVar(&subject, "subject", "", "Email subject (supports {{placeholders}})")
	cobra.CheckErr(Cmd.MarkFlagRequired("subject"))

	Cmd.Flags().StringVar(&replyTo, "reply-to", "", "Reply-to email address")
	Cmd.Flags().StringSliceVar(&attachments, "attachments", []string{}, "Comma-separated list of attachment URLs")

	Cmd.Flags().StringVar(&resendApiKey, "resend-api-key", "", "Resend API key for sending emails")
	cobra.CheckErr(viper.BindPFlag(configs.MANGOMAIL_RESEND_API_KEY, Cmd.Flags().Lookup("resend-api-key")))

	Cmd.Flags().StringVar(&resendApiKeyFile, "resend-api-key-file", "", "Path to file containing Resend API key")
	cobra.CheckErr(viper.BindPFlag(configs.MANGOMAIL_RESEND_API_KEY_FILE, Cmd.Flags().Lookup("resend-api-key-file")))

	Cmd.Flags().StringVar(&logLevel, "log-level", "info", "Log level: debug, info, warn or error")
	cobra.CheckErr(viper.BindPFlag(configs.MANGOMAIL_LOG_LEVEL, Cmd.Flags().Lookup("log-level")))

	Cmd.Flags().Uint64Var(&rateLimit, "rate-limit", 2, "Maximum number of email requests per second")
	cobra.CheckErr(viper.BindPFlag(configs.MANGOMAIL_RATE_LIMIT, Cmd.Flags().Lookup("rate-limit")))

	Cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		var err error
		cfg, err = configs.LoadMangomailConfig()
		if err != nil {
			return err
		}
		return nil
	}

	Cmd.DisableAutoGenTag = true
}

func run(cmd *cobra.Command, args []string) error {
	campaignName := args[0]

	fileReader := file.NewFileReader()

	resendClient := resend.NewClient(cfg.MangomailResendApiKey.Value)

	eventDispatcher := events.NewEventDispatcher()
	emailSentEvent := event.NewEmailSent()
	emailHandler := handler.NewEmailSentHandler(resendClient)

	if err := eventDispatcher.Register(emailSentEvent.GetName(), emailHandler); err != nil {
		return fmt.Errorf("failed to register email handler: %w", err)
	}

	uc := usecase.NewSendCampaignUseCase(emailSentEvent, eventDispatcher, fileReader)

	_, err := uc.Execute(usecase.SendCampaignInputDTO{
		CampaignName: campaignName,
		HTMLPath:     htmlPath,
		SamplePath:   samplePath,
		From:         from,
		Subject:      subject,
		ReplyTo:      replyTo,
		Attachments:  attachments,
		RateLimit:    cfg.MangomailRateLimit,
	})

	return err
}

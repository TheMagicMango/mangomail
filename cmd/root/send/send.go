// (c) Magic Mango and individual authors
// SPDX-License-Identifier: Apache-2.0

package send

import (
	"github.com/TheMagicMango/mangomail/configs"
	"github.com/TheMagicMango/mangomail/internal/domain/event"
	"github.com/TheMagicMango/mangomail/internal/domain/event/handler"
	"github.com/TheMagicMango/mangomail/internal/infra/version"
	"github.com/TheMagicMango/mangomail/internal/usecase"
	"github.com/TheMagicMango/mangomail/pkg/events"
	"github.com/TheMagicMango/mangomail/pkg/reader/file"
	"github.com/go-playground/validator/v10"
	"github.com/resend/resend-go/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	htmlPath         string
	samplePath       string
	from             string
	subject          string
	replyTo          string
	attachments      []string
	resendApiKey     string
	resendApiKeyFile string
	cfg              *configs.MangomailConfig
)

var Cmd = &cobra.Command{
	Use:     "send CAMPAIGN_NAME --html FILE --sample FILE --from EMAIL --subject TEXT [--reply-to EMAIL] [--attachments URLS]",
	Short:   "Send bulk emails using HTML templates and CSV data",
	Long:    `Sends bulk emails using HTML templates and CSV sample data with Resend.`,
	Example: examples,
	Args:    cobra.ExactArgs(1),
	Run:     run,
	Version: version.BuildVersion,
}

const examples = `# Send a campaign with all required flags:
mangomail send my-campaign --html templates/welcome.html --sample data/sample.csv --from "sender@example.com" --subject "Welcome {{name}}"

# Send with reply-to and attachments:
mangomail send newsletter --html templates/newsletter.html --sample data/subscribers.csv --from "news@example.com" --subject "Monthly Newsletter" --reply-to "contact@example.com" --attachments "https://example.com/file1.pdf,https://example.com/file2.pdf"`

func init() {
	Cmd.Flags().StringVar(&htmlPath, "html", "", "Path to HTML template file (required)")
	cobra.CheckErr(Cmd.MarkFlagRequired("html"))

	Cmd.Flags().StringVar(&samplePath, "sample", "", "Path to CSV sample file (required)")
	cobra.CheckErr(Cmd.MarkFlagRequired("sample"))

	Cmd.Flags().StringVar(&from, "from", "", "Sender email address (required)")
	cobra.CheckErr(Cmd.MarkFlagRequired("from"))

	Cmd.Flags().StringVar(&subject, "subject", "", "Email subject (supports {{placeholders}}) (required)")
	cobra.CheckErr(Cmd.MarkFlagRequired("subject"))

	Cmd.Flags().StringVar(&replyTo, "reply-to", "", "Reply-to email address")
	Cmd.Flags().StringSliceVar(&attachments, "attachments", []string{}, "Comma-separated list of attachment URLs")

	Cmd.Flags().StringVar(&resendApiKey, "resend-api-key", "", "Resend API key for sending emails")
	cobra.CheckErr(viper.BindPFlag("MANGOMAIL_RESEND_API_KEY", Cmd.Flags().Lookup("resend-api-key")))

	Cmd.Flags().StringVar(&resendApiKeyFile, "resend-api-key-file", "", "Path to file containing Resend API key")
	cobra.CheckErr(viper.BindPFlag("MANGOMAIL_RESEND_API_KEY_FILE", Cmd.Flags().Lookup("resend-api-key-file")))

	Cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		var err error
		cfg, err = configs.LoadMangomailConfig()
		if err != nil {
			return err
		}
		return nil
	}
}

func run(cmd *cobra.Command, args []string) {
	fileReader := file.NewFileReader()

	htmlContent, err := fileReader.LoadHTML(htmlPath)
	cobra.CheckErr(err)

	csvRows, err := fileReader.LoadCSV(samplePath)
	cobra.CheckErr(err)

	resendClient := resend.NewClient(cfg.ResendApiKey.Value)

	eventDispatcher := events.NewEventDispatcher()
	emailSentEvent := event.NewEmailSent()
	emailHandler := handler.NewEmailSentHandler(resendClient)

	err = eventDispatcher.Register(emailSentEvent.GetName(), emailHandler)
	cobra.CheckErr(err)

	uc := usecase.NewSendCampaignUseCase(emailSentEvent, eventDispatcher)

	input := usecase.SendCampaignInputDTO{
		From:        from,
		Subject:     subject,
		ReplyTo:     replyTo,
		Attachments: attachments,
	}

	validate := validator.New()
	err = validate.Struct(input)
	cobra.CheckErr(err)

	err = uc.Execute(input, htmlContent, csvRows, cfg.ResendRateLimit)
	cobra.CheckErr(err)
}

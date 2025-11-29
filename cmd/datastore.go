package cmd

import (
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/peng225/orochi/internal/datastore/api/server"
	"github.com/peng225/orochi/internal/datastore/handler"
	"github.com/peng225/orochi/internal/datastore/registrar"
	"github.com/peng225/orochi/internal/datastore/service"
	"github.com/spf13/cobra"
)

// datastoreCmd represents the datastore command
var datastoreCmd = &cobra.Command{
	Use:   "datastore",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		levelStr, err := cmd.Flags().GetString(getFlagName())
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: parseLogLevel(levelStr),
		}))
		slog.SetDefault(logger)

		baseURL, err := cmd.Flags().GetString("base-url")
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		var port string
		hostPort := strings.Split(baseURL, ":")
		switch len(hostPort) {
		case 2:
			port = "80"
		case 3:
			port = hostPort[2]
		default:
			slog.Error("Invalid BASE_URL env.", "BASE_URL", baseURL)
			os.Exit(1)
		}
		mgrBaseURL, err := cmd.Flags().GetString("manager-base-url")
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}

		err = registrar.Register(cmd.Context(), baseURL, mgrBaseURL)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}

		objHandler := handler.NewObjectHandler(service.NewObjectStore())
		healthHandler := handler.NewHealthHandler()
		h := server.Handler(struct {
			*handler.ObjectHandler
			*handler.HealthHandler
		}{
			objHandler,
			healthHandler,
		})
		err = http.ListenAndServe(net.JoinHostPort("", port), h)
		if err != nil {
			slog.Error("Server failed to start.", "err", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(datastoreCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// datastoreCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// datastoreCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// FIXME: support https.
	datastoreCmd.Flags().String("base-url", "", "Base url (e.g. http://example.com:8080)")
	datastoreCmd.Flags().String("manager-base-url", "", "Base url for manager (e.g. http://example.com:8080)")
	setLogLevelFlag(datastoreCmd)

	err := datastoreCmd.MarkFlagRequired("base-url")
	if err != nil {
		panic(err)
	}
	err = datastoreCmd.MarkFlagRequired("manager-base-url")
	if err != nil {
		panic(err)
	}
}

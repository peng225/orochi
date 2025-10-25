package cmd

import (
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/peng225/orochi/internal/gateway/api/object/server"
	"github.com/peng225/orochi/internal/gateway/handler"
	"github.com/peng225/orochi/internal/gateway/service"
	"github.com/spf13/cobra"
)

// gatewayCmd represents the gateway command
var gatewayCmd = &cobra.Command{
	Use:   "gateway",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// FIXME: should get the datastore address from somewhere else.
		objHandler := handler.NewObjectHandler(service.NewDataStoreObjectStore("http://localhost:8081"))
		h := server.Handler(objHandler)
		port, err := cmd.Flags().GetString("port")
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		err = http.ListenAndServe(net.JoinHostPort("", port), h)
		if err != nil {
			slog.Error("Server failed to start.", "err", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(gatewayCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// gatewayCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gatewayCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	gatewayCmd.Flags().StringP("port", "p", "8080", "Port number")
}

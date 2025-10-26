package cmd

import (
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/peng225/orochi/internal/datastore/api/server"
	"github.com/peng225/orochi/internal/datastore/handler"
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
		objHandler := handler.NewObjectHandler(service.NewObjectStore())
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
	rootCmd.AddCommand(datastoreCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// datastoreCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// datastoreCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	datastoreCmd.Flags().StringP("port", "p", "8082", "Port number")
}

/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/peng225/orochi/internal/manager/api/server"
	"github.com/peng225/orochi/internal/manager/handler"
	"github.com/peng225/orochi/internal/manager/infra/postgresql"
	"github.com/peng225/orochi/internal/manager/service"
	"github.com/spf13/cobra"
)

// managerCmd represents the manager command
var managerCmd = &cobra.Command{
	Use:   "manager",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		dsRepo := postgresql.NewPSQLDatastoreRepository()
		defer dsRepo.Close()
		dsHandler := handler.NewDatastoreHandler(
			service.NewDatastoreService(dsRepo),
		)
		h := server.Handler(dsHandler)
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
	rootCmd.AddCommand(managerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// managerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// managerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	managerCmd.Flags().StringP("port", "p", "8080", "Port number")
}

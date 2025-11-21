package cmd

import (
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/peng225/orochi/internal/gateway/api/server"
	"github.com/peng225/orochi/internal/gateway/handler"
	"github.com/peng225/orochi/internal/gateway/infra/datastore"
	"github.com/peng225/orochi/internal/gateway/infra/postgresql"
	"github.com/peng225/orochi/internal/gateway/service"
	"github.com/peng225/orochi/internal/pkg/psqlutil"
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
		levelStr, err := cmd.Flags().GetString(getFlagName())
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: parseLogLevel(levelStr),
		}))
		slog.SetDefault(logger)
		port, err := cmd.Flags().GetString("port")
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}

		db := psqlutil.InitDB()
		defer db.Close()
		tx := psqlutil.NewTransaction(db)
		dsRepo := postgresql.NewDatastoreRepository(db)
		omRepo := postgresql.NewObjectMetadataRepository(db)
		ovRepo := postgresql.NewObjectVersionRepository(db)
		bucketRepo := postgresql.NewBucketRepository(db)
		lgRepo := postgresql.NewLocationGroupRepository(db)
		eccRepo := postgresql.NewECConfigRepository(db)
		objService := service.NewObjectStore(
			tx, nil, datastore.NewClientFactory(),
			dsRepo, omRepo, ovRepo, bucketRepo, lgRepo, eccRepo,
		)
		objHandler := handler.NewObjectHandler(objService)
		h := server.Handler(objHandler)
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
	gatewayCmd.Flags().StringP("port", "p", "8081", "Port number")
	setLogLevelFlag(gatewayCmd)
}

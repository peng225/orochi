package cmd

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/peng225/orochi/internal/async/infra/datastore"
	"github.com/peng225/orochi/internal/async/infra/gateway"
	"github.com/peng225/orochi/internal/async/infra/postgresql"
	"github.com/peng225/orochi/internal/async/process"
	"github.com/peng225/orochi/internal/pkg/psqlutil"
	"github.com/spf13/cobra"
)

// asyncCmd represents the async command
var asyncCmd = &cobra.Command{
	Use:   "async",
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
		period, err := cmd.Flags().GetDuration("period")
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		gwBaseURLs, err := cmd.Flags().GetStringSlice("gateway-base-url")
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}

		db := psqlutil.InitDB()
		defer db.Close()
		tx := psqlutil.NewTransaction(db)
		bucketRepo := postgresql.NewBucketRepository(db)
		jobRepo := postgresql.NewJobRepository(db)
		dsRepo := postgresql.NewDatastoreRepository(db)
		gwClients := make([]process.GatewayClient, len(gwBaseURLs))
		for i := range len(gwClients) {
			gwClients[i] = gateway.NewClient(gwBaseURLs[i])
		}
		dscFactory := datastore.NewClientFactory()

		p := process.NewProcessor(period, tx, jobRepo, bucketRepo, dsRepo, gwClients, dscFactory)
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
		defer stop()
		p.Start(ctx)
	},
}

func init() {
	rootCmd.AddCommand(asyncCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// asyncCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// asyncCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	asyncCmd.Flags().Duration("period", 10*time.Second, "Time period to process async jobs.")
	asyncCmd.Flags().StringSlice("gateway-base-url", nil, "A list of gateway base URL.")
	setLogLevelFlag(asyncCmd)

	err := asyncCmd.MarkFlagRequired("gateway-base-url")
	if err != nil {
		panic(err)
	}
}

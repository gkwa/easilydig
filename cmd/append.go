package cmd

import (
	"context"
	"os"

	"github.com/gkwa/easilydig/core"
	"github.com/spf13/cobra"
)

var appendCmd = &cobra.Command{
	Use:   "append [files...]",
	Short: "Append usage metrics from JSON files to DynamoDB",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		logger := LoggerFrom(cmd.Context())

		dynamoSvc, err := core.NewDynamoDBService(context.Background())
		if err != nil {
			logger.Error(err, "Failed to create DynamoDB service")
			os.Exit(1)
		}

		repo := core.NewMetricRepository(dynamoSvc)
		service := core.NewMetricService(repo, logger)

		if err := service.AppendFiles(context.Background(), args); err != nil {
			logger.Error(err, "Failed to append metrics")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(appendCmd)
}

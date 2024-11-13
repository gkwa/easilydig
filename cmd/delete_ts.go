package cmd

import (
	"context"
	"os"

	"github.com/gkwa/easilydig/core"
	"github.com/spf13/cobra"
)

var deleteTsCmd = &cobra.Command{
	Use:   "deletets [timestamp]",
	Short: "Delete records containing specified timestamp",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		logger := LoggerFrom(cmd.Context())

		dynamoSvc, err := core.NewDynamoDBService(context.Background())
		if err != nil {
			logger.Error(err, "Failed to create DynamoDB service")
			os.Exit(1)
		}

		repo := core.NewMetricRepository(dynamoSvc)
		service := core.NewMetricService(repo, logger)

		if err := service.DeleteByTimestamp(context.Background(), args[0]); err != nil {
			logger.Error(err, "Failed to delete records")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteTsCmd)
}

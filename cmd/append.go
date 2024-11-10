package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/gkwa/easilydig/core"
	"github.com/spf13/cobra"
)

var appendCmd = &cobra.Command{
	Use:   "append [file]",
	Short: "Append usage metrics from JSON file to DynamoDB",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		logger := LoggerFrom(cmd.Context())

		data, err := os.ReadFile(args[0])
		if err != nil {
			logger.Error(err, "Failed to read input file")
			os.Exit(1)
		}

		var metric core.UsageMetric
		if err := json.Unmarshal(data, &metric); err != nil {
			logger.Error(err, "Failed to parse JSON")
			os.Exit(1)
		}

		dynamoSvc, err := core.NewDynamoDBService(context.Background())
		if err != nil {
			logger.Error(err, "Failed to create DynamoDB service")
			os.Exit(1)
		}

		repo := core.NewMetricRepository(dynamoSvc)
		if err := repo.Append(context.Background(), metric); err != nil {
			logger.Error(err, fmt.Sprintf("Failed to append metric to DynamoDB: %v", err))
			os.Exit(1)
		}

		logger.Info("Successfully appended metric to DynamoDB")
	},
}

func init() {
	rootCmd.AddCommand(appendCmd)
}

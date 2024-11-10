package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/gkwa/easilydig/core"
	"github.com/spf13/cobra"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch all usage metrics from DynamoDB",
	Run: func(cmd *cobra.Command, args []string) {
		logger := LoggerFrom(cmd.Context())

		dynamoSvc, err := core.NewDynamoDBService(context.Background())
		if err != nil {
			logger.Error(err, "Failed to create DynamoDB service")
			os.Exit(1)
		}

		repo := core.NewMetricRepository(dynamoSvc)
		metrics, err := repo.FetchAll(context.Background())
		if err != nil {
			logger.Error(err, "Failed to fetch metrics from DynamoDB")
			os.Exit(1)
		}

		output, err := json.MarshalIndent(metrics, "", "  ")
		if err != nil {
			logger.Error(err, "Failed to marshal metrics to JSON")
			os.Exit(1)
		}

		fmt.Println(string(output))
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)
}

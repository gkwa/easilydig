package core

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/go-logr/logr"
)

type MetricService struct {
	repo   MetricRepositoryInterface
	logger logr.Logger
}

func NewMetricService(repo MetricRepositoryInterface, logger logr.Logger) *MetricService {
	return &MetricService{
		repo:   repo,
		logger: logger,
	}
}

func (s *MetricService) AppendFiles(ctx context.Context, patterns []string) error {
	var allFiles []string

	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return fmt.Errorf("failed to expand glob pattern %s: %w", pattern, err)
		}

		allFiles = append(allFiles, matches...)
	}

	if len(allFiles) == 0 {
		return fmt.Errorf("no files found matching the provided patterns")
	}

	for _, file := range allFiles {
		if err := s.appendFile(ctx, file); err != nil {
			s.logger.Error(err, "Failed to process file", "file", file)
			continue
		}

		s.logger.Info("Successfully appended metric to DynamoDB", "file", file)
	}

	return nil
}

func (s *MetricService) appendFile(ctx context.Context, file string) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var metric UsageMetric
	if err := json.Unmarshal(data, &metric); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	if err := s.repo.Append(ctx, metric); err != nil {
		return fmt.Errorf("failed to append metric to DynamoDB: %w", err)
	}

	return nil
}

func (s *MetricService) DeleteByTimestamp(ctx context.Context, ts string) error {
	if err := s.repo.DeleteByTimestamp(ctx, ts); err != nil {
		return fmt.Errorf("failed to delete records: %w", err)
	}

	s.logger.Info("Successfully deleted records", "timestamp", ts)

	return nil
}

func NewDynamoDBService(ctx context.Context) (*dynamodb.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		return nil, err
	}

	return dynamodb.NewFromConfig(cfg), nil
}

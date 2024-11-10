package core

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type MetricRepositoryInterface interface {
	Append(ctx context.Context, metric UsageMetric) error
	FetchAll(ctx context.Context) ([]UsageMetric, error)
}

type DynamoDBService interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	Scan(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
}

type MetricRepository struct {
	dynamoSvc DynamoDBService
	tableName string
}

func NewMetricRepository(svc DynamoDBService) *MetricRepository {
	return &MetricRepository{
		dynamoSvc: svc,
		tableName: "usage-metrics",
	}
}

func (r *MetricRepository) Append(ctx context.Context, metric UsageMetric) error {
	item, err := attributevalue.MarshalMap(metric)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(r.tableName),
	}

	_, err = r.dynamoSvc.PutItem(ctx, input)
	return err
}

func (r *MetricRepository) FetchAll(ctx context.Context) ([]UsageMetric, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(r.tableName),
	}

	result, err := r.dynamoSvc.Scan(ctx, input)
	if err != nil {
		return nil, err
	}

	var metrics []UsageMetric
	err = attributevalue.UnmarshalListOfMaps(result.Items, &metrics)
	if err != nil {
		return nil, err
	}

	return metrics, nil
}

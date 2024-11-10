package core

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDynamoDBService struct {
	mock.Mock
}

func (m *MockDynamoDBService) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*dynamodb.PutItemOutput), args.Error(1)
}

func (m *MockDynamoDBService) Scan(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*dynamodb.ScanOutput), args.Error(1)
}

func TestMetricRepository_Append(t *testing.T) {
	mockSvc := &MockDynamoDBService{}
	repo := NewMetricRepository(mockSvc)

	tests := []struct {
		name    string
		metric  UsageMetric
		mockFn  func()
		wantErr bool
	}{
		{
			name: "successful append",
			metric: UsageMetric{
				Date:        "2024-01-01",
				Amount:      100.0,
				AmountUnits: "GB",
				ScrapedAt:   time.Now(),
			},
			mockFn: func() {
				mockSvc.On("PutItem", mock.Anything, mock.Anything).Return(&dynamodb.PutItemOutput{}, nil)
			},
			wantErr: false,
		},
		{
			name: "dynamodb error",
			metric: UsageMetric{
				Date:        "2024-01-01",
				Amount:      100.0,
				AmountUnits: "GB",
				ScrapedAt:   time.Now(),
			},
			mockFn: func() {
				mockSvc.On("PutItem", mock.Anything, mock.Anything).Return(&dynamodb.PutItemOutput{}, errors.New("dynamodb error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc.Mock = mock.Mock{} // Reset mock
			tt.mockFn()

			err := repo.Append(context.Background(), tt.metric)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestMetricRepository_FetchAll(t *testing.T) {
	mockSvc := &MockDynamoDBService{}
	repo := NewMetricRepository(mockSvc)

	tests := []struct {
		name    string
		mockFn  func()
		want    []UsageMetric
		wantErr bool
	}{
		{
			name: "successful fetch",
			mockFn: func() {
				mockSvc.On("Scan", mock.Anything, mock.Anything).Return(&dynamodb.ScanOutput{
					Items: []map[string]types.AttributeValue{
						{
							"date":        &types.AttributeValueMemberS{Value: "2024-01-01"},
							"amount":      &types.AttributeValueMemberN{Value: "100"},
							"amountUnits": &types.AttributeValueMemberS{Value: "GB"},
							"scrapedAt":   &types.AttributeValueMemberS{Value: "2024-01-01T00:00:00Z"},
						},
					},
				}, nil)
			},
			want: []UsageMetric{
				{
					Date:        "2024-01-01",
					Amount:      100.0,
					AmountUnits: "GB",
					ScrapedAt:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			wantErr: false,
		},
		{
			name: "dynamodb error",
			mockFn: func() {
				mockSvc.On("Scan", mock.Anything, mock.Anything).Return(&dynamodb.ScanOutput{}, errors.New("dynamodb error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc.Mock = mock.Mock{} // Reset mock
			tt.mockFn()

			got, err := repo.FetchAll(context.Background())

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

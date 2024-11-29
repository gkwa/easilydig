package core

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDynamoDBService struct {
	mock.Mock
}

func (m *MockDynamoDBService) PutItem(
	ctx context.Context,
	params *dynamodb.PutItemInput,
	optFns ...func(*dynamodb.Options),
) (*dynamodb.PutItemOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*dynamodb.PutItemOutput), args.Error(1)
}

func (m *MockDynamoDBService) Scan(
	ctx context.Context,
	params *dynamodb.ScanInput,
	optFns ...func(*dynamodb.Options),
) (*dynamodb.ScanOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*dynamodb.ScanOutput), args.Error(1)
}

func (m *MockDynamoDBService) DeleteItem(
	ctx context.Context,
	params *dynamodb.DeleteItemInput,
	optFns ...func(*dynamodb.Options),
) (*dynamodb.DeleteItemOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*dynamodb.DeleteItemOutput), args.Error(1)
}

func TestMetricRepository_DeleteByTimestamp(t *testing.T) {
	mockSvc := &MockDynamoDBService{}
	repo := NewMetricRepository(mockSvc)

	tests := []struct {
		name    string
		ts      string
		mockFn  func()
		wantErr bool
	}{
		{
			name: "successful delete",
			ts:   "2024-01-01T00:00:00Z",
			mockFn: func() {
				mockSvc.On("Scan", mock.Anything, mock.Anything).Return(&dynamodb.ScanOutput{
					Items: []map[string]types.AttributeValue{
						{
							"date": &types.AttributeValueMemberS{Value: "2024-01-01"},
							"scrapedAt": &types.AttributeValueMemberS{
								Value: "2024-01-01T00:00:00Z",
							},
						},
					},
				}, nil)
				mockSvc.On("DeleteItem", mock.Anything, mock.Anything).
					Return(&dynamodb.DeleteItemOutput{}, nil)
			},
			wantErr: false,
		},
		{
			name: "scan error",
			ts:   "2024-01-01T00:00:00Z",
			mockFn: func() {
				mockSvc.On("Scan", mock.Anything, mock.Anything).
					Return(&dynamodb.ScanOutput{}, errors.New("scan error"))
			},
			wantErr: true,
		},
		{
			name: "delete error",
			ts:   "2024-01-01T00:00:00Z",
			mockFn: func() {
				mockSvc.On("Scan", mock.Anything, mock.Anything).Return(&dynamodb.ScanOutput{
					Items: []map[string]types.AttributeValue{
						{
							"date": &types.AttributeValueMemberS{Value: "2024-01-01"},
							"scrapedAt": &types.AttributeValueMemberS{
								Value: "2024-01-01T00:00:00Z",
							},
						},
					},
				}, nil)
				mockSvc.On("DeleteItem", mock.Anything, mock.Anything).
					Return(&dynamodb.DeleteItemOutput{}, errors.New("delete error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc.Mock = mock.Mock{} // Reset mock

			tt.mockFn()

			err := repo.DeleteByTimestamp(context.Background(), tt.ts)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

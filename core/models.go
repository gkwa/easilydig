package core

import "time"

type UsageMetric struct {
	Date         string    `dynamodbav:"date"         json:"date"`
	Amount       float64   `dynamodbav:"amount"       json:"amount"`
	AmountUnits  string    `dynamodbav:"amountUnits"  json:"amountUnits"`
	Total        float64   `dynamodbav:"total"        json:"total"`
	TotalUnits   string    `dynamodbav:"totalUnits"   json:"totalUnits"`
	Overage      float64   `dynamodbav:"overage"      json:"overage"`
	OverageUnits string    `dynamodbav:"overageUnits" json:"overageUnits"`
	ScrapedAt    time.Time `dynamodbav:"scrapedAt"    json:"scrapedAt"`
}

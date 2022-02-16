package data

import (
	"context"
	"encoding/json"
	"github.com/mmcdole/gofeed"
	"github.com/olivere/elastic/v7"
	"github.com/pkg/errors"
)

const indexTitle = "index-items"

var searchingFields = []string{
	"title",
	"description",
	"published",
	"guid",
	"categories",
}

type ItemsQ struct {
	client  *elastic.Client
	clauses []string
}

func NewItemsQ(client *elastic.Client) *ItemsQ {
	return &ItemsQ{
		client:  client,
		clauses: searchingFields,
	}
}

func (q *ItemsQ) BulkIndex(items []gofeed.Item) error {
	if len(items) == 0 {
		return nil
	}

	var requests []elastic.BulkableRequest

	for _, item := range items {
		requests = append(requests, elastic.NewBulkIndexRequest().
			Doc(item).
			Id(item.GUID))
	}

	_, err := q.client.Bulk().
		Index(indexTitle).
		Add(requests...).
		Do(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (q *ItemsQ) Search(key string, limit int, offset int) ([]gofeed.Item, error) {
	if limit == 0 {
		limit = 10
	}

	items := make([]gofeed.Item, 0)

	query := elastic.NewMultiMatchQuery(key, q.clauses...).Type("phrase_prefix")

	res, err := q.client.
		Search().
		Index(indexTitle).
		Query(query).
		From(offset).
		Size(limit).
		Do(context.Background())
	if err != nil {
		return items, errors.Wrap(err, "failed to exec search query")
	}

	for _, hit := range res.Hits.Hits {
		var item gofeed.Item
		err = json.Unmarshal(hit.Source, &item)
		if err != nil {
			return items, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (q *ItemsQ) Aggregation(key string) ([]*elastic.AggregationBucketHistogramItem, error) {
	query := elastic.NewMultiMatchQuery(key, q.clauses...).Type("phrase_prefix")
	agg := elastic.NewDateHistogramAggregation().Field("publishedParsed").
		MinDocCount(0).CalendarInterval("1d")

	res, err := q.client.
		Search().
		Index(indexTitle).
		Query(query).
		Aggregation("dates_with_holes", agg).
		Do(context.Background())
	if err != nil {
		return nil, err
	}

	h, ok := res.Aggregations.DateHistogram("dates_with_holes")
	if !ok {
		return nil, errors.New("failed to get hist")
	}

	return h.Buckets, nil
}

package service

import (
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/pkg/errors"
	"os"
	"strconv"

	"encoding/csv"

	"github.com/zh0glikk/elastic-app/internal/config"
	"github.com/zh0glikk/elastic-app/internal/service/data"
)

func Aggregate(cfg *config.Config, key string, output string) error {
	if cfg == nil {
		return errors.New("nil cfg")
	}

	client, err := elastic.NewClient(elastic.SetURL(cfg.ElasticCfg.URL), elastic.SetSniff(false))
	if err != nil {
		return errors.Wrap(err, "failed to init elastic client")
	}

	itemsQ := data.NewItemsQ(client)

	buckets, err := itemsQ.Aggregation(key)
	if err != nil {
		return errors.Wrap(err, "failed to exec aggregation q")
	}

	records := make([][]string, 0)

	for _, b := range buckets {
		records = append(records, []string{
			*b.KeyAsString,
			strconv.FormatInt(b.DocCount, 10),
		})
	}

	f, err := os.Create(fmt.Sprintf("%s.csv", output))

	wr := csv.NewWriter(f)
	wr.WriteAll(records)
	wr.Flush()

	return nil
}

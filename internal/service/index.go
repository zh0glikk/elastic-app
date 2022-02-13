package service

import (
	"encoding/json"
	"fmt"
	"github.com/mmcdole/gofeed"
	"github.com/olivere/elastic/v7"
	"github.com/pkg/errors"
	"os"

	"github.com/zh0glikk/elastic-app/internal/config"
	"github.com/zh0glikk/elastic-app/internal/service/data"
)

func Index(fileTitle string, cfg *config.Config) error {
	if cfg == nil {
		return errors.New("nil cfg")
	}

	client, err := elastic.NewClient(elastic.SetURL(cfg.ElasticCfg.URL), elastic.SetSniff(false))
	if err != nil {
		return errors.Wrap(err, "failed to init elastic client")
	}

	itemsQ := data.NewItemsQ(client)

	dd, err := os.ReadFile(fmt.Sprintf("%s.json", fileTitle))
	if err != nil {
		return errors.Wrap(err, "failed to read file")
	}

	items := make([]gofeed.Item, 0)

	err = json.Unmarshal(dd, &items)
	if err != nil {
		return errors.Wrap(err, "failed to read file")
	}

	//FIXME: could it fail with enormous amount of items????
	err = itemsQ.BulkIndex(items)
	if err != nil {
		return errors.Wrap(err, "failed to bulk index")
	}

	return nil
}

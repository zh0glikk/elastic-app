package service

import (
	"encoding/json"
	"fmt"
	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
	"os"
)

func Convert(fileTitle string) error {
	file, err := os.Open(fmt.Sprintf("%s.xml", fileTitle))
	if err != nil {
		return errors.Wrap(err, "failed to open file")
	}
	defer file.Close()

	fp := gofeed.NewParser()

	feed, err := fp.Parse(file)
	if err != nil {
		return errors.Wrap(err, "failed to parse file")
	}

	res, err := json.Marshal(feed.Items)
	if err != nil {
		return errors.Wrap(err, "failed to marshal items")
	}

	err = os.WriteFile(fmt.Sprintf("%s.json", fileTitle), res, 0644)
	if err != nil {
		return errors.Wrap(err, "failed to write to file")
	}

	return nil
}

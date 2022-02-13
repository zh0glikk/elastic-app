package service

import (
	"encoding/json"
	"fmt"
	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
	"os"
)

func Merge(fileTitles []string, output string) error {
	result := make([]gofeed.Item, 0)

	for _, fileTitle := range fileTitles {
		dd, err := os.ReadFile(fmt.Sprintf("%s.json", fileTitle))
		if err != nil {
			return errors.Wrap(err, "failed to read file")
		}

		items := make([]gofeed.Item, 0)

		err = json.Unmarshal(dd, &items)
		if err != nil {
			return errors.Wrap(err, "failed to read file")
		}

		result = append(result, items...)
	}

	res, err := json.Marshal(result)
	if err != nil {
		return errors.Wrap(err, "failed to marshal items")
	}

	err = os.WriteFile(fmt.Sprintf("%s.json", output), res, 0644)
	if err != nil {
		return errors.Wrap(err, "failed to write to file")
	}

	return nil
}

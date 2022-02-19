package cli

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/zh0glikk/elastic-app/internal/config"
	"github.com/zh0glikk/elastic-app/internal/service"
)

const defaultConfigPath = "./config.json"

func Run(args []string) bool {
	var cfg config.Config

	app := cli.NewApp()

	err := config.SetupConfig(defaultConfigPath, &cfg)
	if err != nil {
		logrus.WithError(err).Error("failed to parse config")
		return false
	}

	//flow:
	//1.convert each file to json
	//2.merge all files to one
	//3.index with provided json
	//4.run web

	app.Commands = cli.Commands{
		{
			Name: "convert",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name: "file",
				},
			},
			Action: func(c *cli.Context) error {
				return service.Convert(c.String("file"))
			},
		},
		{
			Name: "merge",
			Flags: []cli.Flag{
				&cli.StringSliceFlag{
					Name: "file",
				},
				&cli.StringFlag{
					Name:  "output",
					Value: "merged",
				},
			},
			Action: func(c *cli.Context) error {
				return service.Merge(c.StringSlice("file"), c.String("output"))
			},
		},
		{
			Name: "index",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name: "file",
					Aliases: []string{"f"},
				},
			},
			Action: func(c *cli.Context) error {
				return service.Index(c.String("file"), &cfg)
			},
		},
		{
			Name: "web",
			Action: func(_ *cli.Context) error {
				return service.NewService(&cfg).Run(context.Background())
			},
		},
		{
			Name: "aggregate",
			Flags: []cli.Flag{
				&cli.StringSliceFlag{
					Name: "key",
				},
				&cli.StringFlag{
					Name:  "output",
					Value: "aggregated",
				},
			},
			Action: func(c *cli.Context) error {
				return service.Aggregate(&cfg, c.String("key"), c.String("output"))
			},
		},
		{
			Name: "extract",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name: "file",
					Aliases: []string{"f"},
				},
			},
			Action: func(c *cli.Context) error {
				return service.Extract(c.String("file"))
			},
		},
	}

	if err := app.Run(args); err != nil {
		logrus.WithError(err).Error("app finished")
		return false
	}
	return true
}

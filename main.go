package main

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/joho/godotenv"
	"github.com/urfave/cli"
)

var build = "0" // build number set at compile-time

func main() {
	app := cli.NewApp()
	app.Name = "niftycloud-object-storage plugin"
	app.Usage = "niftycloud-object-storage plugin"
	app.Action = run
	app.Version = fmt.Sprintf("1.0.%s", build)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "access-key",
			Usage:  "access key",
			EnvVar: "PLUGIN_ACCESS_KEY,NIFTY_ACCESS_KEY_ID",
		},
		cli.StringFlag{
			Name:   "secret-key",
			Usage:  "secret key",
			EnvVar: "PLUGIN_SECRET_KEY,NIFTY_SECRET_KEY",
		},
		cli.StringFlag{
			Name:   "bucket",
			Usage:  "bucket",
			EnvVar: "PLUGIN_BUCKET,STORAGE_BUCKET",
		},
		cli.StringFlag{
			Name:   "region",
			Usage:  "region",
			Value:  "jp-east-2",
			EnvVar: "PLUGIN_REGION,STORAGE_REGION",
		},
		cli.StringFlag{
			Name:   "acl",
			Usage:  "upload files with acl",
			Value:  "private",
			EnvVar: "PLUGIN_ACL",
		},
		cli.StringFlag{
			Name:   "source",
			Usage:  "upload files from source folder",
			EnvVar: "PLUGIN_SOURCE",
		},
		cli.StringFlag{
			Name:   "target",
			Usage:  "upload files to target folder",
			EnvVar: "PLUGIN_TARGET",
		},
		cli.StringFlag{
			Name:   "strip-prefix",
			Usage:  "strip the prefix from the target",
			EnvVar: "PLUGIN_STRIP_PREFIX",
		},
		cli.StringSliceFlag{
			Name:   "exclude",
			Usage:  "ignore files matching exclude pattern",
			EnvVar: "PLUGIN_EXCLUDE",
		},
		cli.StringFlag{
			Name:  "env-file",
			Usage: "source env file",
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(c *cli.Context) error {
	if c.String("env-file") != "" {
		_ = godotenv.Load(c.String("env-file"))
	}

	plugin := Plugin{
		Key:         c.String("access-key"),
		Secret:      c.String("secret-key"),
		Bucket:      c.String("bucket"),
		Region:      c.String("region"),
		Access:      c.String("acl"),
		Source:      c.String("source"),
		Target:      c.String("target"),
		StripPrefix: c.String("strip-prefix"),
		Exclude:     c.StringSlice("exclude"),
	}

	return plugin.Exec()
}

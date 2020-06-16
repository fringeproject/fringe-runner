package commands

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/fringeproject/fringe-runner/common"
)

type ParseFileCommand struct {
}

func (s *ParseFileCommand) ParseRawFile(filePath string, ctx *common.ModuleContext) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		ips := common.SearchAllIP(line)
		hostnames := []string{}
		for _, hostname := range common.SearchAllHostname(line) {
			found := false
			for _, ip := range ips {
				if hostname == ip {
					found = true
					break
				}
			}

			if !found {
				hostnames = append(hostnames, hostname)
			}
		}

		for _, hostname := range hostnames {
			err = ctx.CreateNewAssetAsHostname(hostname)
			if err != nil {
				logrus.Infof("Cannot create hostname: %s", hostname)
			}
		}

		for _, ip := range ips {
			err = ctx.CreateNewAssetAsIP(ip)
			if err != nil {
				logrus.Infof("Cannot create ip: %s", ip)
			}
		}
	}

	return nil
}

func (s *ParseFileCommand) Execute(c *cli.Context, config *common.FringeConfig) error {
	filePath := c.String("path")

	if !common.FileExists(filePath) {
		return fmt.Errorf("The file does not exists.")
	}

	asset := common.Asset{
		Value: "",
		Type:  "",
	}

	ctx, err := common.NewModuleContext(asset, config)
	if err != nil {
		logrus.Warn("Cannot crate module context.")
		logrus.Debug(err)
		return err
	}

	err = s.ParseRawFile(filePath, ctx)
	if err != nil {
		return err
	}

	assetsJSON, err := json.MarshalIndent(ctx.NewAssets, "", "\t")
	if err != nil {
		return fmt.Errorf("Couldn't format the assets to JSON.")
	}

	fmt.Println(string(assetsJSON))

	return nil
}

func newParseFileCommand() *ParseFileCommand {
	return &ParseFileCommand{}
}

// Register the command
func init() {
	common.RegisterCommand("parse", "Parse a file to find new assets", newParseFileCommand(), []cli.Flag{
		&cli.StringFlag{
			Name:    "path",
			Aliases: []string{"p"},
			Usage:   "Path to the file to parse",
		},
	})
}

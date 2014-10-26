package cmd

import "github.com/codegangsta/cli"

var CmdBackup = cli.Command{
	Name:        "backup",
	Usage:       "backup blog data into archive",
	Description: "backup blog data into archive",
	Action:      DoBackup,
	Flags:       []cli.Flag{},
}

func DoBackup(ctx *cli.Context) {

}

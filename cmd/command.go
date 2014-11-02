package cmd

import "github.com/codegangsta/cli"

var (
	CmdRegister *Command
)

func init() {
	CmdRegister = &Command{
		cmdMap:    make(map[string]cli.Command),
		beforeCmd: []func(*cli.Context){},
		afterCmd:  []func(*cli.Context){},
	}

	CmdRegister.Register("web", CmdWeb)
	CmdRegister.Register("upgrade", CmdUpgrade)
	CmdRegister.Register("backup", CmdBackup)

	CmdRegister.Before(installCheck)
	CmdRegister.Before(upgradeCheck)
}

type Command struct {
	cmdMap    map[string]cli.Command
	beforeCmd []func(*cli.Context)
	afterCmd  []func(*cli.Context)
}

func (cm *Command) Before(fn func(*cli.Context)) {
	cm.beforeCmd = append(cm.beforeCmd, fn)
}

func (cm *Command) After(fn func(*cli.Context)) {
	cm.afterCmd = append(cm.afterCmd, fn)
}

func (cm *Command) Register(name string, cmd cli.Command) {
	cm.cmdMap[name] = cmd
}

func (cm *Command) Get(name string) *cli.Command {
	cmd, ok := cm.cmdMap[name]
	if !ok {
		return nil
	}
	oldAction := cmd.Action
	cmd.Action = func(ctx *cli.Context) {
		for _, fn := range cm.beforeCmd {
			fn(ctx)
		}
		oldAction(ctx)
		for _, fn := range cm.afterCmd {
			fn(ctx)
		}
	}
	return &cmd
}

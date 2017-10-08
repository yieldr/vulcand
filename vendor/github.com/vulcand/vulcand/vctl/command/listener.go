package command

import (
	"github.com/codegangsta/cli"
	"github.com/vulcand/vulcand/engine"
)

func NewListenerCommand(cmd *Command) cli.Command {
	return cli.Command{
		Name:  "listener",
		Usage: "Operations with socket listeners",
		Subcommands: []cli.Command{
			{
				Name:   "ls",
				Usage:  "List all listeners",
				Flags:  []cli.Flag{},
				Action: cmd.printListenersAction,
			},
			{
				Name:  "show",
				Usage: "Show listener details",
				Flags: []cli.Flag{
					cli.StringFlag{Name: "id", Usage: "listener id"},
				},
				Action: cmd.printListenerAction,
			},
			{
				Name:  "upsert",
				Usage: "Update or insert a listener",
				Flags: append([]cli.Flag{
					cli.StringFlag{Name: "id", Usage: "id"},
					cli.StringFlag{Name: "proto", Usage: "protocol, either http or https"},
					cli.StringFlag{Name: "net", Value: "tcp", Usage: "network, tcp or unix"},
					cli.StringFlag{Name: "addr", Value: "tcp", Usage: "address to bind to, e.g. 'localhost:31000'"},
					cli.StringFlag{Name: "scope", Usage: "scope expression limits the listener, e.g. 'Hostname(`myhost`)'"},
					cli.StringFlag{Name: "proxy-header", Value: "none", Usage: "none or PROXY_V1"},
				}, getTLSFlags()...),
				Action: cmd.upsertListenerAction,
			},
			{
				Name:   "rm",
				Usage:  "Remove a listener",
				Action: cmd.deleteListenerAction,
				Flags: []cli.Flag{
					cli.StringFlag{Name: "id", Usage: "id"},
				},
			},
		},
	}
}

func (cmd *Command) upsertListenerAction(c *cli.Context) error {
	var settings *engine.HTTPSListenerSettings
	if c.String("proto") == engine.HTTPS {
		s, err := getTLSSettings(c)
		if err != nil {
			return err
		}
		settings = &engine.HTTPSListenerSettings{TLS: *s}
	}
	listener, err := engine.NewListener(c.String("id"), c.String("proto"), c.String("net"), c.String("addr"), c.String("scope"), c.String("proxy-header"), settings)
	if err != nil {
		return err
	}
	if err := cmd.client.UpsertListener(*listener); err != nil {
		return err
	}
	cmd.printOk("listener upserted")
	return nil
}

func (cmd *Command) deleteListenerAction(c *cli.Context) error {
	if err := cmd.client.DeleteListener(engine.ListenerKey{Id: c.String("id")}); err != nil {
		return err
	}
	cmd.printOk("listener deleted")
	return nil
}

func (cmd *Command) printListenersAction(c *cli.Context) error {
	ls, err := cmd.client.GetListeners()
	if err != nil {
		return err
	}
	cmd.printListeners(ls)
	return nil
}

func (cmd *Command) printListenerAction(c *cli.Context) error {
	l, err := cmd.client.GetListener(engine.ListenerKey{Id: c.String("id")})
	if err != nil {
		return err
	}
	cmd.printListener(l)
	return nil
}

package command

import (
	"github.com/codegangsta/cli"
	"github.com/vulcand/vulcand/engine"
)

func NewBackendCommand(cmd *Command) cli.Command {
	return cli.Command{
		Name:  "backend",
		Usage: "Operations with backends",
		Subcommands: []cli.Command{
			{
				Name:   "upsert",
				Usage:  "Update or insert a new backend to vulcan",
				Action: cmd.upsertBackendAction,
				Flags: append(append([]cli.Flag{
					cli.StringFlag{Name: "id", Usage: "backend id"}},
					backendOptions()...),
					getTLSFlags()...),
			},
			{
				Name:   "rm",
				Usage:  "Remove backend from vulcan",
				Action: cmd.deleteBackendAction,
				Flags: []cli.Flag{
					cli.StringFlag{Name: "id", Usage: "backend id"},
				},
			},
			{
				Name:   "ls",
				Usage:  "List backends",
				Action: cmd.listBackendsAction,
			},
			{
				Name:   "show",
				Usage:  "Show backend",
				Action: cmd.printBackendAction,
				Flags: []cli.Flag{
					cli.StringFlag{Name: "id", Usage: "backend id"},
				},
			},
		},
	}
}

func (cmd *Command) upsertBackendAction(c *cli.Context) error {
	settings, err := getBackendSettings(c)
	if err != nil {
		return err
	}
	b, err := engine.NewHTTPBackend(c.String("id"), settings)
	if err != nil {
		return err
	}
	cmd.printResult("%s upserted", b, cmd.client.UpsertBackend(*b))
	return nil
}

func (cmd *Command) deleteBackendAction(c *cli.Context) error {
	if err := cmd.client.DeleteBackend(engine.BackendKey{Id: c.String("id")}); err != nil {
		return err
	}
	cmd.printOk("backend deleted")
	return nil
}

func (cmd *Command) printBackendAction(c *cli.Context) error {
	bk := engine.BackendKey{Id: c.String("id")}
	b, err := cmd.client.GetBackend(bk)
	if err != nil {
		return err
	}
	srvs, err := cmd.client.GetServers(bk)
	if err != nil {
		return err
	}
	cmd.printBackend(b, srvs)
	return nil
}

func (cmd *Command) listBackendsAction(c *cli.Context) error {
	out, err := cmd.client.GetBackends()
	if err != nil {
		return err
	}
	cmd.printBackends(out)
	return nil
}

func getBackendSettings(c *cli.Context) (engine.HTTPBackendSettings, error) {
	s := engine.HTTPBackendSettings{}

	s.Timeouts.Read = c.Duration("readTimeout").String()
	s.Timeouts.Dial = c.Duration("dialTimeout").String()
	s.Timeouts.TLSHandshake = c.Duration("handshakeTimeout").String()

	s.KeepAlive.Period = c.Duration("keepAlivePeriod").String()
	s.KeepAlive.MaxIdleConnsPerHost = c.Int("maxIdleConns")

	tlsSettings, err := getTLSSettings(c)
	if err != nil {
		return s, err
	}
	s.TLS = tlsSettings
	return s, nil
}

func backendOptions() []cli.Flag {
	return []cli.Flag{
		// Timeouts
		cli.DurationFlag{Name: "readTimeout", Usage: "read timeout"},
		cli.DurationFlag{Name: "dialTimeout", Usage: "dial timeout"},
		cli.DurationFlag{Name: "handshakeTimeout", Usage: "TLS handshake timeout"},

		// Keep-alive parameters
		cli.StringFlag{Name: "keepAlivePeriod", Usage: "keep-alive period"},
		cli.IntFlag{Name: "maxIdleConns", Usage: "maximum idle connections per host"},
	}
}

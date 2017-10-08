package command

import (
	"fmt"
	"io"
	"os"

	"github.com/codegangsta/cli"
	"github.com/vulcand/vulcand/secret"
)

func NewKeyCommand(cmd *Command) cli.Command {
	return cli.Command{
		Name:  "secret",
		Usage: "Operations with vulcan encryption keys",
		Subcommands: []cli.Command{
			{
				Name:   "new_key",
				Usage:  "Generate new seal key",
				Action: cmd.generateKeyAction,
				Flags: []cli.Flag{
					cli.StringFlag{Name: "file, f", Usage: "File to write to"},
				},
			},
			{
				Name:   "seal_keypair",
				Usage:  "Seal key pair",
				Action: cmd.sealKeyPairAction,
				Flags: []cli.Flag{
					cli.StringFlag{Name: "file, f", Usage: "File to write to"},
					cli.StringFlag{Name: "sealKey", Usage: "Seal key - used to encrypt and seal certificate and private key"},
					cli.StringFlag{Name: "privateKey", Usage: "Path to a private key"},
					cli.StringFlag{Name: "cert", Usage: "Path to a certificate"},
				},
			},
		},
	}
}

func (cmd *Command) generateKeyAction(c *cli.Context) error {
	key, err := secret.NewKeyString()
	if err != nil {
		return fmt.Errorf("unable to generate key: %v", err)
	}
	stream, closer, err := getStream(c)
	if err != nil {
		return err
	}
	if closer != nil {
		defer closer.Close()
	}
	_, err = stream.Write([]byte(key))
	if err != nil {
		return fmt.Errorf("failed writing to output stream, error %s", err)
	}
	return nil
}

func (cmd *Command) sealKeyPairAction(c *cli.Context) error {
	// Read the key and get a box
	box, err := readBox(c.String("sealKey"))
	if err != nil {
		return err
	}

	// Read keyPairificate
	stream, closer, err := getStream(c)
	if err != nil {
		return err
	}
	if closer != nil {
		defer closer.Close()
	}

	keyPair, err := readKeyPair(c.String("cert"), c.String("privateKey"))
	if err != nil {
		return fmt.Errorf("failed to read key pair: %s", err)
	}

	bytes, err := secret.SealKeyPairToJSON(box, keyPair)
	if err != nil {
		return fmt.Errorf("failed to seal key pair: %s", err)
	}

	_, err = stream.Write(bytes)
	if err != nil {
		return fmt.Errorf("failed writing to output stream, error %s", err)
	}
	return nil
}

func getStream(c *cli.Context) (io.Writer, io.Closer, error) {
	if c.String("file") != "" {
		file, err := os.OpenFile(c.String("file"), os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to open file %s, error: %s", c.String("file"), err)
		}
		return file, file, nil
	}
	return os.Stdout, nil, nil
}

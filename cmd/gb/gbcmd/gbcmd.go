package gbcmd

import (
	"context"
	"ghostbb.io/gb/cmd/gb/internal/cmd"
	gbcmd "ghostbb.io/gb/os/gb_cmd"
)

// Command manages the CLI command of `gb`.
// This struct can be globally accessible and extended with custom struct.
type Command struct {
	*gbcmd.Command
}

// GetCommand retrieves and returns the root command of CLI `gb`.
func GetCommand(ctx context.Context) (*Command, error) {
	root, err := gbcmd.NewFromObject(cmd.GB)
	if err != nil {
		panic(err)
	}
	err = root.AddObject(
		cmd.Version,
		cmd.Pack,
		cmd.Install,
	)
	if err != nil {
		return nil, err
	}
	command := &Command{
		root,
	}
	return command, nil
}

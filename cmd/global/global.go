package global

import "github.com/alecthomas/kong"

type Flags struct {
	Debug   bool             `short:"d" help:"Enable debug mode"`
	Version kong.VersionFlag `short:"v" help:"Print version and quit"`
}

package files

import (
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/internal/system"
)

var rm bool

var Command = &cobra.Command{
	Use:     "files",
	Aliases: []string{"file"},
	Run:     system.Help,
}

func init() {
	Command.AddCommand(&cobra.Command{Use: "list", Aliases: []string{"ls"}, Run: list})
	Command.AddCommand(&cobra.Command{Use: "remove", Aliases: []string{"rm"}, Run: remove})
	compressCommand := &cobra.Command{Use: "compress", Run: compress}
	compressCommand.PersistentFlags().BoolVar(&rm, "rm", false, "Remove after compression ends.")
	Command.AddCommand(compressCommand)

	decompressCommand := &cobra.Command{Use: "decompress", Run: decompress}
	decompressCommand.PersistentFlags().BoolVar(&rm, "rm", false, "Remove after decompression ends.")
	Command.AddCommand(decompressCommand)

	backupCommand := &cobra.Command{Use: "backup", Run: backup}
	backupCommand.PersistentFlags().BoolVar(&rm, "rm", false, "Remove after decompression ends.")
	Command.AddCommand(backupCommand)
}

func check(err error) {
	if err != nil {
		slog.Error(err.Error())
		return
	}
}

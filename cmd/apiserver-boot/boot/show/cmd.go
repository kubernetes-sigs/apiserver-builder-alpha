package show

import "github.com/spf13/cobra"

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Command group for show the installed aggregated apiserver.",
	Long:  `Command group for show the installed aggregated apiserver.`,
	Example: `
# Show the current status for foo resource.
apiserver-boot show resource foo
`,
	Run: RunShow,
}

func AddShow(cmd *cobra.Command) {
	cmd.AddCommand(showCmd)
	AddShowResource(showCmd)
}

func RunShow(cmd *cobra.Command, args []string) {
	cmd.Help()
}

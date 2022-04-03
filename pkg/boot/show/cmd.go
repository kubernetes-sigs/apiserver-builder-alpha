package show

import (
	"os"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var streams = genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}
var clientFactory *genericclioptions.ConfigFlags

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
	AddApiserver(showCmd)
}

func RunShow(cmd *cobra.Command, args []string) {
	cmd.Help()
}

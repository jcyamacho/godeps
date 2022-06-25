package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/jcyamacho/godeps/internal/mermaid"
	"github.com/jcyamacho/godeps/internal/module"

	"github.com/spf13/cobra"
)

var skipModules []string
var skipIndirectDeps bool

func init() {
	flowchartCmd.Flags().StringSliceVarP(&skipModules, "skip", "s", []string{}, "skip modules")
	flowchartCmd.Flags().BoolVarP(&skipIndirectDeps, "skip-indirect", "i", true, "skip indirect dependencies")
	rootCmd.AddCommand(flowchartCmd)
}

func isSkiped(mod *module.Module) bool {
	for _, skip := range skipModules {
		if strings.Contains(mod.Name, skip) {
			return true
		}
	}
	return false
}

var flowchartCmd = &cobra.Command{
	Use:   "flowchart [path]",
	Short: "Print a mermaid flowchart dependency graph",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		graph, err := module.ScanDir(args[0], skipIndirectDeps)
		if err != nil {
			log.Fatal(err)
		}

		if len(skipModules) > 0 {
			graph = graph.Prune(isSkiped)
		}

		flowchart := mermaid.Flowchart(graph)
		fmt.Println(flowchart)
	},
}

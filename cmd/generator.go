package cmd

import (
	"fmt"
	"os"

	"github.com/pasha-codefresh/argo-cd-contrib-insights-generator/pkg"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate Argo CD contrib meeting insights",
	RunE: func(cmd *cobra.Command, args []string) error {

		fmt.Println("- **Weekly GitHub Stats** (Extracted with [insight-generator](https://github.com/crenshaw-dev/argo-cd-contrib-insights-generator))")
		if err := printProjectStats("Argo CD", pkg.ArgoCDRepo); err != nil {
			return err
		}
		if err := printProjectStats("Argo Rollouts", pkg.ArgoRolloutsRepo); err != nil {
			return err
		}

		fmt.Print("- **Weekly LFX Insights Stats**\n")

		topReviewersStats, err := pkg.NewTopReviewersStatsGenerator().Generate()
		if err != nil {
			return err
		}
		fmt.Println(topReviewersStats)

		topMergersStats, err := pkg.NewTopMergersStatsGenerator().Generate()
		if err != nil {
			return err
		}
		fmt.Println(topMergersStats)

		return nil
	},
}

func printProjectStats(name, repo string) error {
	fmt.Printf("  - **%s:**\n", name)

	createdIssuesStats, err := pkg.NewCreatedIssuesStatsGenerator(repo).Generate()
	if err != nil {
		return err
	}
	fmt.Printf("    - %s\n", createdIssuesStats)

	createdPRsStats, err := pkg.NewCreatedPRsStatsGenerator(repo).Generate()
	if err != nil {
		return err
	}
	fmt.Printf("    - %s\n", createdPRsStats)

	staleIssuesStats, err := pkg.NewStaleIssuesStatsGenerator(repo).Generate()
	if err != nil {
		return err
	}
	fmt.Printf("    - %s\n", staleIssuesStats)
	return nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

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

		fmt.Println(" - Argo CD:")
		createdIssuesStats, link, err := pkg.NewCreatedIssuesStatsGenerator().Generate()
		if err != nil {
			return err
		}
		fmt.Println(createdIssuesStats)
		fmt.Println(link)

		createdPRsStats, link, err := pkg.NewCreatedPRsStatsGenerator().Generate()
		if err != nil {
			return err
		}
		fmt.Println(createdPRsStats)
		fmt.Println(link)

		staleIssuesStats, link, err := pkg.NewStaleIssuesStatsGenerator().Generate()
		if err != nil {
			return err
		}
		fmt.Println(staleIssuesStats)
		fmt.Println(link)

		fmt.Println(" - Top Reviewers:")

		topReviewersStats, link, err := pkg.NewTopReviewersStatsGenerator().Generate()
		if err != nil {
			return err
		}
		fmt.Println(topReviewersStats)
		fmt.Println(link)

		fmt.Println(" - Top Mergers:")

		topMergersStats, link, err := pkg.NewTopMergersStatsGenerator().Generate()
		if err != nil {
			return err
		}
		fmt.Println(topMergersStats)
		fmt.Println(link)

		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

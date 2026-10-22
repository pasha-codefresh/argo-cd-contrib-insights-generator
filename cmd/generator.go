package cmd

import (
	"fmt"
	"os"

	"github.com/pasha-codefresh/argo-cd-contrib-insights-generator/pkg"
	"github.com/spf13/cobra"
)

func newRootCmd() *cobra.Command {
	var repo string
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate Argo CD contrib meeting insights",
		RunE: func(cmd *cobra.Command, args []string) error {

			if err := pkg.ValidateRepo(repo); err != nil {
				return err
			}

			fmt.Printf(" - %s:\n", repo)
			createdIssuesStats, link, err := pkg.NewCreatedIssuesStatsGenerator(repo).Generate()
			if err != nil {
				return err
			}
			fmt.Println(createdIssuesStats)
			fmt.Println(link)

			createdPRsStats, link, err := pkg.NewCreatedPRsStatsGenerator(repo).Generate()
			if err != nil {
				return err
			}
			fmt.Println(createdPRsStats)
			fmt.Println(link)

			staleIssuesStats, link, err := pkg.NewStaleIssuesStatsGenerator(repo).Generate()
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
	cmd.Flags().StringVar(&repo, "repo", pkg.DefaultRepo, "GitHub repository (owner/name) to collect issue and PR stats from")
	return cmd
}

func Execute() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

package main

import (
	"os"
	"strconv"

	"github.com/atsman/gsh/pkg/client"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func main() {
	githubClient := client.New()

	var searchCmd = &cobra.Command{
		Use:   "search [query]",
		Short: "Search github repos by query",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			searchHandler(githubClient, args[0])
		},
	}

	var trendingCmd = &cobra.Command{
		Use:   "trending -s [since]",
		Short: "Search trending repos",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			since := cmd.Flags().Lookup("since")
			trendingHandler(githubClient, since.Value.String())
		},
	}

	trendingCmd.Flags().StringP("since", "s", "mountly", "-s [since]")

	var rootCmd = &cobra.Command{Use: "gsh"}
	rootCmd.AddCommand(searchCmd, trendingCmd)
	rootCmd.Execute()
}

func searchHandler(gClient *client.GithubClient, query string) {
	sr, err := gClient.Search(query)
	if err != nil {
		panic(err)
	}
	printResults(sr)
}

func trendingHandler(gClient *client.GithubClient, since string) {
	sr, err := gClient.Trending(since)
	if err != nil {
		panic(err)
	}
	printResults(sr)
}

func shortString(str string, maxLength int) string {
	if len(str) < maxLength {
		return str
	}
	return str[:maxLength] + "..."
}

func printResults(sr *client.SearchResult) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"#", "Name", "Stars", "Description", "Url"})

	for i, repo := range sr.Items {
		table.Append([]string{
			strconv.Itoa(i + 1),
			shortString(repo.Name, 20),
			strconv.Itoa(repo.Stars),
			shortString(repo.Description, 100),
			repo.URL,
		})
	}

	table.Render()
}

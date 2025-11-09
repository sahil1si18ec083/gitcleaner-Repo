package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v59/github"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

const cliVersion = "v1.0.0"

var rootCmd = &cobra.Command{
	Use:   "gitcleaner",
	Short: "A tool to clean up git repositories",
	Long:  `gitcleaner is a CLI tool that helps you clean up your git repositories by removing unnecessary files and optimizing the repository size.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gitcleaner: Use 'gitcleaner list' or 'gitcleaner delete' commands.")

	},
}
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all your GitHub repositories",
	Long:  "Fetches and displays all repositories from your GitHub account.",
	Run: func(cmd *cobra.Command, args []string) {
		repos, _, _ := getRepoHelper()
		if len(repos) == 0 {
			fmt.Println("No repositories found.")
			return
		}

		fmt.Println("Your repositories:")
		for _, r := range repos {
			fmt.Println("-", r.GetName())
		}
	},
}
var (
	deleteAllFlag    bool
	deleteOnlyFlag   string
	deleteExceptFlag string
	dryRunFlag       bool
)

func contains(except []string, r string) bool {
	for _, e := range except {
		if e == r {
			return true
		}
	}
	return false
}
func getRepoHelper() ([]*github.Repository, *github.Client, context.Context) {
	fmt.Println("Fetching your repositories...")

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		fmt.Println("❌ Please set your GITHUB_TOKEN environment variable first.")
		return []*github.Repository{}, nil, context.Background()
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	opt := &github.RepositoryListByAuthenticatedUserOptions{
		ListOptions: github.ListOptions{PerPage: 50},
	}

	repos, _, err := client.Repositories.ListByAuthenticatedUser(ctx, opt)
	if err != nil {
		fmt.Println("❌ Error while fetching repos:", err)
		return []*github.Repository{}, client, ctx
	}

	return repos, client, ctx

}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete repositories in bulk",
	Long:  "Delete your GitHub repositories with various options (all, only, except, dry-run).",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Deleting repositories...")
		repos, client, ctx := getRepoHelper()

		repoNames := []string{}
		for _, r := range repos {
			repoNames = append(repoNames, r.GetName())
			fmt.Println("-", r.GetName())
		}
		toDelete := []string{}
		if deleteAllFlag {
			toDelete = repoNames

		} else if deleteOnlyFlag != "" {
			toDelete = strings.Split(deleteOnlyFlag, ",")

		} else if deleteExceptFlag != "" {
			except := strings.Split(deleteExceptFlag, ",")
			for _, r := range repoNames {
				if !contains(except, r) {
					toDelete = append(toDelete, r)
				}

			}
		} else {
			fmt.Println("⚠️ Please specify a flag: --all, --only, or --except")
			return
		}
		if dryRunFlag {
			fmt.Println("🧪 Dry run mode — repositories that *would* be deleted:")
			for _, r := range toDelete {
				fmt.Println(" -", r)
				return
			}
		}

		fmt.Println("⚠️ The following repositories will be deleted:")
		for _, r := range toDelete {
			fmt.Println(" -", r)
		}
		fmt.Print("Are you sure you want to continue? (y/n): ")
		var confirmation string
		fmt.Scanln(&confirmation)
		switch confirmation {
		case "y", "Y":
			user, _, err := client.Users.Get(ctx, "")
			if err != nil {
				fmt.Println("Error fetching user:", err)
				return
			}
			username := user.GetLogin()
			for _, r := range toDelete {
				_, err := client.Repositories.Delete(ctx, username, r)
				if err != nil {
					fmt.Println(err)
					// fmt.Printf("%v", err)
				}
			}
			fmt.Println("✅ Repositories deleted successfully.")

		case "n", "N":
			fmt.Println("❎ Deletion cancelled.")
			return
		default:
			fmt.Println("❌ Invalid input. Please enter 'y' or 'n'.")

		}

	},
}
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the current version of gitcleaner",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("gitcleaner %s\n", cliVersion)
	},
}

func main() {
	deleteCmd.Flags().BoolVarP(&deleteAllFlag, "all", "a", false, "Delete all repositories")
	deleteCmd.Flags().StringVarP(&deleteOnlyFlag, "only", "o", "", "Delete only specific repositories (comma-separated)")
	deleteCmd.Flags().StringVarP(&deleteExceptFlag, "except", "e", "", "Delete all except these repositories (comma-separated)")
	deleteCmd.Flags().BoolVarP(&dryRunFlag, "dry-run", "d", false, "Preview deletions without actually deleting")
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(deleteCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

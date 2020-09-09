package cmd

import (
	"log"

	"github.com/Hamaiz/go-rest-eg/database"
	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "check deletes records from database",
	Long: `check looks for all the accounts that expired when people didnt use the token
		and deletes them.
		It is like a corn job. It runs every half hour.
		`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("check called")
		database.DeleteAccount()
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}

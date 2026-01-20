package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/nicholasjconn/agent-layer/internal/doctor"
)

func newDoctorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Report missing secrets, disabled servers, and common misconfigurations",
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := getwd()
			if err != nil {
				return err
			}

			fmt.Printf("🏥 Checking Agent Layer health in %s...\n\n", root)

			var allResults []doctor.Result

			// 1. Check Structure
			allResults = append(allResults, doctor.CheckStructure(root)...)

			// 2. Check Config
			configResults, cfg := doctor.CheckConfig(root)
			allResults = append(allResults, configResults...)

			if cfg != nil {
				// 3. Check Secrets
				allResults = append(allResults, doctor.CheckSecrets(cfg)...)

				// 4. Check Agents
				allResults = append(allResults, doctor.CheckAgents(cfg)...)
			}

			hasFail := false
			for _, r := range allResults {
				printResult(r)
				if r.Status == doctor.StatusFail {
					hasFail = true
				}
			}

			fmt.Println()
			if hasFail {
				color.Red("❌ Some checks failed. Please address the issues above.")
				return fmt.Errorf("doctor checks failed")
			} else {
				color.Green("✅ All systems go! Agent Layer is ready.")
			}

			return nil
		},
	}
}

func printResult(r doctor.Result) {
	var status string
	switch r.Status {
	case doctor.StatusOK:
		status = color.GreenString("[OK]  ")
	case doctor.StatusWarn:
		status = color.YellowString("[WARN]")
	case doctor.StatusFail:
		status = color.RedString("[FAIL]")
	}

	fmt.Printf("%s %-10s %s\n", status, r.CheckName, r.Message)
	if r.Recommendation != "" {
		fmt.Printf("       💡 %s\n", r.Recommendation)
	}
}

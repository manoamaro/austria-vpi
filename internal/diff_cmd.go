package internal

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var (
	fromYear  int
	fromMonth int
	toYear    int
	toMonth   int
)

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Calculate the VPI difference between two periods",
	Long:  "Calculate the VPI difference between two periods",
	RunE: func(cmd *cobra.Command, args []string) error {
		if vpiFrom, vpiTo, diff, err := calculateDiff(fromYear, fromMonth, toYear, toMonth); err != nil {
			return err
		} else {
			fmt.Printf("VPI from %d-%02d: %.1f\n", fromYear, fromMonth, vpiFrom)
			fmt.Printf("VPI to %d-%02d: %.1f\n", toYear, toMonth, vpiTo)
			fmt.Printf("The consumer price index has changed by %.1f\n", diff)
			return nil
		}
	},
}

func DiffCmd() *cobra.Command {
	diffCmd.PersistentFlags().IntVar(&fromYear, "from-year", 2015, "The year to start from")
	diffCmd.PersistentFlags().IntVar(&fromMonth, "from-month", 1, "The month to start from")
	diffCmd.PersistentFlags().IntVar(&toYear, "to-year", time.Now().Year(), "The year to end at")
	diffCmd.PersistentFlags().IntVar(&toMonth, "to-month", int(time.Now().Month())-1, "The month to end at")
	return diffCmd
}

func calculateDiff(fromYear, fromMonth, toYear, toMonth int) (float64, float64, float64, error) {
	if csvRecords, err := getRawCSV(); err != nil {
		return 0, 0, 0, err
	} else if records, err := parseRecords(filterForVPI0(csvRecords)); err != nil {
		return 0, 0, 0, err
	} else {
		recordsMap := make(map[int]map[int]float64)
		for _, record := range records {
			if _, ok := recordsMap[record.Year]; !ok {
				recordsMap[record.Year] = make(map[int]float64)
			}
			recordsMap[record.Year][record.Month] = record.VPI
		}
		vpiFrom := recordsMap[fromYear][fromMonth]
		vpiTo := recordsMap[toYear][toMonth]
		vpiDiff := vpiTo - vpiFrom
		vpiDiffPercent := vpiDiff / vpiFrom * 100
		return vpiFrom, vpiTo, vpiDiffPercent, nil
	}
}

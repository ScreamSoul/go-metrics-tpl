package versions

import "fmt"

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func PrintBuildInfo() {
	fmt.Printf("Build version: %s\n\r", buildVersion)
	fmt.Printf("Build date: %s\n\r", buildDate)
	fmt.Printf("Build commit: %s\n\r", buildCommit)
}

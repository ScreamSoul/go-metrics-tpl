// Correctly prints default build information
package versions_test

import (
	"testing"

	"github.com/screamsoul/go-metrics-tpl/internal/versions"
)

func TestPrintBuildInfo_DefaultValues(t *testing.T) {
	// Call the function under test
	versions.PrintBuildInfo()
}

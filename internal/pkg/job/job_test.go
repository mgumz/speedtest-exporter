package job

import (
	"testing"
)

func Test_NormalizeSpeedtestErrorMsg(t *testing.T) {

	fixtures := []struct {
		msg        string
		normalized string
	}{
		// $> speedtest-go --json --custom-url=https://invalid.example.com
		// `Fatal: latency: --, err: server connect timeout`
		{"Fatal: latency: --, err: server connect timeout",
			"fatal: latency: --- err: server connect timeout"},
	}

	for i := range fixtures {
		normalized := normalizeSpeedtestErrorMsg(fixtures[i].msg)
		if normalized != fixtures[i].normalized {
			t.Fatalf("expected %q for %q, got %q", fixtures[i].normalized, fixtures[i].msg, normalized)
		}
	}
}

package job

import (
	"strings"
	"testing"
)

func Test_JobFileParsing(t *testing.T) {
	fixtures := []struct {
		File     string
		Expected []*Job
	}{}

	for i := range fixtures {
		r := strings.NewReader(fixtures[i].File)
		jobs, err := ParseJobs(r, "speedtest-go")
		if err != nil {
			t.Fatalf("error parsing: %s\n%s", err, fixtures[i].File)
		}
		t.Logf("jobs: %d", len(jobs))
	}
}

func Test_ParseSpeedtestArgs(t *testing.T) {

	fixtures := []struct {
		Args     string
		Expected []string
	}{
		{"a", []string{"a"}},
		{"a b", []string{"a", "b"}},
		{" a  b", []string{"a", "b"}},
	}

	for i := range fixtures {
		args, _ := parseSpeedtestArgs(fixtures[i].Args)
		if len(args) != len(fixtures[i].Expected) {
			t.Fatalf("error parsing speedtest-go args: expected %q, got %q",
				fixtures[i].Expected, args)
		}
	}

}

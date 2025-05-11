package job

import (
	"testing"

	"github.com/mgumz/speedtest-exporter/internal/pkg/timeshift"
)

func Test_ParseScheduleTimeshift(t *testing.T) {

	type expT struct {
		err      error
		schedule string
		tsMode   timeshift.Mode
		tsSpec   string
	}

	fixtures := []struct {
		spec   string
		expect expT
	}{
		{"@every 1h", expT{nil, "@every 1h", timeshift.None, ""}},
		{"trailing space  ", expT{nil, "trailing space", timeshift.None, ""}},
		{"  prefix space", expT{nil, "prefix space", timeshift.None, ""}},

		// disabled timeshift.RandomDeviation
		//{"@every 1h ±1h", expT{nil, "@every 1h", timeshift.RandomDeviation, "1h"}},
		//{"@every 1h ±1h ", expT{nil, "@every 1h", timeshift.RandomDeviation, "1h"}},
		//{" @every 1h ±1h ", expT{nil, "@every 1h", timeshift.RandomDeviation, "1h"}},
		//{"@every 1h ± 1h", expT{nil, "@every 1h", timeshift.RandomDeviation, "1h"}},
		//{"@every 1h ± 1h ", expT{nil, "@every 1h", timeshift.RandomDeviation, "1h"}},
		//{"@every 1h ± +1h", expT{nil, "@every 1h", timeshift.RandomDeviation, "+1h"}},

		{"@every 1h ~1h", expT{nil, "@every 1h", timeshift.RandomDelay, "1h"}},
		{"@every 1h ~1h ", expT{nil, "@every 1h", timeshift.RandomDelay, "1h"}},
		{" @every 1h ~1h ", expT{nil, "@every 1h", timeshift.RandomDelay, "1h"}},
	}

	for i, f := range fixtures {
		pSched, pTsMode, pTsSpec, err := parseSchedule(f.spec)
		if f.expect.err != err {
			t.Fatalf("fixture %d: expected error: %v, got %q", i, f.expect.err, err)
		}
		if f.expect.schedule != pSched {
			t.Fatalf("fixture %d: expected schedule: %q, got %q", i, f.expect.schedule, pSched)
		}
		if f.expect.tsMode != pTsMode {
			t.Fatalf("fixture %d: expected tsMode: %d, got %d", i, f.expect.tsMode, pTsMode)
		}
		if f.expect.tsSpec != pTsSpec {
			t.Fatalf("fixture %d: expected tsSpec: %q, got %q", i, f.expect.tsSpec, pTsSpec)
		}
	}

}

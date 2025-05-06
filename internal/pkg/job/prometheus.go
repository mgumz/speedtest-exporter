package job

import (
	"fmt"
	"maps"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/mgumz/speedtest-exporter/internal/pkg/prometheus"
	"github.com/mgumz/speedtest-exporter/internal/pkg/speedtest"
)

const (
	integerBase int = 10
)

var (
	speedTestMM = []prometheus.MetricMeta{
		{Name: "speedtest_runs_total", MType: "counter", Help: "number of speedtest runs"},
		{Name: "speedtest_jobfile_parsed_total", MType: "counter", Help: "number of jobfile related parse runs"},
		{Name: "speedtest_report_duration_seconds", MType: "gauge", Help: "duration of last speedtest run (in seconds)"},
		{Name: "speedtest_report_count_hubs", MType: "gauge", Help: "number of hops visited in the last speedtest run"},
	}
)

// ServeHTTP writes prometheus styled metrics about the last executed `speedtest-go`
// run, see https://prometheus.io/docs/instrumenting/exposition_formats/#line-format
//
// NOTE: at the moment, no use of github.com/prometheus/client_golang/prometheus
// because overhead in size and complexity. once speedtest-exporter requires features
// like push-gateway-export or graphite export or the like, we switch.
func (c *Collector) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	c.mu.Lock()
	defer c.mu.Unlock()

	prometheus.WriteMeta(w, speedTestMM)
	speedtest.WriteMetricsHelpType(w)

	fmt.Fprintf(w, "speedtest_jobfile_parsed_total{reason=%q} %d\n",
		"failed", c.metrics.jobFileWatch.failedTotal)
	fmt.Fprintf(w, "speedtest_jobfile_parsed_total{reason=%q} %d\n",
		"unchanged", c.metrics.jobFileWatch.unchangedTotal)
	fmt.Fprintf(w, "speedtest_jobfile_parsed_total{reason=%q} %d\n",
		"changed", c.metrics.jobFileWatch.changedTotal)

	if len(c.jobs) == 0 {
		fmt.Fprintln(w, "# no speedtest jobs defined (yet).")
		return
	}

	fmt.Fprintf(w, "# %d speedtest jobs defined\n", len(c.jobs))

	for _, job := range c.jobs {

		if !job.DataAvailable() {
			continue
		}

		// the original job.Report might be changed in the background by a
		// successful run of speedtest. copy (pointer to) the report to have
		// something safe to work on
		result := job.Result
		ts := job.Launched.UTC()
		d := job.Duration

		labels := map[string]string{} // FIXME: provide some basic labels: result.Labels()
		labels["speedtest_exporter_job"] = job.Label
		tsMs := ts.UnixNano() / int64(time.Millisecond)

		errMsg := ""
		if result.ErrorMsg != "" {
			errMsg = fmt.Sprintf(" # (err: %q)", result.ErrorMsg)
		}
		fmt.Fprintf(w, "# speedtest run %s: %s -- %s%s\n", job.Label, ts.Format(time.RFC3339Nano), job.CmdLine, errMsg)

		l := labels2Prom(labels)

		for k, v := range job.Runs {
			fmt.Fprintf(w, "speedtest_runs_total{%s%s} %d %d\n",
				l, fmt.Sprintf(",error=%q", k), v, tsMs)
		}

		fmt.Fprintf(w, "speedtest_duration_seconds{%s} %f %d\n",
			l, float64(d)/float64(time.Second), tsMs)

		if len(job.Result.Server) == 0 {
			continue
		}

		for _, server := range result.Server {

			slabels := server.Labels()
			maps.Copy(slabels, labels)

			server.WriteMetrics(w, labels2Prom(labels), tsMs)
		}
	}

}

func labels2Prom(labels map[string]string) string {
	sl := make(sort.StringSlice, 0, len(labels))
	for k, v := range labels {
		sl = append(sl, fmt.Sprintf("%s=%q", k, v))
	}
	sl.Sort()
	return strings.Join(sl, ",")
}

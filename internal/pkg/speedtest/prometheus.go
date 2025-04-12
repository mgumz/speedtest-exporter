package speedtest

import (
	"fmt"
	"io"
	"time"

	"github.com/mgumz/speedtest-exporter/internal/pkg/prometheus"
)

func (server *Server) WriteMetrics(w io.Writer, labels string, ts int64) {

	fmt.Fprintf(w, "speedtest_dl_speed_rate{%s} %d %d\n", labels, int64(server.DlSpeed), ts)
	fmt.Fprintf(w, "speedtest_ul_speed_rate{%s} %d %d\n", labels, int64(server.UlSpeed), ts)
	fmt.Fprintf(w, "speedtest_jitter_seconds{%s} %f %d\n", labels, nanoToSeconds(server.Jitter), ts)
	fmt.Fprintf(w, "speedtest_latency_seconds{%s} %f %d\n", labels, nanoToSeconds(server.Latency), ts)
	fmt.Fprintf(w, "speedtest_max_latency_seconds{%s} %f %d\n", labels, nanoToSeconds(server.MaxLatency), ts)
	fmt.Fprintf(w, "speedtest_min_latency_seconds{%s} %f %d\n", labels, nanoToSeconds(server.MinLatency), ts)
	fmt.Fprintf(w, "speedtest_test_duration_ping_seconds{%s} %f %d\n", labels, nanoToSeconds(server.TestDuration.Ping), ts)
	fmt.Fprintf(w, "speedtest_test_duration_download_seconds{%s} %f %d\n", labels, nanoToSeconds(server.TestDuration.Download), ts)
	fmt.Fprintf(w, "speedtest_test_duration_upload_seconds{%s} %f %d\n", labels, nanoToSeconds(server.TestDuration.Upload), ts)
	fmt.Fprintf(w, "speedtest_test_duration_total_seconds{%s} %f %d\n", labels, nanoToSeconds(server.TestDuration.Total), ts)
	fmt.Fprintf(w, "speedtest_packetloss_sent_packets{%s} %d %d\n", labels, server.PacketLoss.Sent, ts)
	fmt.Fprintf(w, "speedtest_packetloss_dup_packets{%s} %d %d\n", labels, server.PacketLoss.Dup, ts)
	fmt.Fprintf(w, "speedtest_packetloss_max_packets{%s} %d %d\n", labels, server.PacketLoss.Max, ts)
}

func nanoToSeconds(v int64) float64 {
	return time.Duration(v).Truncate(time.Microsecond).Seconds()
}

var (
	speedTestMM = []prometheus.MetricMeta{
		{Name: "speedtest_runs_total", MType: "counter", Help: "tracks how many speedtest runs were executed"},
		{Name: "speedtest_duration_seconds", MType: "gauge", Help: "duration of last performed speedtest run"},
		{Name: "speedtest_dl_speed_rate", MType: "gauge", Help: "downloaded bytes per second"},
		{Name: "speedtest_ul_speed_rate", MType: "gauge", Help: "uploaded bytes per second"},
		{Name: "speedtest_jitter_seconds", MType: "gauge", Help: "jitter"},
		{Name: "speedtest_max_latency_seconds", MType: "gauge", Help: "maximum latency"},
		{Name: "speedtest_min_latency_seconds", MType: "gauge", Help: "minimum latency"},
		{Name: "speedtest_test_duration_ping_seconds", MType: "gauge", Help: "test duration for ping-phase"},
		{Name: "speedtest_test_duration_download_seconds", MType: "gauge", Help: "test duration for download-phase"},
		{Name: "speedtest_test_duration_upload_seconds", MType: "gauge", Help: "test duration for upload-phase"},
		{Name: "speedtest_test_duration_total_seconds", MType: "gauge", Help: "total test duration"},
		{Name: "speedtest_packetloss_sent_packets", MType: "gauge", Help: "amount of sent packets"},
		{Name: "speedtest_packetloss_dup_packets", MType: "gauge", Help: "amount of duplicated packets"},
		{Name: "speedtest_packetloss_max_packets", MType: "gauge", Help: "maximum of packet loss"},
	}
)

func WriteMetricsHelpType(w io.Writer) {
	prometheus.WriteMeta(w, speedTestMM)
}

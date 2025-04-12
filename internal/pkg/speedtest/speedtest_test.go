package speedtest

import (
	"strings"
	"testing"
)

func Test_SpeedtestReportDecoding(t *testing.T) {
	body := `
{
    "timestamp": "2025-04-12 08:29:00.958",
    "user_info": {
        "IP": "192.0.2.1",
        "Lat": "0.0",
        "Lon": "0.0",
        "Isp": "ACME Example"
    },
    "servers": [
        {
            "url": "http://example.com:8080/speedtest/upload.php",
            "lat": "0.0",
            "lon": "0.0",
            "name": "Example City",
            "country": "Example Country",
            "sponsor": "ACME Example",
            "id": "00001",
            "host": "",
            "distance": 1.23,
            "latency": 7521392,
            "max_latency": 9164731,
            "min_latency": 5934208,
            "jitter": 1015667,
            "dl_speed": 31191357.833144464,
            "ul_speed": 10086856.246238846,
            "test_duration": {
                "ping": 2346185184,
                "download": 10401042321,
                "upload": 11402219402,
                "total": 24149446907
            },
            "packet_loss": { "sent": 0, "dup": 0, "max": 0 }
        }
    ]
}
`

	result := &Result{}
	if err := result.Decode(strings.NewReader(body)); err != nil {
		t.Fatalf("error decoding: %s\n%s", err, body)
	}

	if result.Server[0].ID != "0001" {
		t.Fatalf("error parsing speedtest report: expected %q, got %q\n%v",
			"dst.example.com",
			result.Server[0].ID,
			result)
	}
}

func Test_SpeedtestEmptyServers(t *testing.T) {
	body := `
{
	"timestamp": "2025-04-12 08:29:00.958",
	"user_info": {
	    "IP": "192.0.2.1",
	    "Lat": "0.0",
	    "Lon": "0.0",
	    "Isp": "ACME Example"
	},
	"servers": [
	]
}
`

	result := &Result{}
	if err := result.Decode(strings.NewReader(body)); err != nil {
		t.Fatalf("error decoding: %s\n%s", err, body)
	}

	if len(result.Server) != 0 {
		t.Fatalf("error: expected [] servers, got %d", len(result.Server))
	}
}

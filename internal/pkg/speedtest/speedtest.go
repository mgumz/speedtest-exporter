package speedtest

import (
	"encoding/json"
	"io"
)

const (
	integerBase int = 10
)

type Result struct {
	UserInfo UserInfo `json:"user_info"`
	Server   []Server `json:"servers"`
	ErrorMsg string   // carrying the error message of speedtest-go
}

type UserInfo struct {
	IP  string `json:"IP"`
	Lat string `json:"Lat"`
	Lon string `json:"Lon"`
	Isp string `json:"Isp"`
}

type Server struct {
	URL          string  `json:"url"`
	Lat          string  `json:"lat"`
	Lon          string  `json:"lon"`
	Name         string  `json:"name"`
	Country      string  `json:"country"`
	Sponsor      string  `json:"sponsor"`
	ID           string  `json:"id"`
	Host         string  `json:"host"`
	Distance     float64 `json:"distance"`
	Latency      int64   `json:"latency"`
	MaxLatency   int64   `json:"max_latency"`
	MinLatency   int64   `json:"min_latency"`
	Jitter       int64   `json:"jitter"`
	DlSpeed      float64 `json:"dl_speed"`
	UlSpeed      float64 `json:"ul_speed"`
	TestDuration struct {
		Ping     int64 `json:"ping"`
		Download int64 `json:"download"`
		Upload   int64 `json:"upload"`
		Total    int64 `json:"total"`
	} `json:"test_duration"`
	PacketLoss struct {
		Sent int64 `json:"sent"`
		Dup  int64 `json:"dup"`
		Max  int64 `json:"max"`
	} `json:"packet_loss"`
}

func (result *Result) Decode(r io.Reader) error {
	dec := json.NewDecoder(r)
	res := Result{}
	if err := dec.Decode(&res); err != nil {
		return err
	}
	*result = res
	return nil
}

func (result *Result) Empty() bool {
	return len(result.Server) == 0
}

func (server *Server) Labels() map[string]string {
	return map[string]string{
		"server_id":  server.ID,
		"server_url": server.URL,
		"sponsor":    server.Sponsor,
	}
}

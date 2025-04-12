package main

import (
	"fmt"
	"net/http"
)

const rootHtml = `<!doctype html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<title>speedtest-exporter</title>
	<style>
		* { font-family: sans-serif; }
		body { margin: auto 20%; }
		h1 { margin-top: 3em; }
	</style>
</head>
<body>
	<h1>speedtest-exporter</h1>
	<p><a href="https://github.com/mgumz/speedtest-exporter">https://github.com/mgumz/speedtest-exporter<a><br>
	see <a href="/metrics">/metrics</a>.
	</p>
</body>
`

func handleHealth(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "OK")
}

func handleRoot(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, rootHtml)
}

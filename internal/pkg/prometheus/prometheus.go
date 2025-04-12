package prometheus

import (
	"fmt"
	"io"
)

type MetricMeta struct {
	Name  string
	MType string
	Help  string
}

func WriteMeta(w io.Writer, metrics []MetricMeta) {

	for _, m := range metrics {
		fmt.Fprintln(w, "# HELP", m.Name, m.Help)
		fmt.Fprintln(w, "# TYPE", m.Name, m.MType)
	}
}

package job

import (
	"crypto/sha256"
	"io"
	"os"
)

// speedestJobFile definition
//
// # comments, ignore everything after #
// ^space*$ - empty lines
// <label> -- <schedule> -- <speedtest-go-flags>

func ParseJobFile(filename, speedtest string) (jobs Jobs, cksum []byte, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return []*Job{}, []byte{}, err
	}
	defer f.Close()

	h := sha256.New()
	r := io.TeeReader(f, h)

	jobs, err = ParseJobs(r, speedtest)

	return jobs, h.Sum(nil), err
}

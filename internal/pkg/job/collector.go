package job

import (
	"sync"
)

type Collector struct {
	jobs []JobMeta

	mu sync.Mutex

	metrics struct {
		jobFileWatch struct {
			changedTotal   int64
			unchangedTotal int64
			failedTotal    int64
		}
	}

	opts collectorOpts
}

type collectorOpts struct {
	doRenderDeprecatedMetrics bool
}

func NewCollector() *Collector {
	return new(Collector)
}

func (c *Collector) SetRenderDeprecatedMetrics(f bool) *Collector {
	c.opts.doRenderDeprecatedMetrics = f
	return c
}

func (c *Collector) RemoveJob(label string) bool {

	c.mu.Lock()
	defer c.mu.Unlock()

	jobs := []JobMeta{}

	for i := range c.jobs {
		if c.jobs[i].Label != label {
			jobs = append(jobs, c.jobs[i])
		}
	}

	if len(jobs) < len(c.jobs) {
		c.jobs = jobs
		return true
	}
	return false
}

func (c *Collector) AddJob(job JobMeta) bool {

	c.mu.Lock()
	defer c.mu.Unlock()

	for i := range c.jobs {
		if job.Label == c.jobs[i].Label {
			return false
		}
	}
	c.jobs = append(c.jobs, job)

	return true
}

func (c *Collector) UpdateJob(job JobMeta) bool {

	c.mu.Lock()
	defer c.mu.Unlock()

	for i := range c.jobs {
		if c.jobs[i].Label == job.Label {
			c.jobs[i] = job
			return true
		}
	}

	return false
}

func (c *Collector) IncMetricJobFileFailed() {
	c.mu.Lock()
	c.metrics.jobFileWatch.failedTotal += 1
	c.mu.Unlock()
}
func (c *Collector) IncMetricJobFileUnchanged() {
	c.mu.Lock()
	c.metrics.jobFileWatch.unchangedTotal += 1
	c.mu.Unlock()
}
func (c *Collector) IncMetricJobFileChanged() {
	c.mu.Lock()
	c.metrics.jobFileWatch.changedTotal += 1
	c.mu.Unlock()
}

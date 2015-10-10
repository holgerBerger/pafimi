package main

import (
	"log"
	"time"
)

// job states
const (
	Waiting   = iota
	Preparing = iota
	Running   = iota
	Completed = iota
)

// Job contains copy job description
type Job struct {
	Jobid      uint64
	User       string
	Src        string
	Dst        string
	Totalfiles uint64
	Totalbytes uint64
	Filesdone  uint64
	Bytesdone  uint64
	State      int
	Submittime time.Time
	Starttime  time.Time
	Endtime    time.Time
}

// NewJob creates a new job
func NewJob(jobid uint64, user string, src string, dst string) Job {
	job := Job{Jobid: jobid,
		User:       user,
		Src:        src,
		Dst:        dst,
		Totalfiles: 0,
		Totalbytes: 0,
		Filesdone:  0,
		Bytesdone:  0,
		State:      Waiting,
		Submittime: time.Now(),
		Starttime:  time.Unix(0, 0),
		Endtime:    time.Unix(0, 0)}
	return job
}

// Start the job
func (p *Job) Start() {
	v := *p
	v.Starttime = time.Now()
	v.State = Running
	*p = v
}

// Finish the job
func (p *Job) Finish() {
	v := *p
	v.Endtime = time.Now()
	v.State = Completed
	*p = v
	log.Println("copied", v.Filesdone, "files with", v.Bytesdone/(1024*1024), "MB")
}

// IncrFiles increment file number and number of bytes copied
func (p *Job) IncrFiles(size int64) {
	v := *p
	v.Filesdone++
	v.Bytesdone += uint64(size)
	*p = v
}

package main

import (
	"time"
)

// job states
const (
	Waiting   = iota
	Preparing = iota
	Running   = iota
	Completed = iota
)

// copy job description
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

// New creates a new job
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

func (self Job) Finish() {
	self.Endtime = time.Now()
	self.State = Completed
}

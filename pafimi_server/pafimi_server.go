package main

import (
	"fmt"
	"github.com/holgerBerger/pafimi/config"
	"log"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
)

// dummy type to identify RPC
type PafimiServerT int

// increasing global jobid, local to server
var jobid uint64 = 0

// port for RPC
var port string = "1234"

// global (server local) table of jobs, and its mutex
var (
	joblist  map[uint64]Job
	jobmutex sync.Mutex
)

// struct for channel to pass name and size to loadbalancer
type FilenameAndSize struct {
	name string
	size int64
}

// StartServer runs the endless rpc server FIX TLS needed
func StartServer() {
	l, e := net.Listen("tcp", "0.0.0.0:"+port)
	if e != nil {
		log.Fatal("listen error:", e)
	}

	// this serves endless
	rpc.Accept(l)
}

// AddRequest take RPC request and create a  jobs to server, returns jobid
func (*PafimiServerT) AddRequest(request config.Request, result *string) error {

	// increase jobid, make it atomic
	atomic.AddUint64(&jobid, 1)

	log.Println("new request with jobid", jobid, ", copy", request.Src, "to", request.Dst, "for user", request.User)

	// check if src exists
	f, err := os.Open(request.Src)
	if err != nil {
		*result = "Error: " + err.Error()
		return nil
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		*result = "Error: " + err.Error()
		return nil
	}

	// if it exists and is a directory, start go routine
	if fi.Mode().IsDir() {
		// go asynchron
		go ExecuteJob(request, jobid)
	} else {
		*result = "Error: src is not directory."
	}

	*result = "Jobid: " + strconv.FormatUint(jobid, 10)
	return nil
}

// ExecuteJob running as go routine to
//  - find filenames
//  - partition the files into lists
//  - contact other servers to execute the list
func ExecuteJob(request config.Request, jobid uint64) {
	// put job into job table
	jobmutex.Lock()
	joblist[jobid] = NewJob(jobid, request.User, request.Src, request.Dst)
	jobmutex.Unlock()

	// FIXME here we should check if numer of running jobs is large,
	// and postpone some work, use another channel?

	// create channel and worker go routine
	filechan := make(chan FilenameAndSize, 1000)
	go loadBalancer(jobid, filechan)
	// find files and send them to worker
	getFileList(request.Src, filechan)
	// done, stop worker
	close(filechan)
	// log.Println("finished job", jobid)
}

// loadBalancer receives data through channel and distributes over servers
// running for each job in parallel
func loadBalancer(jobid uint64, filechan chan FilenameAndSize) {

	servers := make([]*rpc.Client, 0)

	for _, server := range config.Conf.Client.Servers {
		// connect to rpc server
		client, err := rpc.Dial("tcp", server)
		if err != nil {
			log.Fatal("could not connect to server.")
		} else {
			servers = append(servers, client)
		}
	}

	started := false
	fileList := make([]string, 0, 100000)
	for file := range filechan {
		if !started {
			started = true
			t := joblist[jobid]
			t.Start()
			joblist[jobid] = t
		}

		// gather batches, we do not want to copy one by one
		// as we receive, but all files of one directory in one batch
		if file.name != "." && file.size != -1 {
			// gather a batch of files
			fileList = append(fileList, file.name)
		} else {
			// process batches, stripe over servers
			for s, rpcclient := range servers {
				templist := make([]string, 0)

				for i := s; i < len(fileList); i += len(servers) {
					templist = append(templist, fileList[i])
				}

				if len(templist) > 0 {
					// do RPC call
					var reply string
					err := rpcclient.Call("PafimiServerT.CopySlice", templist, &reply)
					if err != nil {
						log.Fatal(err)
					}
				}
			}
			// empty slice
			fileList = make([]string, 0, 100000)
		}

		t := joblist[jobid]
		t.IncrFiles(file.size)
		joblist[jobid] = t
	}
	t := joblist[jobid]
	t.Finish()
	joblist[jobid] = t
	log.Println("finished workers on job", jobid)
}

// getFileList append file names to filelist, idea a) is implemented here
// get directory list, sort subdirectories and files into seperate lists.
// process files fist (=copy them with loadbalancer) and descend in
// subdirectories afterwards, they should be created first
func getFileList(dir string, filechan chan FilenameAndSize) {
	fmt.Println("getFileList in", dir)
	f, _ := os.Open(dir)
	defer f.Close()

	filelist := make([]os.FileInfo, 0, 100000)
	dirlist := make([]os.FileInfo, 0, 1000)

	// get first 1000 files
	fi, err := f.Readdir(1000)
	for err == nil {
		for _, fi := range fi {
			if fi.Mode().IsDir() {
				dirlist = append(dirlist, fi)
			} else {
				filelist = append(filelist, fi)
			}
		}
		// get further files
		fi, err = f.Readdir(1000)
	}

	// create current directory in target
	// FIXME TODO
	// only needed if not created before descending, which is more efficient,
	// but what about top level?

	// process files first
	for _, fi := range filelist {
		filechan <- FilenameAndSize{name: dir + "/" + fi.Name(), size: fi.Size()}
	}
	// flush
	filechan <- FilenameAndSize{name: ".", size: -1}

	// FIXME TODO all child dirs should be created in target here

	// process child directories
	for _, fi := range dirlist {
		getFileList(dir+"/"+fi.Name(), filechan)
	}
}

// argument handling + rpc server start
func main() {

	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	// read config
	config.ReadConf()

	joblist = make(map[uint64]Job)

	server := new(PafimiServerT)
	rpc.Register(server)

	// start Copyworker, which is endless
	go CopyWorker()

	// start RPC server
	StartServer()
}

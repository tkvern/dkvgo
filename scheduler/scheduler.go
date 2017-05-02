package scheduler

import (
	"log"
	"net"
	"sync"

	"dkvgo/job"
	"dkvgo/job/store"
	"dkvgo/scheduler/tracker"
)

// DkvScheduler d
type DkvScheduler struct {
	sync.WaitGroup
	Opts        *_options
	tcpListener net.Listener
	TaskPool    *TaskPool
	Store       store.JobStore
	runningJobs map[int]*job.Job
}

func newDkvScheduler() *DkvScheduler {
	var sched = &DkvScheduler{
		Opts: Options,
		//Store: store.NewMockStore(),
		Store: store.NewDatabaseStore(Options.DBType, Options.DBAddr),
	}
	sched.TaskPool = NewTaskPool(sched)
	return sched
}

func (s *DkvScheduler) Main() {
	tracker.InitWithStore(s.Store)
	tcpListener, err := net.Listen("tcp", s.Opts.TCPAddr)
	if err != nil {
		log.Fatalf("FATAL: listen %s failed - %s\n", s.Opts.TCPAddr, err)
	}
	s.tcpListener = tcpListener
	log.Printf("TCP listen on %s\n", tcpListener.Addr())
	s.Add(1)
	s.Add(1)
	go s.runTcpServer()
	go s.runApiServer()
	s.Wait()
}

func (s *DkvScheduler) runTcpServer() {
	defer s.Done()
	TCPServer(s.tcpListener, &ProtocolLoop{ctx: s})
}

func (s *DkvScheduler) runApiServer() {
	defer s.Done()
	log.Printf("HTTP listen on %s", s.Opts.HTTPAddr)
	APIServer(s.Opts.HTTPAddr).ListenAndServe()
}

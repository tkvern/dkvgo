package scheduler

import (
	"log"
	"net"
	"sync"

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

// Main entry
func (s *DkvScheduler) Main() {
	tracker.InitWithStore(s.Store)
	tcpListener, err := net.Listen("tcp", s.Opts.TCPAddr)
	if err != nil {
		log.Fatalf("FATAL: listen %s failed - %s\n", s.Opts.TCPAddr, err)
	}
	s.tcpListener = tcpListener
	log.Printf("TCP listen on %s\n", tcpListener.Addr())
	s.Add(2)
	go s.runTCPServer()
	go s.runAPIServer()
	s.Wait()
}

// GetSplitNum return a number than how to split the job
func (s *DkvScheduler) GetSplitNum() int {
	splitNum := s.Opts.SplitNum
	if splitNum <= 0 {
		splitNum = 2
	}
	return splitNum
}

func (s *DkvScheduler) runTCPServer() {
	defer s.Done()
	TCPServer(s.tcpListener, &ProtocolLoop{ctx: s})
}

func (s *DkvScheduler) runAPIServer() {
	defer s.Done()
	log.Printf("HTTP listen on %s", s.Opts.HTTPAddr)
	APIServer(s.Opts.HTTPAddr).ListenAndServe()
}

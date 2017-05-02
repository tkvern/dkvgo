# 
all: sched workerd

sched: app/sched.go
	go build -o $(GOPATH)/bin/sched app/sched.go


workerd:
	go build -o $(GOPATH)/bin/workerd app/workerd.go



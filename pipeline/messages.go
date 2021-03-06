package pipeline

import (
	"time"
	"sync"
	"github.com/thomasschoeftner/go-cli/task"
	"fmt"
)

var _id = 0
var idLock = sync.Mutex{}
func id() int {
	idLock.Lock()
	defer idLock.Unlock()
	_id = _id + 1
	return _id - 1
}

const (
	cmdProcess = 0
	cmdStop    = 1
	cmdCancel  = 2
)

type Command struct {
	id int
	kind int
	issuedAt time.Time
	remark *string
	job task.Job
}

func newCommand(kind int, job task.Job, remark *string) Command {
	return Command{id(), kind, time.Now(), remark, job}
}

func safeCopy(job task.Job) task.Job {
	if job == nil {
		return nil
	} else {
		return job.Copy()
	}
}

func Process(params map[string]string, remark string) Command {
	return newCommand(cmdProcess, safeCopy(task.Job(params)), &remark)
}

func Stop() Command {
	return newCommand(cmdStop, nil, nil)
}

func Cancel(reason string) Command {
	return newCommand(cmdCancel, nil, &reason)
}

func (c Command) isStop() bool {
	return c.kind == cmdStop
}

func (c Command) isCancel() bool {
	return c.kind == cmdCancel
}

func (c Command) String() string {
	kind := ""
	if c.kind == cmdProcess {
		kind = "PROCESS"
	} else if c.kind == cmdCancel {
		kind = "CANCEL"
	} else if c.kind == cmdStop {
		kind = "STOP"
	} else {
		kind = "UNKNOWN"
	}
	var remark string
	if c.remark == nil {
		remark = "without remark"
	} else {
		remark = fmt.Sprintf("with remark: \"%s\"", *c.remark)
	}
	return fmt.Sprintf("%s command [%d] issued at %s %s and params %v", kind, c.id, c.issuedAt.String(), remark, c.job)
}

const (
	evtClosed   = 10
	evtCanceled = 11
	evtError    = 12
	evtJobDone  = 13
)

type Event struct {
	id int
	kind int
	createdAt time.Time
	data eventParams
}
type eventParams map[string]interface{}
const  (
	evtparam_statistics = "statistics"
	evtparam_reason     = "reason"
	evtparam_error      = "error"
	evtparam_job        = "job"
)
func newEvent(kind int, data eventParams) Event {
	return Event{id(), kind, time.Now(), data}
}

func closed(stats *Statistics) Event {
	return newEvent(evtClosed,  eventParams{evtparam_statistics : stats})
}

func canceled(reason string) Event {
	return newEvent(evtCanceled, eventParams {evtparam_reason : reason})
}

func errorIn(job task.Job, err error) Event {
	return newEvent(evtError, eventParams{evtparam_error : err, evtparam_job : safeCopy(job)})
}

func done(job task.Job) Event {
	return newEvent(evtJobDone, eventParams{evtparam_job : safeCopy(job)})
}

func (e Event) IsClosed() (bool, *Statistics) {
	if e.kind == evtClosed {
		return true, e.data[evtparam_statistics].(*Statistics)
	} else {
		return false, nil
	}
}

func (e Event) IsCanceled() (bool, string) {
	if e.kind == evtCanceled {
		return true, e.data[evtparam_reason].(string)
	} else {
		return false, ""
	}
}

func (e Event) IsError() (bool, error, task.Job) {
	if e.kind == evtError {
		return true, e.data[evtparam_error].(error), e.data[evtparam_job].(task.Job)
	} else {
		return false, nil, nil
	}
}

func (e Event) IsDone() (bool, task.Job) {
	if e.kind == evtJobDone {
		return true, e.data[evtparam_job].(task.Job)
	} else {
		return false, nil
	}
}
package tasklane

import (
	"errors"
	"time"
)

// Multiple TaskQueue => One TaskLane
// One Buffered Channel + One Blocking Channel + One Goroutine(Worker) => One TaskQueue
type TaskLane struct {
	laneSize       int
	queueSize      int
	queueList      []chan func()
	universalQueue chan func()
	timeout        time.Duration
}

func (taskLane *TaskLane) startQueue(index int) {
	for {
		taskFunc := <-taskLane.queueList[index]
		select {
		case taskLane.queueList[index+taskLane.laneSize] <- taskFunc:
		default:
			select {
			case taskLane.queueList[index+taskLane.laneSize] <- taskFunc:
			case taskLane.universalQueue <- taskFunc:
			}
		}
	}
}

func (taskLane *TaskLane) startWorker(index int) {
	var taskFunc func()
	for {
		select {
		case taskFunc = <-taskLane.queueList[index+taskLane.laneSize]:
		default:
			select {
			case taskFunc = <-taskLane.queueList[index+taskLane.laneSize]:
			case taskFunc = <-taskLane.universalQueue:
			}
		}
		taskFunc()
	}
}

func New(laneSize, queueSize int) *TaskLane {
	// create channels for TaskQueue
	queueList := make([]chan func(), laneSize*2)
	for i := 0; i < laneSize; i++ {
		queueList[i] = make(chan func(), queueSize)
		queueList[i+laneSize] = make(chan func())
	}
	universalQueue := make(chan func())

	// create TaskLane
	taskLane := &TaskLane{
		laneSize:       laneSize,
		queueSize:      queueSize,
		queueList:      queueList,
		universalQueue: universalQueue,
		timeout:        time.Second * 5,
	}

	// start TaskQueue
	for i := 0; i < taskLane.laneSize; i++ {
		go taskLane.startQueue(i)
		go taskLane.startWorker(i)
	}

	return taskLane
}

func (taskLane *TaskLane) Status() []int {
	status := make([]int, taskLane.laneSize)
	for i := 0; i < taskLane.laneSize; i++ {
		status[i] = len(taskLane.queueList[i])
	}
	return status
}

func (taskLane *TaskLane) ShortestQueueIndex() int {
	index := 0
	for i := 1; i < taskLane.laneSize; i++ {
		if len(taskLane.queueList[i]) < len(taskLane.queueList[index]) {
			index = i
		}
	}
	return index
}

func (taskLane *TaskLane) PushTask(taskFunc func(), index int) error {
	select {
	case taskLane.queueList[index] <- taskFunc:
		return nil
	case <-time.After(time.Second * 5):
		return errors.New("PushTaskTimeout")
	}
}

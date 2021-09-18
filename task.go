package fstask

import (
	"errors"
	"github.com/fsnotify/fsnotify"
	"sync"
)

// Task document task
type Task struct {
	Folder string   // folder
	Rule   string   // rule
	Action []string // fire action: create、 write、 rename、remove、chmod
	Handle func(event fsnotify.Event)
}

// FsTask document task instance
type FsTask struct {
	watcher   *fsnotify.Watcher
	ch        chan int
	taskQueue sync.Map
	lock      sync.Mutex
}

func (f *FsTask) Add(task Task) error {
	var taskList []Task

	if task.Folder != "" {
		f.watcher.Add(task.Folder)
	}

	taskID := MD5([]byte("ft_" + task.Rule))
	taskQueue, ok := f.taskQueue.Load(taskID)
	if ok {
		taskList, ok = taskQueue.([]Task)
		if !ok {
			return errors.New("incorrect task queue")
		}
	}

	f.lock.Lock()
	taskList = append(taskList, task)
	f.lock.Unlock()

	f.taskQueue.Store(task.Rule, taskList)
	return nil
}

func (f *FsTask) Wait() {
	<-f.ch
}

func (f *FsTask) Close() error {
	return f.watcher.Close()
}

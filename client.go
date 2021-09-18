package fstask

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// New initialize FsTask
func New(name string) (*FsTask, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	fsTask := FsTask{watcher: watcher, ch: make(chan int)}

	go func() {
		defer close(fsTask.ch)

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				fsTask.taskQueue.Range(func(key, value interface{}) bool {
					taskID, ok := key.(string)
					if !ok {
						return false
					}
					taskQueue, ok := fsTask.taskQueue.Load(taskID)
					if !ok {
						return false
					}

					taskList, ok := taskQueue.([]Task)
					if !ok {
						return false
					}

					for _, task := range taskList {
						reg := regexp.MustCompile(task.Rule)

						if reg.MatchString(event.Name) && Include(strings.ToLower(event.Op.String()), task.Action) {
							eventKey := MD5([]byte(event.Name + event.Op.String() + strconv.Itoa(rand.Int())))

							//ensures that it is executed only once in a given period of time
							Debounce(eventKey, func() {
								task.Handle(event)
							}, 1*time.Second)
						}
					}

					return true
				})

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	if err = watcher.Add(name); err != nil {
		return nil, err
	}

	return &fsTask, nil
}

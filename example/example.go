package main

import (
	"github.com/NICEXAI/fstask"
	"github.com/fsnotify/fsnotify"
	"log"
)

func main() {
	fsTask, err := fstask.New("./example/config")
	if err != nil {
		log.Println(err)
		return
	}
	defer fsTask.Close()

	fsTask.Add(fstask.Task{
		Rule:   ".*settings.yaml",
		Action: []string{"write"},
		Handle: func(event fsnotify.Event) {
			log.Println(event.Name, "config file change")
		},
	})

	fsTask.Wait()
}

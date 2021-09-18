# fstask

### Installation

Run the following command under your project:

> go get -u github.com/NICEXAI/fstask

### Basic Usage

If I want to listen for changes to the settings.yaml file and perform a response, you can:

```go
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
```

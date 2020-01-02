package main

import (
	"flag"
	"fmt"
	"github.com/radovskyb/watcher"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func watch(path string) {
	if path == "." {
		fmt.Println("watching default folder")
	}
	exists, err := pathExists(path)
	if err != nil {
		log.Fatal(err)
	}
	if !exists {
		log.Fatal(path, " does not exists")
	}
	configure(path)
}

func procEvent(event watcher.Event) {
	if event.Op.String() == "CREATE" && event.Name() != "" {
		if strings.Contains(event.Name(), ".mp4") {
			getFrame(event.Path, "./output/"+event.Name()+".jpeg")
		}
	}
}

func procEventt(jobs chan watcher.Event) {
	go func() {

		for {
			select {
			case event := <-jobs:
				fmt.Println("<-", event)
				procEvent(event)
			}
		}
	}()
}

func getFrame(source string, output string) {
	proc := "ffmpeg-4.2.2-amd64-static/ffmpeg"
	fmt.Println("processing ", source, proc)
	cmd := exec.Command(proc, "-y", "-i", source, "-ss", "00:00:01.000", "-vframes", "1", output)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	fmt.Printf("combined out:\n%s\n", string(out))
}

func configure(path string) {
	w := watcher.New()
	w.SetMaxEvents(25)
	w.FilterOps(watcher.Create)
	go procEventt(w.Event)
	if err := w.Add(path); err != nil {
		log.Fatalln(err)
	}
	go func() {
		w.Wait()
		w.TriggerEvent(watcher.Create, nil)
		w.TriggerEvent(watcher.Remove, nil)
	}()
	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
}

func main() {
	folderFlag := flag.String("folder", ".", "sets the folder to watch")
	flag.Parse()
	watch(*folderFlag)
}

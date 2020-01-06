package main

import (
	"flag"
	"fmt"
	"github.com/radovskyb/watcher"
	"io/ioutil"
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
			i := strings.LastIndex(event.Name(), ".")
			getFrame(event.Path, event.Name()[0:i]+".jpeg")

			//getFrame(event.Path, "./output/"+event.Name()+".jpeg")
		}
	}
}

func procEventt(jobs chan watcher.Event) {
	go func() {

		for {
			select {
			case event := <-jobs:
				procEvent(event)
			}
		}
	}()
}

func getFrame(source string, fileName string) {
	fmt.Println("->>", output+fileName)
	cmd := exec.Command(ffmpeg, "-y", "-i", source, "-ss", "00:00:01.000", "-vframes", "1", output+fileName)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	fmt.Printf("combined out:\n%s\n", string(out))
}

func configure(path string) {
	w := watcher.New()
	w.SetMaxEvents(concurrentJobs)
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

func initProcess(in string, out string) {
	inf, err := ioutil.ReadDir(in)
	if err != nil {
		log.Fatal(err)
	}

	outf, err := ioutil.ReadDir(out)
	if err != nil {
		log.Fatal(err)
	}

	for _, inff := range inf {
		matches := 0
		for _, off := range outf {
			if removeExtension(inff.Name()) == removeExtension(off.Name()) {
				matches++
			}
		}
		if matches == 0 {
			fileName := removeExtension(inff.Name()) + ".jpeg"
			fmt.Println(in + inff.Name())
			getFrame(in+inff.Name(), fileName)
		}
	}

}

func removeExtension(file string) string {
	i := strings.LastIndex(file, ".")
	return file[0:i]
}

var input, output, ffmpeg string
var concurrentJobs int

func main() {
	ffmpegFlag := flag.String("ffmpeg", "ffmpeg-4.2.2-amd64-static/ffmpeg", "ffmpeg location")
	i := flag.String("in", ".", "sets the folder to watch")
	o := flag.String("out", ".", "output folder")
	concEvents := flag.Int("concurrentJobs", 10, "sets the amount of maximum concurrent jobs execution")
	flag.Parse()
	input = *i
	output = *o
	ffmpeg = *ffmpegFlag
	concurrentJobs = *concEvents
	initProcess(input, output)
	watch(input)
}

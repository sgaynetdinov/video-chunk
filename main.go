package main

import (
	"flag"
	"fmt"
	"sync"

	yandexdisk "github.com/sgaynetdinov/go-yandex-disk"
)

var (
	workerCount     = flag.Int("w", 10, "worker count")
	configFile      = flag.String("c", "", "config file")
	yandexdiskToken = flag.String("token", "", "yandex disk token")
)

var client = yandexdisk.NewClient(*yandexdiskToken)

func workerScheduler(url, folder string) {
	fmt.Printf("\n\nStart download: %s\n", folder)
	var wg sync.WaitGroup

	client.CreateFolder(folder)

	finishCh := make(chan bool)
	urlCh := generateUrl(url, 0)

	for i := 1; i < *workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				select {
				case <-finishCh:
					return
				default:
					url := <-urlCh
					err := Worker(url, folder)
					if err != nil {
						finishCh <- true
						return
					}
					fmt.Print(".")
				}
			}
		}()
	}

	for {
		select {
		case <-finishCh:
			return
		}
	}

	defer wg.Wait()
}

func main() {
	flag.Parse()

	videoItems, err := parseConfig(*configFile)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Worker count: %d\n", *workerCount)
	fmt.Printf("Video count: %d\n", len(videoItems))

	for _, video := range videoItems {
		workerScheduler(video.Url, video.Folder)
	}
}

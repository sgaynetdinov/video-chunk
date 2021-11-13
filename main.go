package main

import (
	"flag"
	"fmt"
	"sync"

	yandexdisk "github.com/sgaynetdinov/go-yandex-disk"
)

var (
	client          *yandexdisk.Client
	workerCount     = flag.Int("w", 10, "worker count")
	configFile      = flag.String("c", "", "config file")
	yandexdiskToken = flag.String("token", "", "yandex disk token")
)

func workerScheduler(url, folder string) {
	var wg sync.WaitGroup

	err := client.MkdirAll(folder)
	if err != nil {
		panic(err)
	}

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
					err := RetryWorker(5, url, folder)
					if err != nil {
						if err != ErrNotFound {
							fmt.Println("Error")
						}
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

	client = yandexdisk.NewClient(*yandexdiskToken)

	videoItems, err := parseConfig(*configFile)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Worker count: %d\n", *workerCount)
	fmt.Printf("Video count: %d\n", len(videoItems))

	for _, video := range videoItems {
		fmt.Printf("\n\nStart download: %s\n", folder)
		workerScheduler(video.Url, video.Folder)
	}
}

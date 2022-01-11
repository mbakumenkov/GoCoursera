package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

// Function pipeline
func ExecutePipeline(jobs ...job) {
	wg := &sync.WaitGroup{}
	var in, out chan interface{}

	for _, j := range jobs {
		in = out
		// "You can never expect to have more than 100 input items" (c)
		out = make(chan interface{}, 100)
		wg.Add(1)

		go func(in, out chan interface{}, j job) {
			defer wg.Done()
			defer close(out)

			j(in, out)
		}(in, out, j)
	}

	wg.Wait()
}

// Async read from [in] channel, apply [function] and write to [out]
func asyncApplyFuncToString(in, out chan interface{}, function func(string) string) {
	wg := &sync.WaitGroup{}
	for data := range in {
		wg.Add(1)
		go func(data string) {
			defer wg.Done()

			out <- function(data)
		}(fmt.Sprintf("%v", data))
	}
	wg.Wait()
}

// Calculate crc32(data) + '~` + crc32(md5(data))
func SingleHash(in, out chan interface{}) {
	md5mtx := &sync.Mutex{}
	asyncApplyFuncToString(in, out, func(data string) string {
		hashCrcCh := make(chan string)
		go func() {
			defer close(hashCrcCh)

			hashCrcCh <- DataSignerCrc32(fmt.Sprintf("%v", data))
		}()

		hashCrcMd5Ch := make(chan string)
		go func() {
			defer close(hashCrcMd5Ch)

			md5mtx.Lock()
			md5hash := DataSignerMd5(fmt.Sprintf("%v", data))
			md5mtx.Unlock()

			hashCrcMd5Ch <- DataSignerCrc32(md5hash)
		}()

		return fmt.Sprintf("%s~%s", <-hashCrcCh, <-hashCrcMd5Ch)
	})
}

// Calculate crc32(th + data(from [in] channel)), where th is [0..5] and concatenate it in order of th
func MultiHash(in, out chan interface{}) {
	asyncApplyFuncToString(in, out, func(data string) string {
		resWg := &sync.WaitGroup{}
		resMtx := &sync.Mutex{}
		var results [6]string

		for th := 0; th < 6; th++ {
			resWg.Add(1)

			go func(th int) {
				defer resWg.Done()

				hash := DataSignerCrc32(fmt.Sprintf("%d%s", th, data))

				resMtx.Lock()
				results[th] = hash
				resMtx.Unlock()
			}(th)
		}
		resWg.Wait()
		return strings.Join(results[:], "")
	})
}

// Join all results from [in] channel with "_" and write it to [out] channel
func CombineResults(in, out chan interface{}) {
	var results []string

	for data := range in {
		results = append(results, fmt.Sprintf("%v", data))
	}

	sort.Strings(results)
	out <- strings.Join(results, "_")
}

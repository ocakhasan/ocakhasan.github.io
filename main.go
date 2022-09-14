package main

import "sync"

func handleWebsites(websites []string) error {
	errChan := make(chan error, 1)
	semaphores := make(chan struct{}, 5) // aynı anda 5 iş çalıştır
	var wg sync.WaitGroup
	wg.Add(len(websites))
	for _, website := range websites {
		semaphores <- struct{}{} // semaphore acquire et
		go func() {
			defer func() {
				wg.Done()
				<-semaphores
			}()
			if err := handle(website); err != nil {
				errChan <- err
			}
		}()
	}
	wg.Wait()
	close(semaphores)
	close(errChan)
	return <-errChan
}

package player

import "os"

// WriteFileChan takes in a file and returns a channel that when you
// send on it, that line will be written to the file
func WriteFileChan(file *os.File, lines <-chan string) (done <-chan string) {
	doneRW := make(chan string)
	go func() {
		for line := range lines {
			_, err := file.WriteString(line + "\n")
			if err != nil {
				panic(err)
			}
		}
		close(doneRW)
	}()
	return doneRW
}

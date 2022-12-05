package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"runtime/pprof"
	"time"

	"github.com/douglasmakey/pocketknife/tracker"
)

func main() {

	f, err := os.Create("netsock.prof")
	if err != nil {
		log.Fatal(err)
	}

	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	defer tracker.LogTimeTrack(time.Now(), "10k request to network socket")

	for i := 0; i < 10000; i++ {
		conn, err := net.Dial("tcp", "localhost:8000")
		if err != nil {
			log.Fatal(err)
		}

		msg := "I'm a Kungfu Dev"
		if _, err := conn.Write([]byte(msg)); err != nil {
			log.Fatal(err)
		}

		buf := make([]byte, len(msg))
		n, err := conn.Read(buf)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(buf[:n]))
	}
}

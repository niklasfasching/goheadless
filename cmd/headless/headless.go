package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/niklasfasching/goheadless"
)

var code = flag.String("c", "", "code to run after files have been imported")
var args = flag.String("a", "", "window.args - split via strings.Fields")

func main() {
	log.SetFlags(0)
	flag.Parse()
	h := &goheadless.Runner{}
	if err := h.Start(); err != nil {
		log.Fatal(err)
	}
	defer h.Stop()
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		h.Stop()
		os.Exit(1)
	}()
	html := goheadless.HTML(*code, flag.Args(), strings.Fields(*args))
	messages := h.Run(context.Background(), html)
	for m := range messages {
		if m.Method == "clear" && len(m.Args) == 1 {
			exitCode, ok := m.Args[0].(float64)
			if !ok {
				os.Exit(-1)
			}
			os.Exit(int(exitCode))
		} else if m.Method == "info" {
			log.Println(goheadless.Colorize(m))
		} else {
			log.Println(append([]interface{}{m.Method}, m.Args...)...)
		}
	}
}

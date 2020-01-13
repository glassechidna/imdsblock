package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
)

func main() {
	usage := "usage: imdsblock (start|pre-start|post-stop)"
	if len(os.Args) == 1 {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "start":
		start()
	case "pre-start":
		preStart()
	case "post-stop":
		postStop()
	default:
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(1)
	}
}

func start() {
	router := mux.NewRouter()

	// just return 404 for anything iam (this is what aws does)
	router.PathPrefix("/{version}/meta-data/iam").Handler(router.NotFoundHandler)

	// otherwise pass it all through
	u, _ := url.Parse("http://169.254.169.254")
	passthrough := httputil.NewSingleHostReverseProxy(u)
	router.PathPrefix("/").Handler(passthrough)

	err := http.ListenAndServe(":51999", router)
	panic(err)
}

func preStart() {
	if !iptables(opAdd) {
		os.Exit(1)
	}
}

func postStop() {
	if !iptables(opRemove) {
		os.Exit(1)
	}
}

const opAdd = "-A"
const opRemove = "-D"

func iptables(op string) bool {
	args := []string{
		"-t", "nat",
		op, "PREROUTING",
		"-p", "tcp",
		"-d", "169.254.169.254",
		"--dport", "80",
		"-j", "DNAT",
		"--to-destination", "127.0.0.1:51999",
	}

	cmd := exec.Command("/sbin/iptables", args...)
	err := cmd.Run()
	return err == nil
}

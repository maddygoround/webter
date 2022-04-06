package main

import "os"

func main() {
	stunServer := "stun:stun.l.google.com:19302"
	oneWay := true
	cmd := []string{"bash", "-l"}
	for i, arg := range os.Args {
		if arg == "-cmd" {
			cmd = os.Args[i+1:]
			os.Args = os.Args[:i]
		}
	}
	hc := hostSession{
		oneWay: oneWay,
		cmd:    cmd,
	}
	hc.stunServers = []string{stunServer}
	_ = hc.run()

}

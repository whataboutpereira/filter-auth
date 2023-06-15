
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"log"
)

type session struct {
	id   string
	rdns string
	ip   string
	user string
}

var version string
var sessions = make(map[string]*session)

func linkConnect(s *session, params []string) {
	if len(params) != 4 {
		log.Fatal("invalid input, shouldn't happen")
	}

	s.rdns = params[0]

	if strings.HasPrefix(params[2], "unix:") {
		s.ip = "127.0.0.1"
	} else if params[2][0] == '[' {
		s.ip = strings.Split(strings.Split(params[2], "]")[0], "[")[1]
	} else {
		s.ip = strings.Split(params[2], ":")[0]
	}
}

func linkAuth(s *session, params []string) {
	if len(params) < 2 {
		log.Fatal("invalid input, shouldn't happen")
	}

	var res string

	if version < "0.7" {
		res = params[len(params) - 1]
		s.user = strings.Join(params[0:len(params)-1], "|")
	} else {
		res = params[0]
		s.user = strings.Join(params[1:], "|")
	}

	if res != "pass" {
		fmt.Fprintf(os.Stderr, "failed authentication from user=%s address=%s host=%s\n", s.user, s.ip, s.rdns)
		return
	}
}

func linkDisconnect(s *session, params []string) {
	if len(params) != 0 {
		log.Fatal("invalid input, shouldn't happen")
	}
	delete(sessions, s.id)
}

func main() {
	input := bufio.NewScanner(os.Stdin)
	reporters := map[string]func(*session, []string) {
		"link-connect":    linkConnect,
		"link-disconnect": linkDisconnect,
		"link-auth":       linkAuth,
	}

	for {
		if !input.Scan() {
			os.Exit(0)
		}

		if input.Text() == "config|ready" {
			for k := range reporters {
				fmt.Printf("register|report|smtp-in|%s\n", k)
			}

			fmt.Println("register|ready")

			break
		}
	}

	for {
		if !input.Scan() {
			os.Exit(0)
		}

		bits := strings.Split(input.Text(), "|")
		if len(bits) < 6 {
			os.Exit(1)
		}

		version = bits[1]

		if bits[0] == "report" {
			if bits[4] == "link-connect" {
				sessions[bits[5]] = &session{ id : bits[5] }
			}

			s := sessions[bits[5]]

			if v, ok := reporters[bits[4]]; ok {
				v(s, bits[6:])
			} else {
				os.Exit(1)
			}
		} else {
			os.Exit(1)
		}
	}
}

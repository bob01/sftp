package main

import (
	"fmt"
	"github.com/bob01/sftp/internal"
	"golang.org/x/crypto/ssh"
	"encoding/hex"
	"flag"
	"io/ioutil"
	"encoding/json"
	"os"
)


type Conf struct {
	Addr string `json:"addr"`
	User string `json:"user"`
	Password string `json:"password"`
	Path string `json:"path"`
	Flags uint `json:"flags"`
}

func main()  {

	// usage
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "SFTP v6 protocol STAT test v1.8.0a")
		flag.PrintDefaults()
	}

	// parse args
	var conf Conf
	confFile := flag.String("conf", "", "JSON configuration file eg.\n{\n  \"addr\": \"<host>:<port>\",\n  \"user\": \"username\",\n  \"password\": \"<password>\",\n  \"path\": \"<sftp-object-path>\"\n}")

	addr := flag.String("addr", conf.Addr, "Host address <host>[:<port>]")
	username := flag.String("user", conf.User, "Username")
	password := flag.String("password", conf.Password, "Password")
	path := flag.String("path", conf.Path, "File system object path")
	flags := flag.Uint("flags", conf.Flags, "Attribute flags")

	flag.Parse()

	// if conf file specified repeat after reading to use conf values as defaults
	if len(*confFile) != 0 {
		// read
		raw, e := ioutil.ReadFile(*confFile)
		if e != nil {
			fmt.Println(e)
			return
		}

		// unmarshal
		var config Conf
		json.Unmarshal(raw, &config)

		// parse again to use conf values as defaults for flag values
		flag.Parse()
	}

	// use values from confFile file if provided for args not explicitly specified
	if len(*confFile) != 0 {
		// read
		raw, e := ioutil.ReadFile(*confFile)
		if e != nil {
			fmt.Println(e)
			return
		}

		// unmarshal
		var config Conf
		json.Unmarshal(raw, &config)

		if len(*addr) == 0 {
			addr = &config.Addr
		}
		if len(*username) == 0 {
			username = &config.User
		}
		if len(*password) == 0 {
			password = &config.Password
		}
		if len(*path) == 0 {
			path = &config.Path
		}
		if *flags == 0 {
			flags = &config.Flags
		}
	}

	// open ssh connection
	config := &ssh.ClientConfig{
		User: *username,
		HostKeyCallback:ssh.InsecureIgnoreHostKey(),
		Auth:[]ssh.AuthMethod{
			ssh.Password(*password),
		},
	}
	config.SetDefaults()
	conn, e := ssh.Dial("tcp", *addr, config)
	if e != nil {
		fmt.Println(e)
		return
	}
	defer conn.Close()

	// open version 6 sftp session
	client, e := sftp.NewClient(conn)
	if e != nil {
		fmt.Println(e)
		return
	}
	defer client.Close()

	// request STAT
	b, e := client.StatP6(*path, uint32(*flags))
	if e != nil {
		fmt.Println(e)
		return
	}

	// append 32 bytes
	extra32 := make([]byte, 32)
	b = append(b, extra32...)

	// display & parse packet
	fmt.Printf("%s\n", hex.Dump(b))
	fmt.Printf("%s", sftp.DumpP6Attrs(b))

	fmt.Println("done.")
}

package main

import (
	"fmt"
	"github.com/bob01/sftp/internal"
	"golang.org/x/crypto/ssh"
	"encoding/hex"
	"flag"
)


func main()  {

	addr := flag.String("addr", "", "Host address <host>[:<port>]")
	username := flag.String("user", "", "Username")
	password := flag.String("password", "", "Password")
	path := flag.String("path", "", "File system object path")
	flags := flag.Uint("flags", sftp.SSH_FILEXFER_ATTR_SIZE, "Attribute flags")
	flag.Parse()

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

	client, e := sftp.NewClient(conn)
	if e != nil {
		fmt.Println(e)
		return
	}
	defer client.Close()

	b, e := client.StatP6(*path, uint32(*flags))
	if e != nil {
		fmt.Println(e)
		return
	}

	// extra 32 bytes
	extra32 := make([]byte, 32)
	b = append(b, extra32...)

	fmt.Printf("%s\n", hex.Dump(b))
	fmt.Printf("%s", sftp.DumpP6Attrs(b))

	fmt.Println("done.")
}


//infos, e := client.ReadDir(".")
//if e != nil {
//	fmt.Println(e)
//	return
//}
//for _, fi := range infos {
//	fmt.Printf(".name='%v', size=%d, mtime=%v, dir=%v\n", fi.Name(), fi.Size(), fi.ModTime(), fi.IsDir())
//}

//fi, e := client.Stat("1059.diff")
//if e != nil {
//	fmt.Println(e)
//	return
//}
//fmt.Printf("STAT: name='%v', size=%d, mtime=%v, dir=%v\n", fi.Name(), fi.Size(), fi.ModTime(), fi.IsDir())
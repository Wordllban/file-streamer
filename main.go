package main

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type FileServer struct{}

func (fs *FileServer) start() {
	ln, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go fs.readLoop(conn)
	}

}

func (fs *FileServer) readLoop(conn net.Conn) {
	buf := new(bytes.Buffer)
	for {
		var size int64
		binary.Read(conn, binary.LittleEndian, &size)
		nBytes, err := io.CopyN(buf, conn, size)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(buf.Bytes())
		fmt.Printf("Received %d bytes over network \n", nBytes)
	}
}

// simple function to mock file streaming over network, almost like from client side
func mockSendFile(size int) error {
	file := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, file)
	if err != nil {
		return err
	}

	conn, err := net.Dial("tcp", ":3000")
	if err != nil {
		return err
	}

	binary.Write(conn, binary.LittleEndian, int64(size))
	nBytes, err := io.CopyN(conn, bytes.NewReader(file), int64(size))
	if err != nil {
		return err
	}

	fmt.Printf("Written %d bytes over the network\n", nBytes)

	return nil
}

func main() {
	go func() {
		time.Sleep(4 * time.Second)
		mockSendFile(200000)
	}()

	server := &FileServer{}
	server.start()
}

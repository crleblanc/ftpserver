package main

import (
	"fmt"
	"github.com/jlaffaye/ftp" // ftp client for testing
	"strconv"
	"testing"
	"time"
)

// Integration tests for ftpserver using a running client and server

var (
	testClient *ftp.ServerConn
)

func init() {
	// run the FTP server in a goroutine so it can respond to our test client
	go main()

	// give the server some time to start
	time.Sleep(time.Millisecond * 5)

	u := ftpServer.Settings.Host + ":" + strconv.Itoa(ftpServer.Settings.Port)
	var err error
	if testClient, err = ftp.DialTimeout(u, time.Second); err != nil {
		panic(err)
	}
}

// TODO: we need a test server and test client.  The server needs to run main() in a
// goroutine, the client can be blocking.  Test all implemented FTP commands with expected
// passing and failing tests.  Should be able to reproduce the panic I found and fix.  Use the
// 'virtual driver' for testing.

func TestLogin(t *testing.T) {
	var err error
	testCases := []struct {
		user       string
		passwd     string
		shouldFail bool
	}{
		{"", "", false},
		{"bad", "", true},
		{"", "bad", true},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("user=%s passwd=%s shouldFail=%v", tc.user, tc.passwd, tc.shouldFail), func(t *testing.T) {

			if err = testClient.Login(tc.user, tc.passwd); (err != nil) != tc.shouldFail {
				t.Error(err)
			}

			if tc.shouldFail {
				return
			}

			if err = testClient.Quit(); err != nil {
				t.Error(err)
			}
		})
	}

}

//func TestPutGet(t *testing.T) {
//	var err error
//	testCases := []struct {
//		local       string
//		remote     string
//		shouldFail bool
//	}{
//		{"", "", false},
//		{"", "", false},
//		{"", "", false},
//	}
//
//	for _, tc := range testCases {
//		t.Run(fmt.Sprintf("user=%s passwd=%s shouldFail=%v", tc.local, tc.remote, tc.shouldFail), func(t *testing.T) {
//
//			//if err = testClient.Login(tc.user, tc.passwd); (err != nil) != tc.shouldFail {
//			//	t.Error(err)
//			//}
//
//			if tc.shouldFail {
//				return
//			}
//
//		})
//	}
//
//}

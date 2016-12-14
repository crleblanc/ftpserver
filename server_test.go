package main

import (
	"bytes"
	"fmt"
	"github.com/jlaffaye/ftp" // ftp client for testing
	"gopkg.in/inconshreveable/log15.v2"
	"io"
	"io/ioutil"
	"strconv"
	"testing"
	"time"
)

// Integration tests for ftpserver using a running client and server

func init() {

	// disable logging for clean test reports
	logger := log15.Root()
	logger.SetHandler(log15.DiscardHandler())

	// run the FTP server in a goroutine so it can respond to our test client
	go main()

	// give the server some time to start
	time.Sleep(time.Millisecond * 5)
}

func getTestClient() (tc *ftp.ServerConn, err error) {
	u := ftpServer.Settings.Host + ":" + strconv.Itoa(ftpServer.Settings.Port)
	if tc, err = ftp.DialTimeout(u, time.Second); err != nil {
		return nil, err
	}

	return tc, nil
}

// TODO: we need a test server and test client.  The server needs to run main() in a
// goroutine, the client can be blocking.  Test all implemented FTP commands with expected
// passing and failing tests.  Should be able to reproduce the panic I found and fix.  Use the
// 'virtual driver' for testing.

func TestLogin(t *testing.T) {
	testCases := []struct {
		user       string
		passwd     string
		shouldFail bool
	}{
		{"", "", false},
		{"bad", "", true},
		{"", "bad", true},
	}

	var err error
	var testClient *ftp.ServerConn

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("user=%s passwd=%s shouldFail=%v", tc.user, tc.passwd, tc.shouldFail),
			func(t *testing.T) {

				if testClient, err = getTestClient(); err != nil {
					t.Fatal(err)
				}
				defer testClient.Quit()

				if err = testClient.Login(tc.user, tc.passwd); (err != nil) != tc.shouldFail {
					t.Error(err)
				}
			})
	}

}

func TestPutGet(t *testing.T) {
	testCases := []struct {
		local      string
		remote     string
		shouldFail bool
	}{
		{"localfile", "remotefile", false},
		//{"localfile", "remotefile/invalid_directory", false},
	}

	testString := "test data"
	data := bytes.NewBufferString(testString)

	var err error
	var testClient *ftp.ServerConn

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("local=%s remote=%s shouldFail=%v", tc.local, tc.remote, tc.shouldFail),
			func(t *testing.T) {

				if testClient, err = getTestClient(); err != nil {
					t.Fatal(err)
				}
				defer testClient.Quit()

				if err = testClient.Login("", ""); err != nil {
					t.Error(err)
				}

				// upload the file
				if err = testClient.Stor(tc.remote, data); (err != nil) != tc.shouldFail {
					t.Error(err)
				}

				// get the file back
				var reader io.ReadCloser
				if reader, err = testClient.Retr(tc.remote); (err != nil) != tc.shouldFail {
					t.Error(err)
				}
				defer reader.Close()

				var dataRead []byte
				if dataRead, err = ioutil.ReadAll(reader); err != nil {
					t.Error(err)
				}

				if testString != string(dataRead) {
					t.Errorf("Strings do not match, expected: [%s] but saw [%s]", testString, string(dataRead))
				}

			})
	}

}

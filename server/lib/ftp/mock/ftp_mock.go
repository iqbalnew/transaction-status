package FtpMock

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"net/textproto"
	"strconv"
	"strings"
	"sync"
	"time"

	ftpwrapper "bitbucket.bri.co.id/scm/bricams-addons/transaction-status/server/lib/ftp/wrapper"

	"github.com/jlaffaye/ftp"
)

const timeFormat = "20060102150405"

type FtpMock struct {
	address  string
	modtime  string // no-time, std-time, vsftpd
	listener *net.TCPListener
	proto    *textproto.Conn
	commands []string // list of received commands
	lastFull string   // full last command
	rest     int
	fileCont *bytes.Buffer
	dataConn *mockDataConn
	sync.WaitGroup
	DataContent string

	// #region Wrapper for negative test
	Conn *ftp.ServerConn
	ftpwrapper.FtpInterface
	ftpwrapper.FtpConnectionInterface
	// #endregion
}

// newFtpMock returns a mock implementation of a FTP server
// For simplication, a mock instance only accepts a signle connection and terminates afer
func NewFtpMock(address string) (*FtpMock, error) {
	return newFtpMockExt(address, "no-time")
}

func newFtpMockExt(address, modtime string) (*FtpMock, error) {
	var err error
	mock := &FtpMock{
		address: address,
		modtime: modtime,
	}

	l, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	tcpListener, ok := l.(*net.TCPListener)
	if !ok {
		return nil, errors.New("listener is not a net.TCPListener")
	}
	mock.listener = tcpListener

	go mock.listen()

	return mock, nil
}

func (mock *FtpMock) listen() {
	// Listen for an incoming connection.
	conn, err := mock.listener.Accept()
	if err != nil {
		panic(fmt.Sprintf("can not accept: %s", err))
	}

	// Do not accept incoming connections anymore
	mock.listener.Close()

	mock.Add(1)
	defer mock.Done()
	defer conn.Close()

	mock.proto = textproto.NewConn(conn)
	mock.printfLine("220 FTP Server ready.")

	for {
		fullCommand, _ := mock.proto.ReadLine()
		mock.lastFull = fullCommand

		cmdParts := strings.Split(fullCommand, " ")

		// Append to list of received commands
		mock.commands = append(mock.commands, cmdParts[0])

		// At least one command must have a multiline response
		switch cmdParts[0] {
		case "FEAT":
			features := "211-Features:\r\n FEAT\r\n PASV\r\n EPSV\r\n UTF8\r\n SIZE\r\n MLST\r\n"
			switch mock.modtime {
			case "std-time":
				features += " MDTM\r\n MFMT\r\n"
			case "vsftpd":
				features += " MDTM\r\n"
			}
			features += "211 End"
			mock.printfLine(features)
		case "USER":
			if cmdParts[1] == "anonymous" {
				mock.printfLine("331 Please send your password")
			} else {
				mock.printfLine("530 This FTP server is anonymous only")
			}
		case "PASS":
			mock.printfLine("230-Hey,\r\nWelcome to my FTP\r\n230 Access granted")
		case "TYPE":
			mock.printfLine("200 Type set ok")
		case "CWD":
			if cmdParts[1] == "missing-dir" {
				mock.printfLine("550 %s: No such file or directory", cmdParts[1])
			} else {
				mock.printfLine("250 Directory successfully changed.")
			}
		case "DELE":
			mock.printfLine("250 File successfully removed.")
		case "MKD":
			mock.printfLine("257 Directory successfully created.")
		case "RMD":
			if cmdParts[1] == "missing-dir" {
				mock.printfLine("550 No such file or directory")
			} else {
				mock.printfLine("250 Directory successfully removed.")
			}
		case "PWD":
			mock.printfLine("257 \"/incoming\"")
		case "CDUP":
			mock.printfLine("250 CDUP command successful")
		case "SIZE":
			if cmdParts[1] == "magic-file" {
				mock.printfLine("213 42")
			} else {
				mock.printfLine("550 Could not get file size.")
			}
		case "PASV":
			p, err := mock.listenDataConn()
			if err != nil {
				mock.printfLine("451 %s.", err)
				break
			}

			p1 := int(p / 256)
			p2 := p % 256

			mock.printfLine("227 Entering Passive Mode (127,0,0,1,%d,%d).", p1, p2)
		case "EPSV":
			p, err := mock.listenDataConn()
			if err != nil {
				mock.printfLine("451 %s.", err)
				break
			}
			mock.printfLine("229 Entering Extended Passive Mode (|||%d|)", p)
		case "STOR":
			if mock.dataConn == nil {
				mock.printfLine("425 Unable to build data connection: Connection refused")
				break
			}
			mock.printfLine("150 please send")
			mock.recvDataConn(false)
		case "APPE":
			if mock.dataConn == nil {
				mock.printfLine("425 Unable to build data connection: Connection refused")
				break
			}
			mock.printfLine("150 please send")
			mock.recvDataConn(true)
		case "LIST":
			if mock.dataConn == nil {
				mock.printfLine("425 Unable to build data connection: Connection refused")
				break
			}

			mock.dataConn.Wait()
			mock.printfLine("150 Opening ASCII mode data connection for file list")
			mock.dataConn.write([]byte("-rw-r--r--   1 ftp      wheel           0 Jan 29 10:29 lo\r\ntotal 1"))
			mock.printfLine("226 Transfer complete")
			mock.closeDataConn()
		case "MLSD":
			if mock.dataConn == nil {
				mock.printfLine("425 Unable to build data connection: Connection refused")
				break
			}

			mock.dataConn.Wait()
			mock.printfLine("150 Opening data connection for file list")
			mock.dataConn.write([]byte("Type=file;Size=0;Modify=20201213202400; lo\r\n"))
			mock.printfLine("226 Transfer complete")
			mock.closeDataConn()
		case "MLST":
			if cmdParts[1] == "multiline-dir" {
				mock.printfLine("250-File data\r\n Type=dir;Size=0; multiline-dir\r\n Modify=20201213202400; multiline-dir\r\n250 End")
			} else {
				mock.printfLine("250-File data\r\n Type=file;Size=42;Modify=20201213202400; magic-file\r\n \r\n250 End")
			}
		case "NLST":
			if mock.dataConn == nil {
				mock.printfLine("425 Unable to build data connection: Connection refused")
				break
			}

			mock.dataConn.Wait()
			mock.printfLine("150 Opening ASCII mode data connection for file list")
			mock.dataConn.write([]byte("/incoming"))
			mock.printfLine("226 Transfer complete")
			mock.closeDataConn()
		case "RETR":
			if mock.dataConn == nil {
				mock.printfLine("425 Unable to build data connection: Connection refused")
				break
			}

			mock.dataConn.Wait()
			mock.printfLine("150 Opening ASCII mode data connection for file list")
			mock.fileCont = new(bytes.Buffer)
			mock.fileCont.WriteString(mock.DataContent)
			mock.dataConn.write(mock.fileCont.Bytes()[mock.rest:])
			mock.rest = 0
			mock.printfLine("226 Transfer complete")
			mock.closeDataConn()
		case "RNFR":
			mock.printfLine("350 File or directory exists, ready for destination name")
		case "RNTO":
			mock.printfLine("250 Rename successful")
		case "REST":
			if len(cmdParts) != 2 {
				mock.printfLine("500 wrong number of arguments")
				break
			}
			rest, err := strconv.Atoi(cmdParts[1])
			if err != nil {
				mock.printfLine("500 REST: %s", err)
				break
			}
			mock.rest = rest
			mock.printfLine("350 Restarting at %s. Send STORE or RETRIEVE to initiate transfer", cmdParts[1])
		case "MDTM":
			var answer string
			switch {
			case mock.modtime == "no-time":
				answer = "500 Unknown command MDTM"
			case len(cmdParts) == 3 && mock.modtime == "vsftpd":
				answer = "213 UTIME OK"
				_, err := time.ParseInLocation(timeFormat, cmdParts[1], time.UTC)
				if err != nil {
					answer = "501 Can't get a time stamp"
				}
			case len(cmdParts) == 2:
				answer = "213 20201213202400"
			default:
				answer = "500 wrong number of arguments"
			}
			mock.printfLine(answer)
		case "MFMT":
			var answer string
			switch {
			case mock.modtime == "std-time" && len(cmdParts) == 3:
				answer = "213 UTIME OK"
				_, err := time.ParseInLocation(timeFormat, cmdParts[1], time.UTC)
				if err != nil {
					answer = "501 Can't get a time stamp"
				}
			default:
				answer = "500 Unknown command MFMT"
			}
			mock.printfLine(answer)
		case "NOOP":
			mock.printfLine("200 NOOP ok.")
		case "OPTS":
			if len(cmdParts) != 3 {
				mock.printfLine("500 wrong number of arguments")
				break
			}
			if (strings.Join(cmdParts[1:], " ")) == "UTF8 ON" {
				mock.printfLine("200 OK, UTF-8 enabled")
			}
		case "REIN":
			mock.printfLine("220 Logged out")
		case "QUIT":
			mock.printfLine("221 Goodbye.")
			return
		default:
			mock.printfLine("500 Unknown command %s.", cmdParts[0])
		}
	}
}

func (mock *FtpMock) printfLine(format string, args ...interface{}) {
	if err := mock.proto.Writer.PrintfLine(format, args...); err != nil {
		panic(err)
	}
}

func (mock *FtpMock) closeDataConn() {
	if mock.dataConn != nil {
		if err := mock.dataConn.Close(); err != nil {
			panic(err)
		}
		mock.dataConn = nil
	}
}

type mockDataConn struct {
	listener *net.TCPListener
	conn     net.Conn
	// WaitGroup is done when conn is accepted and stored
	sync.WaitGroup
}

func (d *mockDataConn) Close() (err error) {
	if d.listener != nil {
		err = d.listener.Close()
	}
	if d.conn != nil {
		err = d.conn.Close()
	}
	return
}

func (d *mockDataConn) write(b []byte) {
	if d.conn == nil {
		panic("data conn is not opened")
	}

	if _, err := d.conn.Write(b); err != nil {
		panic(err)
	}
}

func (mock *FtpMock) listenDataConn() (int64, error) {
	mock.closeDataConn()

	l, err := net.Listen("tcp", mock.address)
	if err != nil {
		return 0, err
	}

	tcpListener, ok := l.(*net.TCPListener)
	if !ok {
		return 0, errors.New("listener is not a net.TCPListener")
	}

	addr := tcpListener.Addr().String()

	_, port, err := net.SplitHostPort(addr)
	if err != nil {
		return 0, err
	}

	p, err := strconv.ParseInt(port, 10, 32)
	if err != nil {
		return 0, err
	}

	dataConn := &mockDataConn{
		listener: tcpListener,
	}
	dataConn.Add(1)

	go func() {
		// Listen for an incoming connection.
		conn, err := dataConn.listener.Accept()
		if err != nil {
			// mock.t.Fatalf("can not accept data conn: %s", err)
			return
		}

		dataConn.conn = conn
		dataConn.Done()
	}()

	mock.dataConn = dataConn
	return p, nil
}

func (mock *FtpMock) recvDataConn(append bool) {
	mock.dataConn.Wait()
	if !append {
		mock.fileCont = new(bytes.Buffer)
	}

	if _, err := io.Copy(mock.fileCont, mock.dataConn.conn); err != nil {
		panic(err)
	}

	mock.printfLine("226 Transfer Complete")
	mock.closeDataConn()
}

func (mock *FtpMock) Addr() string {
	return mock.listener.Addr().String()
}

// Closes the listening socket
func (mock *FtpMock) Close() {
	mock.listener.Close()
}

func (mock *FtpMock) Dial(addr string, options ...ftp.DialOption) (*ftp.ServerConn, error) {
	if addr == "" {
		return nil, errors.New("cannot connect to ftp server")
	}

	return nil, nil
}

func (mock *FtpMock) Login(user string, password string) error {
	if user == "" {
		return errors.New("login failed")
	}
	return nil
}

func (mock *FtpMock) Logout() error {
	return nil
}

func (mock *FtpMock) Quit() error {
	return nil
}

func (mock *FtpMock) List(path string) ([]*ftp.Entry, error) {
	return []*ftp.Entry{}, nil
}

func (mock *FtpMock) ChangeDir(path string) error {
	if path == "" {
		return errors.New("path not found")
	}
	return nil
}

func (mock *FtpMock) Retr(path string) (*ftp.Response, error) {
	if path == "" {
		return nil, errors.New("file not found")
	}

	return &ftp.Response{}, nil
}

func (mock *FtpMock) Stor(path string, r io.Reader) error {
	if path == "" {
		return errors.New("failed write file")
	}

	return nil
}

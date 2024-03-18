package ftpwrapper

import (
	"io"

	"github.com/jlaffaye/ftp"
)

type FtpInterface interface {
	Dial(addr string, options ...ftp.DialOption) (*ftp.ServerConn, error)
}

type FtpConnectionInterface interface {
	ChangeDir(path string) error
	Login(user string, password string) error
	Logout() error
	List(path string) (entries []*ftp.Entry, err error)
	Quit() error
	Retr(path string) (*ftp.Response, error)
	Stor(path string, r io.Reader) error
}

package ftpwrapper

import (
	"io"

	"github.com/jlaffaye/ftp"
)

type FtpWrapper struct {
	FtpInterface
}

type FtpConnectionWrapper struct {
	conn *ftp.ServerConn
	FtpConnectionInterface
}

func (fw *FtpWrapper) Dial(addr string, options ...ftp.DialOption) (*ftp.ServerConn, error) {
	return ftp.Dial(addr, options...)
}

func (fcw *FtpConnectionWrapper) ChangeDir(path string) error {
	return fcw.conn.ChangeDir(path)
}

func (fcw *FtpConnectionWrapper) Login(user string, password string) error {
	return fcw.conn.Login(user, password)
}

func (fcw *FtpConnectionWrapper) List(path string) (entries []*ftp.Entry, err error) {
	return fcw.conn.List(path)
}

func (fcw *FtpConnectionWrapper) Logout() error {
	return fcw.conn.Logout()
}

func (fcw *FtpConnectionWrapper) Quit() error {
	return fcw.conn.Quit()
}

func (fcw *FtpConnectionWrapper) Retr(path string) (*ftp.Response, error) {
	return fcw.conn.Retr(path)
}

func (fcw *FtpConnectionWrapper) Stor(path string, r io.Reader) error {
	return fcw.conn.Stor(path, r)
}

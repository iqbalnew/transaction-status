package ftp

import (
	"bytes"
	"errors"
	"io/ioutil"
	"time"

	ftpwrapper "bitbucket.bri.co.id/scm/bricams-addons/transaction-status/server/lib/ftp/wrapper"

	"github.com/jlaffaye/ftp"
)

type Config struct {
	User     string
	Password string
	Host     string
	Timeout  int64
}

type FtpObj struct {
	config Config
	fw     ftpwrapper.FtpInterface
	conn   ftpwrapper.FtpConnectionInterface
}

func NewFtp(config Config, fw ftpwrapper.FtpInterface, conn ftpwrapper.FtpConnectionInterface) *FtpObj {
	return &FtpObj{config: config, fw: fw, conn: conn}
}

func (fo *FtpObj) InitConnection() error {
	var err error
	if fo.conn == nil {
		fo.conn, err = fo.fw.Dial(fo.config.Host, ftp.DialWithTimeout(time.Duration(fo.config.Timeout)*time.Second))
		if err != nil {
			return err
		}
	}

	return nil
}

func (fo *FtpObj) Connect() error {
	return fo.conn.Login(fo.config.User, fo.config.Password)
}

func (fo *FtpObj) Close() error {
	if fo.conn != nil {
		fo.conn.Logout()
		fo.conn.Quit()
		fo.conn = nil
		return nil
	}
	return errors.New("connection already closed")
}

func (fo *FtpObj) ListFiles(path string) ([]*ftp.Entry, error) {
	return fo.conn.List(path)
}

func (fo *FtpObj) ReadFile(path string, filename string) ([]byte, error) {
	changeErr := fo.conn.ChangeDir(path)
	if changeErr != nil {
		return nil, changeErr
	}

	r, err := fo.conn.Retr(filename)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	return ioutil.ReadAll(r)
}

func (fo *FtpObj) WriteFile(path string, fileName string, content []byte) error {
	changeErr := fo.conn.ChangeDir(path)
	if changeErr != nil {
		return changeErr
	}

	err := fo.conn.Stor(fileName, bytes.NewBuffer(content))
	if err != nil {
		return err
	}

	return nil
}

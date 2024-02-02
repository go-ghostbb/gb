package gbtcp

import (
	"bufio"
	"bytes"
	"crypto/tls"
	gberror "ghostbb.io/gb/errors/gb_error"
	"io"
	"net"
	"time"
)

// Conn is the TCP connection object.
type Conn struct {
	net.Conn                     // Underlying TCP connection object.
	reader         *bufio.Reader // Buffer reader for connection.
	deadlineRecv   time.Time     // Timeout point for reading.
	deadlineSend   time.Time     // Timeout point for writing.
	bufferWaitRecv time.Duration // Interval duration for reading buffer.
}

const (
	// Default interval for reading buffer.
	receiveAllWaitTimeout = time.Millisecond
)

// NewConn creates and returns a new connection with given address.
func NewConn(addr string, timeout ...time.Duration) (*Conn, error) {
	if conn, err := NewNetConn(addr, timeout...); err == nil {
		return NewConnByNetConn(conn), nil
	} else {
		return nil, err
	}
}

// NewConnTLS creates and returns a new TLS connection
// with given address and TLS configuration.
func NewConnTLS(addr string, tlsConfig *tls.Config) (*Conn, error) {
	if conn, err := NewNetConnTLS(addr, tlsConfig); err == nil {
		return NewConnByNetConn(conn), nil
	} else {
		return nil, err
	}
}

// NewConnKeyCrt creates and returns a new TLS connection
// with given address and TLS certificate and key files.
func NewConnKeyCrt(addr, crtFile, keyFile string) (*Conn, error) {
	if conn, err := NewNetConnKeyCrt(addr, crtFile, keyFile); err == nil {
		return NewConnByNetConn(conn), nil
	} else {
		return nil, err
	}
}

// NewConnByNetConn creates and returns a TCP connection object with given net.Conn object.
func NewConnByNetConn(conn net.Conn) *Conn {
	return &Conn{
		Conn:           conn,
		reader:         bufio.NewReader(conn),
		deadlineRecv:   time.Time{},
		deadlineSend:   time.Time{},
		bufferWaitRecv: receiveAllWaitTimeout,
	}
}

// Send writes data to remote address.
func (c *Conn) Send(data []byte, retry ...Retry) error {
	for {
		if _, err := c.Write(data); err != nil {
			// Connection closed.
			if err == io.EOF {
				return err
			}
			// Still failed even after retrying.
			if len(retry) == 0 || retry[0].Count == 0 {
				err = gberror.Wrap(err, `Write data failed`)
				return err
			}
			if len(retry) > 0 {
				retry[0].Count--
				if retry[0].Interval == 0 {
					retry[0].Interval = defaultRetryInternal
				}
				time.Sleep(retry[0].Interval)
			}
		} else {
			return nil
		}
	}
}

// Recv receives and returns data from the connection.
//
// Note that,
//  1. If length = 0, which means it receives the data from current buffer and returns immediately.
//  2. If length < 0, which means it receives all data from connection and returns it until no data
//     from connection. Developers should notice the package parsing yourself if you decide receiving
//     all data from buffer.
//  3. If length > 0, which means it blocks reading data from connection until length size was received.
//     It is the most commonly used length value for data receiving.
func (c *Conn) Recv(length int, retry ...Retry) ([]byte, error) {
	var (
		err        error  // Reading error.
		size       int    // Reading size.
		index      int    // Received size.
		buffer     []byte // Buffer object.
		bufferWait bool   // Whether buffer reading timeout set.
	)
	if length > 0 {
		buffer = make([]byte, length)
	} else {
		buffer = make([]byte, defaultReadBufferSize)
	}

	for {
		if length < 0 && index > 0 {
			bufferWait = true
			if err = c.SetReadDeadline(time.Now().Add(c.bufferWaitRecv)); err != nil {
				err = gberror.Wrap(err, `SetReadDeadline for connection failed`)
				return nil, err
			}
		}
		size, err = c.reader.Read(buffer[index:])
		if size > 0 {
			index += size
			if length > 0 {
				// It reads til `length` size if `length` is specified.
				if index == length {
					break
				}
			} else {
				if index >= defaultReadBufferSize {
					// If it exceeds the buffer size, it then automatically increases its buffer size.
					buffer = append(buffer, make([]byte, defaultReadBufferSize)...)
				} else {
					// It returns immediately if received size is lesser than buffer size.
					if !bufferWait {
						break
					}
				}
			}
		}
		if err != nil {
			// Connection closed.
			if err == io.EOF {
				break
			}
			// Re-set the timeout when reading data.
			if bufferWait && isTimeout(err) {
				if err = c.SetReadDeadline(c.deadlineRecv); err != nil {
					err = gberror.Wrap(err, `SetReadDeadline for connection failed`)
					return nil, err
				}
				err = nil
				break
			}
			if len(retry) > 0 {
				// It fails even it retried.
				if retry[0].Count == 0 {
					break
				}
				retry[0].Count--
				if retry[0].Interval == 0 {
					retry[0].Interval = defaultRetryInternal
				}
				time.Sleep(retry[0].Interval)
				continue
			}
			break
		}
		// Just read once from buffer.
		if length == 0 {
			break
		}
	}
	return buffer[:index], err
}

// RecvLine reads data from the connection until reads char '\n'.
// Note that the returned result does not contain the last char '\n'.
func (c *Conn) RecvLine(retry ...Retry) ([]byte, error) {
	var (
		err    error
		buffer []byte
		data   = make([]byte, 0)
	)
	for {
		buffer, err = c.Recv(1, retry...)
		if len(buffer) > 0 {
			if buffer[0] == '\n' {
				data = append(data, buffer[:len(buffer)-1]...)
				break
			} else {
				data = append(data, buffer...)
			}
		}
		if err != nil {
			break
		}
	}
	return data, err
}

// RecvTill reads data from the connection until reads bytes `til`.
// Note that the returned result contains the last bytes `til`.
func (c *Conn) RecvTill(til []byte, retry ...Retry) ([]byte, error) {
	var (
		err    error
		buffer []byte
		data   = make([]byte, 0)
		length = len(til)
	)
	for {
		buffer, err = c.Recv(1, retry...)
		if len(buffer) > 0 {
			if length > 0 &&
				len(data) >= length-1 &&
				buffer[0] == til[length-1] &&
				bytes.EqualFold(data[len(data)-length+1:], til[:length-1]) {
				data = append(data, buffer...)
				break
			} else {
				data = append(data, buffer...)
			}
		}
		if err != nil {
			break
		}
	}
	return data, err
}

// RecvWithTimeout reads data from the connection with timeout.
func (c *Conn) RecvWithTimeout(length int, timeout time.Duration, retry ...Retry) (data []byte, err error) {
	if err = c.SetDeadlineRecv(time.Now().Add(timeout)); err != nil {
		return nil, err
	}
	defer func() {
		_ = c.SetDeadlineRecv(time.Time{})
	}()
	data, err = c.Recv(length, retry...)
	return
}

// SendWithTimeout writes data to the connection with timeout.
func (c *Conn) SendWithTimeout(data []byte, timeout time.Duration, retry ...Retry) (err error) {
	if err = c.SetDeadlineSend(time.Now().Add(timeout)); err != nil {
		return err
	}
	defer func() {
		_ = c.SetDeadlineSend(time.Time{})
	}()
	err = c.Send(data, retry...)
	return
}

// SendRecv writes data to the connection and blocks reading response.
func (c *Conn) SendRecv(data []byte, length int, retry ...Retry) ([]byte, error) {
	if err := c.Send(data, retry...); err == nil {
		return c.Recv(length, retry...)
	} else {
		return nil, err
	}
}

// SendRecvWithTimeout writes data to the connection and reads response with timeout.
func (c *Conn) SendRecvWithTimeout(data []byte, length int, timeout time.Duration, retry ...Retry) ([]byte, error) {
	if err := c.Send(data, retry...); err == nil {
		return c.RecvWithTimeout(length, timeout, retry...)
	} else {
		return nil, err
	}
}

// SetDeadline sets the deadline for current connection.
func (c *Conn) SetDeadline(t time.Time) (err error) {
	if err = c.Conn.SetDeadline(t); err == nil {
		c.deadlineRecv = t
		c.deadlineSend = t
	}
	if err != nil {
		err = gberror.Wrapf(err, `SetDeadline for connection failed with "%s"`, t)
	}
	return err
}

// SetDeadlineRecv sets the deadline of receiving for current connection.
func (c *Conn) SetDeadlineRecv(t time.Time) (err error) {
	if err = c.SetReadDeadline(t); err == nil {
		c.deadlineRecv = t
	}
	if err != nil {
		err = gberror.Wrapf(err, `SetDeadlineRecv for connection failed with "%s"`, t)
	}
	return err
}

// SetDeadlineSend sets the deadline of sending for current connection.
func (c *Conn) SetDeadlineSend(t time.Time) (err error) {
	if err = c.SetWriteDeadline(t); err == nil {
		c.deadlineSend = t
	}
	if err != nil {
		err = gberror.Wrapf(err, `SetDeadlineSend for connection failed with "%s"`, t)
	}
	return err
}

// SetBufferWaitRecv sets the buffer waiting timeout when reading all data from connection.
// The waiting duration cannot be too long which might delay receiving data from remote address.
func (c *Conn) SetBufferWaitRecv(bufferWaitDuration time.Duration) {
	c.bufferWaitRecv = bufferWaitDuration
}

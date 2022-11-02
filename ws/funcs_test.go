package ws

import (
	"context"
	"testing"
	"time"
)

func TestSingleConn_writePump(t *testing.T) {
	s := &SingleConn{
		conn:                    nil,
		id:                      "",
		ctx:                     context.TODO(),
		cancel:                  nil,
		options:                 nil,
		closed:                  false,
		beforeHandleReceivedMsg: nil,
		handleReceiveMsg:        nil,
		afterHandleReceivedMsg:  nil,
		handleReceiveTaskErrors: nil,
		isReading:               false,
		beforeHandleSendMsg:     nil,
		afterHandleSendMsg:      nil,
		handleSendTaskErrors:    nil,
		isWriting:               false,
		sendChan:                make(chan Msg, 1),
		heartCheck:              3 * time.Second,
		sendTimeOut:             5 * time.Second,
		writeTimeOut:            0,
	}
	//s.writePump()
	s.Serve()

	time.Sleep(time.Hour)
}

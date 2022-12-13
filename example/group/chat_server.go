package main

import (
	"context"
	"git.woa.com/nioliu/wsutil-go/group"
	"git.woa.com/nioliu/wsutil-go/ws"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"log"
	"time"
)

var g *group.Group

func main() {
	var err error
	c := context.Background()
	// init root group
	g, err = group.NewWithContext(c, group.WithUpgrader(ws.NewWrappedGorillaUpgrader()), group.WithMaxConnCnt(10))
	if err != nil {
		log.Fatal(err)
	}

	// init gin
	engine := gin.New()

	// handle ...
	engine.Handle("GET", "/add", addUserToGroup)

	engine.Handle("POST", "/msg", broadcastMsg)

	engine.Handle("POST", "/msg/tags", sendMsgToTag)

	if err := engine.Run("0.0.0.0:9090"); err != nil {
		log.Fatal(err)
	}
}

func addUserToGroup(ctx *gin.Context) {
	conn, err := g.WsUpgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println(zap.Error(err))
		return
	}
	tag := ctx.Query("tag")
	singleConn, err := ws.NewSingleConn(ctx, conn, ws.WithContext(ctx),
		ws.WithHeartCheck(time.Second*10),
		ws.WithReceiveTaskErrors(func(ctx context.Context, id string, err []error) error {
			log.Println(ctx)
			log.Println(id)
			log.Println(err)
			return err[len(err)-1]
		}),
		ws.WithHandleReceiveMsg(
			func(ctx context.Context, id string, msgType int, msg []byte, err []error) error {
				if err != nil {
					log.Println(err)
					return err[0]
				}
				switch msgType {
				case websocket.BinaryMessage:
					log.Println("this is a binary msg: ", string(msg))
				case websocket.TextMessage:
					log.Println("this is a text message: ", string(msg))
				}
				return nil
			}),
		ws.WithTags(tag))
	if err != nil {
		log.Println(zap.Error(err))
		return
	}

	if err = g.AddNewSingleConn(singleConn); err != nil {
		if err := singleConn.Close(); err != nil {
			log.Println(err)
		}
		log.Println(err)
		return
	}

	log.Println(ctx, "add new single conn,id:"+singleConn.GetId())

}

func broadcastMsg(c *gin.Context) {
	msg := c.Query("msg")
	log.Println(msg)
	if err := g.Broadcast(c, ws.Msg{
		Msg:     []byte(msg),
		MsgType: websocket.TextMessage,
	}); err != nil {
		log.Println(err)
		return
	}
}

func sendMsgToTag(c *gin.Context) {
	msg := c.Query("msg")
	tag := c.Query("tag")
	log.Println(msg)
	if err := g.SendMsgWithTags(c, ws.Msg{Msg: []byte(msg), MsgType: websocket.TextMessage},
		true, tag); err != nil {
		log.Println(err)
		return
	}
}

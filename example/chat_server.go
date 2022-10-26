package main

import (
	"context"
	"git.woa.com/nioliu/wsutil-go/group"
	"git.woa.com/nioliu/wsutil-go/ws"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"log"
	"net/http"
	"sync"
	"time"
)

var once = sync.Once{}

func main() {
	engine := gin.New()
	c := context.Background()
	g, err := group.NewGroupWithContext(c, &ws.WrappedGorillaUpgrader{}, group.WithMaxConnCnt(10))
	if err != nil {
		log.Println(err)
		return
	}
	engine.Handle("GET", "/add", func(ctx *gin.Context) {
		upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
			return true
		}}
		conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			log.Println(zap.Error(err))
			return
		}
		singleConn, err := ws.NewSingleConn(ctx, conn, ws.WithContext(ctx),
			ws.WithHeartCheck(time.Second*10), ws.WithReceiveTaskErrors(func(ctx context.Context, id string, err []error) error {
				log.Println(ctx)
				log.Println(id)
				log.Println(err)
				return err[len(err)-1]
			}), ws.WithHandleReceiveMsg(
				func(ctx context.Context, id string, msgType int, msg []byte, err []error) error {
					if err != nil {
						log.Println(err)
					}
					switch msgType {
					case websocket.BinaryMessage:
						log.Println("this is a binary msg: ", string(msg))
					case websocket.TextMessage:
						log.Println("this is a text message: ", string(msg))
					}
					return nil
				}))
		if err != nil {
			log.Println(zap.Error(err))
			return
		}

		if err = singleConn.Serve(); err != nil {
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
	})

	engine.Handle("POST", "/msg", func(c *gin.Context) {
		msg := c.Query("msg")
		log.Println(msg)
		if err := g.Broadcast(c, ws.Msg{
			Msg:     []byte(msg),
			MsgType: websocket.TextMessage,
		}); err != nil {
			log.Println(err)
			return
		}
	})

	if err := engine.Run("0.0.0.0:9090"); err != nil {
		log.Fatal(err)
	}
}

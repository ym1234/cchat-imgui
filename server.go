package main

import (
	// "fmt"

	"github.com/diamondburned/cchat/text"
	"context"
	"time"
	"github.com/inkyblackness/imgui-go/v2"
	"github.com/diamondburned/cchat"
)

type Message struct {
	ID string
	content text.Rich
	time time.Time
	author text.Rich
}

type Server struct {
	server  cchat.Server // either ServerMessage or ServerList
	isActive bool
	messages []*Message
	messagesMap map[string]*Message
	servers []*Server
	main chan func() error
}


func (sers *Server) Start() {
	if x, ok := sers.server.(cchat.ServerList); ok {
		go x.Servers(sers)
	}
}

func (sers *Server) SetServers(newservers []cchat.Server) {
	sers.main <- func() error {
		sers.servers = make([]*Server, len(newservers))
		for i := 0; i < len(newservers); i++ {
			sers.servers[i] = &Server{server: newservers[i], main: sers.main}
			if x, ok := newservers[i].(cchat.ServerList); ok {
				go x.Servers(sers.servers[i])
			}
		}
		return nil
	}
}

func (sers *Server) CreateMessage(info cchat.MessageCreate) {
	sers.main <- func() error {
		msg := &Message{ID: info.ID(), content: info.Content(), author: info.Author().Name(), time: info.Time()}
		sers.messages = append(sers.messages, msg)
		if sers.messagesMap == nil {
			sers.messagesMap = make(map[string]*Message)
		}
		sers.messagesMap[info.ID()] = msg
		return nil
	}
}

func (sers *Server) UpdateMessage(info cchat.MessageUpdate) {


}


func (sers *Server) DeleteMessage(info cchat.MessageDelete) {


}

func (sers *Server) Render() {
	if _, ok := sers.server.(cchat.ServerList); ok {
		imgui.PushID(sers.server.Name().Content + sers.server.ID())
		if imgui.TreeNode(sers.server.Name().Content) {
			for _, k := range sers.servers {
				k.Render()
			}
			if len(sers.servers) == 0 {
				imgui.Text("No channels here!")
			}
			imgui.TreePop()
		}
		imgui.PopID()
	} else {
		x, ok := sers.server.(cchat.ServerMessage)
		if !ok {
		}
		if imgui.Selectable(sers.server.Name().Content)  && !sers.isActive {
			sers.isActive = true
			go x.JoinServer(context.Background(), sers)
		}
		if !sers.isActive  {
			return
		}
		if imgui.Begin(sers.server.Name().Content) {
			x := imgui.WindowSize()
			imgui.PushTextWrapPosV(x.X)
			for _, msg := range sers.messages {
				msg.Render()
			}
			imgui.PopTextWrapPos()
			imgui.SetScrollY(imgui.ScrollMaxY())
		}

		imgui.End()
	}
}


// TODO
func (msg *Message) Render() {
	// Sort so that all ending points are sorted decrementally. We probably
	// don't need SliceStable here, as we're sorting again.
	// sort.Slice(content.Segments, func(i, j int) bool {
	// 	_, i = content.Segments[i].Bounds()
	// 	_, j = content.Segments[j].Bounds()
	// 	return i > j
	// })

	// // Sort so that all starting points are sorted incrementally.
	// sort.SliceStable(content.Segments, func(i, j int) bool {
	// 	i, _ = content.Segments[i].Bounds()
	// 	j, _ = content.Segments[j].Bounds()
	// 	return i < j
	// })
	// content := msg.author.Content
	// for _, k := range msg.author.Segments {
	// 	start, end := k.Bounds()
	// 	if x, ok := k.(text.Attributor); ok {
	// 		attributes := x.Attribute()
	// 		result := 0
	// 		if attributes.Has(text.AttrBold) {

	// 		}
	// 		if attributes.Has(text.AttrItalics) {

	// 		}
	// 		if attributes.Has(text.AttrUnderline) {

	// 		}
	// 		if attributes.Has(text.AttrStrikethrough) {

	// 		}
	// 		if attributes.Has(text.AttrSpoiler) {

	// 		}
	// 		if attributes.Has(text.AttrMonospace) {

	// 		}
	// 		if attributes.Has(text.AttrDimmed) {

	// 		}
	// 	}
	// }
}

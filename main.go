package main

import (
	"fmt"
	// "github.com/davecgh/go-spew/spew"
	"os"

	_ "github.com/diamondburned/cchat-mock"
	_ "github.com/diamondburned/cchat-discord"
	"github.com/diamondburned/cchat/services"
	"github.com/inkyblackness/imgui-go/v2"
	"github.com/ym1234/cchat-imgui/internal"
)

func Run(p Platform, r Renderer) {
	imgui.CurrentIO().SetClipboard(clipboard{platform: p})

	currentServices, err := services.Get()
	if len(err) != 0 {
		for _, k := range err {
			print(k.Error())
		}
		os.Exit(-1)
	}

	myServices := make([]*Service, len(currentServices))
	for i, service := range currentServices {
		ldi := &LoginDialogInfo{
			service: service,
			auth: service.Authenticate(),
			loginInfo: make([]string, len(service.Authenticate().AuthenticateForm())),
			lastErr: nil,
			isActive: false,
			isPending: false,
			newSession: make(chan newSessionMessage)}

		myServices[i] = &Service{service, []*Server{}, ldi}
	}


	execOnMain := make(chan func() error)

	for !p.ShouldStop() {
		p.ProcessEvents()
		p.NewFrame()
		imgui.NewFrame()
		RenderNormalUI(myServices)
		for _, k := range myServices {
			newSession := k.ldi.Render()
			if newSession == nil {
				continue
			}
			session := &Server{server: newSession, servers: nil, main: execOnMain}
			session.Start()
			k.sessions = append(k.sessions, session)
		}

		// TODO check if this can
		select {
		case f := <- execOnMain:
			err := f()
			if err != nil  {
				println(err.Error())
			}
		default:
			break
		}

		imgui.Render()
		r.PreRender([3]float32{0.4, 0.6, 0.7})
		r.Render(p.DisplaySize(), p.FramebufferSize(), imgui.RenderedDrawData())
		p.PostRender()
	}
}

func RenderNormalUI(services []*Service) {
	if imgui.BeginV("Servers", nil, 0) {
		// spew.Println(services)
		for _, k :=  range services {
			// imgui.PushID(k.service.Name().Content + k.service.ID())
			open := imgui.TreeNodeV(k.service.Name().Content,  0)
			// if imgui.BeginPopupContextItemV(k.service.Name().Content, 1) {
			// 	if imgui.Selectable("Connect") {
			// 		k.ldi.isActive = true
			// 	}
			// 	if imgui.Selectable("Disconnect") {

			// 	}
			// 	imgui.EndPopup()
			// }

			if open {
				imgui.Separator()
				if imgui.Selectable("Connect") {
					k.ldi.isActive = true
				}
				if imgui.Selectable("Disconnect All") {
					// TODO
				}
				imgui.Unindent()

				for _, j := range k.sessions {
					j.Render()
				}
				imgui.Indent()
				imgui.Separator()
				imgui.TreePop()
			}
			// imgui.PopID()
		}
	}
	imgui.End()

}

func main() {
	context := imgui.CreateContext(nil)
	defer context.Destroy()
	io := imgui.CurrentIO()

	platform, err := internal.NewGLFW(io, internal.GLFWClientAPIOpenGL3)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(-1)
	}
	defer platform.Dispose()

	renderer, err := internal.NewOpenGL3(io)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(-1)
	}
	defer renderer.Dispose()

	Run(platform, renderer)
}

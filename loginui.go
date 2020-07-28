package main
import (
	"github.com/diamondburned/cchat"
	"context"
	"errors"

	"github.com/inkyblackness/imgui-go/v2"
)
type LoginDialogInfo struct {
	service    cchat.Service
	auth       cchat.Authenticator
	loginInfo  []string
	lastErr    error
	isActive   bool
	isPending  bool
	newSession chan newSessionMessage
	cancel context.CancelFunc
}

type newSessionMessage struct {
	Session cchat.Session
	err     error
}


func (ldi *LoginDialogInfo) Render() (session cchat.Session) {
	if !ldi.isActive {
		return
	}

	if ldi.isPending {
		select {
		case x := <-ldi.newSession:
			if x.err == nil {
				ldi.Reset()
				session = x.Session
				return
			} else {
				ldi.lastErr = x.err
				ldi.isPending = false
			}
		default:
			imgui.SetNextWindowSize(imgui.Vec2{500, 250})
			if imgui.Begin(ldi.service.Name().Content + " - Login") {
				RenderForm(ldi.auth.AuthenticateForm(), ldi.loginInfo, true)
				imgui.Separator()
				imgui.PushStyleVarFloat(imgui.StyleVarAlpha, 0.6)
				imgui.Button("Confirm")
				imgui.PopStyleVar()
				imgui.SameLine()
				if imgui.Button("Cancel") {
					ldi.cancel()
					ldi.cancel = nil
					ldi.isPending = false
					ldi.lastErr = errors.New("Canceled!")
				}
				imgui.End()
			}
			return
		}
	}

	imgui.SetNextWindowSize(imgui.Vec2{500, 250})
	if imgui.Begin(ldi.service.Name().Content + " - Login") {
		enterWasPressed := RenderForm(ldi.auth.AuthenticateForm(), ldi.loginInfo, false)
		imgui.Separator()
		if imgui.Button("Confirm")  || enterWasPressed {
			ldi.isPending = true
			ctx, cancel := context.WithCancel(context.Background())
			go func() {
				session, err := ldi.auth.Authenticate(ldi.loginInfo)
				select {
				case <- ctx.Done():
					return
				default:
					ldi.newSession <- newSessionMessage{session, err}
				}
			}()
			ldi.cancel = cancel
		}
		imgui.SameLine()
		if imgui.Button("Close") {
			ldi.Reset()
		}
		if ldi.lastErr != nil {
			imgui.PushTextWrapPosV(500)
			imgui.Text(ldi.lastErr.Error())
		}
	}
	imgui.End()
	return
}

func (ldi *LoginDialogInfo) Reset() {
	ldi.isActive = false
	ldi.isPending = false
	ldi.lastErr = nil
	for i := 0; i < len(ldi.loginInfo); i++ {
		ldi.loginInfo[i] = ""
	}
}

func RenderForm(form []cchat.AuthenticateEntry, loginInfo []string, readonly bool) bool {
	imgui.ColumnsV(2, "", false)
	imgui.SetColumnWidth(0, 150)
	enterWasPressed := false
	for i, k := range form {
		flags := 0
		if readonly {
			imgui.PushStyleVarFloat(imgui.StyleVarAlpha, 0.6)
			flags |= imgui.InputTextFlagsReadOnly
		}
		if k.Secret {
			flags |= imgui.InputTextFlagsPassword
		}

		imgui.Text(k.Name + ":")
		imgui.NextColumn()

		if k.Multiline {
			enterWasPressed = imgui.InputTextMultilineV(" ##"+k.Name, &loginInfo[i], imgui.Vec2{0, 0}, flags | imgui.InputTextFlagsEnterReturnsTrue, nil) || enterWasPressed
		} else {
			enterWasPressed = imgui.InputTextV(" ##"+k.Name, &loginInfo[i], flags | imgui.InputTextFlagsEnterReturnsTrue, nil) || enterWasPressed
		}
		imgui.NextColumn()
		if readonly {
			imgui.PopStyleVar()
		}
	}
	imgui.Columns()
	return enterWasPressed
}


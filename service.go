package main

import (
	"github.com/diamondburned/cchat"
)

type Service struct {
	service cchat.Service
	sessions []*Server
	ldi *LoginDialogInfo
}


package m_spew

import "github.com/davecgh/go-spew/spew"

var Config = spew.ConfigState{Indent: "\t", DisableMethods: true}

func Dump(o ...any) {
	Config.Dump(o...)
}

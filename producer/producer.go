package producer

import "github.com/blackjack200/rich-go-plus/codec"

type Producer interface {
	Activity() *codec.Activity
	AppId() string
}

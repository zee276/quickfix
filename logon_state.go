package quickfix

import (
	"github.com/quickfixgo/quickfix/fix/tag"
	"github.com/quickfixgo/quickfix/message"
)

type logonState struct{}

func (s logonState) FixMsgIn(session *Session, msg message.Message) (nextState sessionState) {
	msgType := new(message.StringValue)
	if err := msg.Header.GetField(tag.MsgType, msgType); err == nil && msgType.Value == "A" {
		if err := session.handleLogon(msg); err != nil {
			session.log.OnEvent(err.Error())
			return latentState{}
		}

		return inSession{}
	}

	session.log.OnEventf("Invalid Session State: Received Msg %v while waiting for Logon", msg)
	return latentState{}
}

func (s logonState) Timeout(session *Session, e event) (nextState sessionState) {
	if e == logonTimeout {
		session.log.OnEvent("Timed out waiting for logon response")
		return latentState{}
	}

	return s
}

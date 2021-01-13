package detailedlogger

import (
	"bytes"

	"golang.org/x/crypto/ssh"
)

const (
	// https://tools.ietf.org/html/rfc4254#section-9
	sshGlobalRequest         uint8 = 80
	sshGlobalRequestSuccess  uint8 = 81
	sshGlobalRequestFailure  uint8 = 82
	sshChannelOpen           uint8 = 90 // A request from the client or the server to create a channel
	sshChannelOpenSuccess    uint8 = 91
	sshChannelOpenFailure    uint8 = 92
	sshChannelWindowAdjust   uint8 = 93
	sshChannelData           uint8 = 94
	sshChannelExtendedData   uint8 = 95
	sshChannelEOF            uint8 = 96
	sshChannelClose          uint8 = 97
	sshChannelRequest        uint8 = 98 // A request from the client or server specific to this channel
	sshChannelRequestSuccess uint8 = 99
	sshChannelRequestFailure uint8 = 100
)

// Audit messages sent from the server to the client
func (a *DetailedConnectionAuditor) sshServerToClient(conn ssh.ConnMetadata, msg []byte) ([]byte, error) {
	stream := &message{bytes.NewBuffer(msg)}

	var msgType uint8
	msgType, _ = stream.byte()

	switch msgType {
	case sshChannelData:
		serverChan, _ := stream.uint32()
		t, ok := a.serverChannels[serverChan]
		if !ok {
			break
		}
		stream.uint32() // shift off the length of the rest of the data packet
		switch t.session {
		case channelSessionTypeShell:
		case channelSessionTypeExec:
		case channelSessionTypeSFTP:
			// we currently are not logging anything sent from the server to the client
			// a.sftpServerToClient(t, stream)
		}
	case sshChannelExtendedData:

	case sshChannelOpenSuccess:
		// https://tools.ietf.org/html/rfc4254#section-5.1
		clientChan, _ := stream.uint32() // Client first since this message here means the client requested the channel
		serverChan, _ := stream.uint32()
		ch := &channel{
			clientID: clientChan,
			serverID: serverChan,
			session:  channelSessionTypeUnknown,
		}
		ch.init(a)
		a.clientChannels[clientChan] = ch
		a.serverChannels[serverChan] = ch
	case sshChannelOpen: // Ignored, since we only account channels on a success response from the server, or client
	case sshChannelOpenFailure: // Ignored, since we only account channels on a success response from the server, or client

	case sshChannelRequest:
		// The server us making a channel request to the client
		serverChan, _ := stream.uint32()
		ch, ok := a.serverChannels[serverChan]
		if !ok {
			break
		}
		req, _ := stream.str()
		switch {
		case bytes.Equal(req, []byte("subsystem")):
			stream.bool()
			data, _ := stream.str()
			switch {
			case bytes.Equal(data, subsystemNameSFTP):
				ch.session = channelSessionTypeRequestedSFTP
				ch.Params["subsystem"] = "sftp"
			}
		case bytes.Equal(req, subsystemNameExec):
			stream.bool()
			cmd, _ := stream.str()
			ch.session = channelSessionTypeRequestedExec
			ch.Params["subsystem"] = "exec"
			ch.Params["command"] = string(cmd)
		case bytes.Equal(req, subsystemNameShell):
			ch.session = channelSessionTypeRequestedShell
			ch.Params["subsystem"] = "shell"
		}
	case sshChannelRequestSuccess:
		// The server is responding to a channel request from the client
		clientChan, _ := stream.uint32()
		ch, ok := a.clientChannels[clientChan]
		if !ok {
			break
		}
		switch ch.session {
		case channelSessionTypeRequestedSFTP:
			ch.session = channelSessionTypeSFTP
			ch.log(1, "started")
		case channelSessionTypeRequestedExec:
			ch.session = channelSessionTypeExec
			ch.log(1, "exec: %s", ch.Params["command"])
		case channelSessionTypeRequestedShell:
			ch.session = channelSessionTypeShell
			ch.log(1, "started")
		}
	case sshChannelRequestFailure:
		clientChan, _ := stream.uint32()
		ch, ok := a.clientChannels[clientChan]
		if !ok {
		}
		switch ch.session {
		case channelSessionTypeRequestedSFTP:
			ch.log(1, "failed")
		case channelSessionTypeRequestedExec:
			ch.log(1, "failed")
		case channelSessionTypeRequestedShell:
			ch.log(1, "failed")
		}

	case sshGlobalRequest: // Ignored, so far we have no need of this
	case sshGlobalRequestSuccess: // Ignored, so far we have no need of this
	case sshGlobalRequestFailure: // Ignored, so far we have no need of this

	case sshChannelEOF:
	case sshChannelClose:
	case sshChannelWindowAdjust:
	}
	return msg, nil
}

// Audit messages sent from the client to the server
func (a *DetailedConnectionAuditor) sshClientToServer(conn ssh.ConnMetadata, msg []byte) ([]byte, error) {
	stream := &message{bytes.NewBuffer(msg)}

	var msgType uint8
	msgType, _ = stream.byte()

	switch msgType {
	case sshChannelData:
		clientChan, _ := stream.uint32()
		t, ok := a.serverChannels[clientChan]
		if !ok {
			break
		}
		stream.uint32() // shift off the length of the rest of the data packet
		switch t.session {
		case channelSessionTypeExec:
		case channelSessionTypeShell:
		case channelSessionTypeSFTP:
			a.sftpClientToServer(t, stream)
			break
		}
	case sshChannelExtendedData:

	case sshChannelOpenSuccess:
		// https://tools.ietf.org/html/rfc4254#section-5.1
		serverChan, _ := stream.uint32() // Server first since this message here means the server requested the channel
		clientChan, _ := stream.uint32() // Client second
		ch := &channel{
			clientID: clientChan,
			serverID: serverChan,
			session:  channelSessionTypeUnknown,
		}
		ch.init(a)
		a.clientChannels[clientChan] = ch
		a.serverChannels[serverChan] = ch
	case sshChannelOpen: // Ignored, since we only account channels on a success response from the server, or client
	case sshChannelOpenFailure: // Ignored, since we only account channels on a success response from the server, or client

	case sshChannelRequest:
		// The client is making a channel request to the server
		clientChan, _ := stream.uint32()
		ch, ok := a.clientChannels[clientChan]
		if !ok {
			break
		}
		req, _ := stream.str()
		switch {
		case bytes.Equal(req, []byte("subsystem")):
			stream.bool()
			data, _ := stream.str()
			switch {
			case bytes.Equal(data, subsystemNameSFTP):
				ch.Params["subsystem"] = "sftp"
				ch.session = channelSessionTypeRequestedSFTP
			}
		case bytes.Equal(req, subsystemNameExec):
			stream.bool()
			cmd, _ := stream.str()
			ch.session = channelSessionTypeRequestedExec
			ch.Params["subsystem"] = "exec"
			ch.Params["command"] = string(cmd)
		case bytes.Equal(req, subsystemNameShell):
			ch.Params["subsystem"] = "shell"
			ch.session = channelSessionTypeRequestedShell
		}
	case sshChannelRequestSuccess:
		// The client is responding to a channel request sent by the server
		serverChan, _ := stream.uint32()
		ch, ok := a.serverChannels[serverChan]
		if !ok {
			break
		}
		switch ch.session {
		case channelSessionTypeRequestedSFTP:
			ch.session = channelSessionTypeSFTP
		case channelSessionTypeRequestedExec:
			ch.session = channelSessionTypeExec
		case channelSessionTypeRequestedShell:
			ch.session = channelSessionTypeShell
		}
	case sshChannelRequestFailure:

	case sshGlobalRequest: // Ignored, so far we have no need of this
	case sshGlobalRequestSuccess: // Ignored, so far we have no need of this
	case sshGlobalRequestFailure: // Ignored, so far we have no need of this

	case sshChannelEOF:
	case sshChannelClose:
	case sshChannelWindowAdjust:
	}
	return msg, nil
}

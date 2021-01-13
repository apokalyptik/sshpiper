package detailedlogger

import (
	"context"
	"strings"

	"golang.org/x/crypto/ssh"

	"github.com/tg123/sshpiper/sshpiperd/auditor"
)

// DetailedConnectionAuditor is the actual auditor/logger in this package and
// is able to parse and house necessary state for proper auditing and logging
type DetailedConnectionAuditor struct {
	*jsonLogger
	context         context.Context
	clientChannels  map[uint32]*channel
	serverChannels  map[uint32]*channel
	sftpFileHandles map[uint32]string
}

// AddContext adds a context to the current auditor
func (a *DetailedConnectionAuditor) AddContext(ctx context.Context) {
	a.context = ctx
	if user := ctx.Value(ssh.SSHPiperContextDownstreamUsername); user != nil {
		a.Params["downstream_ssh_user"] = user.(string)
	}
	if user := ctx.Value(ssh.SSHPiperContextUpstreamUsername); user != nil {
		a.Params["upstream_ssh_user"] = user.(string)
	}
	if host := ctx.Value(ssh.SSHPiperContextUpstreamAddress); host != nil {
		destAddr := host.(string)
		portIndex := strings.LastIndex(destAddr, ":")
		a.Params["destination_ip"] = destAddr[:portIndex]
		a.Params["destination_port"] = destAddr[portIndex+1:]
	}
	if host := ctx.Value(ssh.SSHPiperContextDownstreamAddress); host != nil {
		remoteAddr := host.(string)
		portIndex := strings.LastIndex(remoteAddr, ":")
		a.Params["source_ip"] = remoteAddr[:portIndex]
		a.Params["source_port"] = remoteAddr[portIndex+1:]
	}
}

func (a *DetailedConnectionAuditor) init() {
	a.clientChannels = map[uint32]*channel{}
	a.serverChannels = map[uint32]*channel{}
	a.sftpFileHandles = map[uint32]string{}
	a.jsonLogger, _ = newJsonLogger("sshpiper")
}

// GetUpstreamHook satisfies an interface for obtaining a generic hook to pass messages through
func (a *DetailedConnectionAuditor) GetUpstreamHook() auditor.Hook {
	return a.sshServerToClient
}

// GetDownstreamHook satisfies an interface for obtaining a generic hook to pass messages through
func (a *DetailedConnectionAuditor) GetDownstreamHook() auditor.Hook {
	return a.sshClientToServer
}

// Close Satisfies the interface requirements
func (a *DetailedConnectionAuditor) Close() error {
	a.clientChannels = nil
	a.serverChannels = nil
	a.sftpFileHandles = nil
	a.jsonLogger = nil
	return nil
}

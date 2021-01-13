package detailedlogger

type channelSessionType uint32

const (
	channelSessionTypeUnknown channelSessionType = iota
	channelSessionTypeRequestedSFTP
	channelSessionTypeSFTP
	channelSessionTypeRequestedExec
	channelSessionTypeExec
	channelSessionTypeRequestedShell
	channelSessionTypeShell
)

var (
	subsystemNameSFTP  = []byte("sftp")
	subsystemNameExec  = []byte("exec")
	subsystemNameShell = []byte("shell")
)

type channel struct {
	*jsonLogger
	clientID         uint32
	serverID         uint32
	session          channelSessionType
	sftpOpenRequests map[uint32]string
	sftpOpenHandles  map[uint32]string
	Loglevel         int
}

func (c *channel) init(a *DetailedConnectionAuditor) {
	c.sftpOpenHandles = map[uint32]string{}
	c.sftpOpenHandles = map[uint32]string{}
	c.jsonLogger = a.jsonLogger.clone()
}

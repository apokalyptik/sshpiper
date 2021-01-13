package detailedlogger

import "bytes"

const (
	// https://tools.ietf.org/html/draft-ietf-secsh-filexfer-02#section-3
	// This is, apparently, the version of the protocol that openssh's built-in subsystem impliments
	sftpInit          uint8 = 1
	sftpVersion       uint8 = 2
	sftpOpen          uint8 = 3
	sftpClose         uint8 = 4
	sftpRead          uint8 = 5
	sftpWrite         uint8 = 6
	sftpLstat         uint8 = 7
	sftpFstat         uint8 = 8
	sftpSetstat       uint8 = 9
	sftpFsetstat      uint8 = 10
	sftpOpendir       uint8 = 11
	sftpReadDir       uint8 = 12
	sftpRemove        uint8 = 13
	sftpMkdir         uint8 = 14
	sftpRmdir         uint8 = 15
	sftpRealpath      uint8 = 16
	sftpStat          uint8 = 17
	sftpRename        uint8 = 18
	sftpReadlink      uint8 = 19
	sftpSymlink       uint8 = 20
	sftpStatus        uint8 = 101
	sftpHandle        uint8 = 102
	sftpData          uint8 = 103
	sftpName          uint8 = 104
	sftpAttrs         uint8 = 105
	sftpExtended      uint8 = 200
	sftpExtendedReply uint8 = 201
)

// sftpServeTtoClient is not currently used. We are not tracking state right now (which would be necessary to, for example,
// announce that a file is being written to as opposed to being read from).  If we decide that we are no longer content
// knowing that a file was opened, but in face want to know whether it was read from or written to, then we would need to
// start tracking things like file handles and response statuses on them to then wait for a read or write op on a valid
// and known file handle.
func (a *DetailedConnectionAuditor) sftpServerToClient(ch *channel, stream *message) {
}

// sftpClientToServer handles messages passed from the client to the server. This is where we gain most of our intelligence
// from.  For limitations and future improvements if desired, please see the command attached to the currently unused
// sftpServeTtoClient function
func (a *DetailedConnectionAuditor) sftpClientToServer(ch *channel, stream *message) {
	stream.uint32()
	packetType, _ := stream.ReadByte()
	stream.uint32()
	switch packetType {
	// ignored for brevity
	case sftpInit:
	// ignored for brevity
	case sftpVersion:
	case sftpOpen:
		path, _ := stream.str()
		ch.Params["path"] = string(path)
		ch.Params["action"] = "openfile"
		ch.log(1, "openfile %s", ch.Params["path"])
	// ignored to avoid tracking more state
	case sftpClose:
	// ignored to avoid tracking more state
	case sftpRead:
	// ignored to avoid tracking more state
	case sftpWrite:
	// ignored for brevity
	case sftpLstat:
	// ignored for brevity
	case sftpFstat:
	// ignored for brevity
	case sftpSetstat:
	// ignored for brevity
	case sftpFsetstat:
	case sftpOpendir:
		path, _ := stream.str()
		ch.Params["path"] = string(path)
		ch.Params["action"] = "opendir"
		ch.log(1, "opendir %s", ch.Params["path"])
	// ignored for brevity
	case sftpReadDir:
	case sftpRemove:
		path, _ := stream.str()
		ch.Params["path"] = string(path)
		ch.Params["action"] = "rm"
		ch.log(1, "rm %s", ch.Params["path"])
	case sftpMkdir:
		path, _ := stream.str()
		ch.Params["path"] = string(path)
		ch.Params["action"] = "mkdir"
		ch.log(1, "mkdir %s", ch.Params["path"])
	case sftpRmdir:
		path, _ := stream.str()
		ch.Params["path"] = string(path)
		ch.Params["action"] = "rmdir"
		ch.log(1, "rmdir %s", ch.Params["path"])
	// ignored for brevity
	case sftpRealpath:
	// ignored for brevity
	case sftpStat:
	case sftpRename:
		from, _ := stream.str()
		to, _ := stream.str()
		ch.Params["action"] = "rename"
		ch.Params["path"] = string(from)
		ch.Params["target"] = string(to)
		ch.log(1, "rename %s to %s", ch.Params["path"], ch.Params["target"])
	// ignored for brevity
	case sftpReadlink:
	case sftpSymlink:
		from, _ := stream.str()
		to, _ := stream.str()
		ch.Params["action"] = "symlink"
		ch.Params["path"] = string(from)
		ch.Params["target"] = string(to)
		ch.log(1, "symlink %s to %s", ch.Params["path"], ch.Params["target"])
	// ignored to avoid tracking more state
	case sftpStatus:
	// ignored to avoid tracking more state
	case sftpHandle:
	// ignored for brevity
	case sftpData:
	// ignored for brevity
	case sftpName:
	// ignored for brevity
	case sftpAttrs:
	case sftpExtended:
		extension, _ := stream.str()
		switch {
		case bytes.Equal(extension, []byte("posix-rename@openssh.com")):
			from, _ := stream.str()
			to, _ := stream.str()
			ch.Params["action"] = "rename"
			ch.Params["path"] = string(from)
			ch.Params["target"] = string(to)
			ch.log(1, "rename %s to %s", ch.Params["path"], ch.Params["target"])
		case bytes.Equal(extension, []byte("hardlink@openssh.com")):
			from, _ := stream.str()
			to, _ := stream.str()
			ch.Params["action"] = "hardlink"
			ch.Params["path"] = string(from)
			ch.Params["target"] = string(to)
			ch.log(1, "hardlink %s to %s", ch.Params["path"], ch.Params["target"])
		case bytes.Equal(extension, []byte("statvfs@openssh.com")):
			path, _ := stream.str()
			ch.Params["path"] = string(path)
			ch.Params["action"] = "statvfs"
			ch.log(1, "hardlink %s", ch.Params["path"])
		case bytes.Equal(extension, []byte("fsync@openssh.com")):
			path, _ := stream.str()
			ch.Params["path"] = string(path)
			ch.Params["action"] = "fsync"
			ch.log(1, "fsync %s", ch.Params["path"])
		default:
			data := stream.Bytes()
			ch.log(1, "Unknown Extension: %s: %s", string(extension), string(data))
		}
	// ignored for brevity
	case sftpExtendedReply:
	}
	// These are only valid for the current message.
	delete(ch.Params, "action")
	delete(ch.Params, "path")
	delete(ch.Params, "target")
}

package detailedlogger

import (
	"log"

	"github.com/tg123/sshpiper/sshpiperd/auditor"
	"golang.org/x/crypto/ssh"
)

// Plugin is basically a generator.  It gets made, and is then asked (with the Create() function) to make
// bespoke auditors for the connection.  Previously we used a naive dumper that didn't need anything bespoke
// so it handled everything. But since we want to actually be able to log specific things (like file accesses
// and ssh commands) we need something more stateful and that thing is the AtomicConnectionAuditor. One
// of which is created per connection and can keep track of which streams are of which session type and
// possibly later on things like which sftp file handle ID's are to which files and so can be logged when
// read from or written to providing greater accuracy in the logs
type Plugin struct {
	Config struct {
		LogFile  string `long:"auditor-detailedlogger-logfile" default:"/dev/stdout" description:"File messages are written to"  env:"SSHPIPERD_AUDITOR_DETAILEDLOGGER_LOGFILE"  ini-name:"auditor-detailedlogger-logfile"`
		LogLevel int    `long:"auditor-detailedlogger-loglevel" default:"1" description:"Log verbosity"  env:"SSHPIPERD_AUDITOR_DETAILEDLOGGER_LOGLEVEL"  ini-name:"auditor-detailedlogger-loglevel"`
	}
}

// GetName returns the name of the plugin (package)
func (p *Plugin) GetName() string {
	return "detailed-logger"
}

// Create actually returns an AtomicConnectionAuditor which can hold more state since the SSH streams
// are stateful creations and to be able to tell a shell from an sftp requires holding onto information
// about that stream for that connection.
func (p *Plugin) Create(conn ssh.ConnMetadata) (auditor.Auditor, error) {
	a := &DetailedConnectionAuditor{}
	a.init()
	a.Filename = p.Config.LogFile
	a.LogLevel = p.Config.LogLevel
	return a, nil
}

// Init satisfies the interface
func (p *Plugin) Init(logger *log.Logger) error {
	return nil
}

// GetUpstreamHook is not actually used. It's just here to satisfy the interface requirement for the plugin struct
func (p *Plugin) GetUpstreamHook() auditor.Hook {
	return func(conn ssh.ConnMetadata, msg []byte) ([]byte, error) {
		return msg, nil
	}
}

// GetDownstreamHook is not actually used. It's just here to satisfy the interface requirement for the plugin struct
func (p *Plugin) GetDownstreamHook() auditor.Hook {
	return func(conn ssh.ConnMetadata, msg []byte) ([]byte, error) {
		return msg, nil
	}
}

// Close satisfies the interface
func (p *Plugin) Close() error {
	return nil
}

// GetOpts satisfies the interface
func (p *Plugin) GetOpts() interface{} {
	return &p.Config
}

// Ran when the package is imported
func init() {
	auditor.Register("detailed-logger", new(Plugin))
}

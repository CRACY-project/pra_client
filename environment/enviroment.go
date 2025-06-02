package environment

// Misc
const (
	logMaxSize     = 10
	_interfaceName = "jimber"
)

// build flags (ldflags)
var (
	environment              string
	version                  string = "0.0.0"
	tag                      string
	buildDate                string
	shaCommitHash            string
	directConnectionsEnabled string
)

func GetEnvironment() string {
	return environment
}

func GetVersion() string {
	return version
}

func GetTag() string {
	return tag
}

func GetBuildDate() string {
	return buildDate
}

func GetShaCommitHash() string {
	return shaCommitHash
}
func GetSignalServerWs() string {
	return signalServerWs
}
func GetSignalServer() string {
	return signalServer
}
func GetLauncherLogPath() string {
	return launcherLogPath
}
func GetDebugLogPath() string {
	return debugLogPath
}
func GetPlatform() string {
	return _platform
}

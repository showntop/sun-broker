package logger

import (
	"os"

	"github.com/Sirupsen/logrus"
)

var Log *logrus.Logger

func init() {
	Log = logrus.New()
	// Log as JSON instead of the default ASCII formatter.
	Log.SetFormatter(&log.JSONFormatter{})

	// Output to stderr instead of stdout, could also be a file.
	Log.SetOutput(os.Stderr)

	// Only log the warning severity or above.
	Log.SetLevel(log.WarnLevel)
}

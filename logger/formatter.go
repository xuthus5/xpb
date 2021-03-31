package logger

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"strings"
)

type LightFormatter struct {
	NoLevel bool
}

func (l *LightFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// output buffer
	b := &bytes.Buffer{}
	// level
	if !l.NoLevel {
		var level = strings.ToUpper(entry.Level.String())
		b.WriteString(level)
		b.WriteString(" ")
	}

	b.WriteString(strings.TrimSpace(entry.Message))
	b.WriteByte('\n')
	return b.Bytes(), nil
}

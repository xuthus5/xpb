package common

import (
	"fmt"
	"os"
)

func CrudeOutput(out string) {
	_, _ = fmt.Fprint(os.Stdout, out)
}

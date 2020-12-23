package monitoringsystem

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime/debug"
	"time"
)

func panicRecover() {
	if r := recover(); r != nil {
		err := fmt.Sprintf("%v", r)
		host, _ := os.Hostname()
		data, _ := json.Marshal(map[string]interface{}{
			"host":        host,
			"app":         "monitoring-lib",
			"file":        "",
			"method":      "panicRecover",
			"type":        "fatal",
			"component":   "Application",
			"code":        "MONITORING.LIBRARY.PANIC",
			"description": err,
			"category":    "UnknownError",
			"doc":         "",
			"ts":          time.Now().String(),
			"ref": map[string]interface{}{
				"stack": string(debug.Stack()),
			},
		})
		fmt.Println(string(data))
	}
}

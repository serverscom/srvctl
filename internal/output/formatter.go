package output

import (
	"encoding/json"
	"fmt"
	"os"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"gopkg.in/yaml.v3"
)

func Format(v interface{}, format string) error {
	switch format {
	case "json":
		data, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			return err
		}
		_, err = os.Stdout.Write(data)
		return err
	case "yaml":
		return yaml.NewEncoder(os.Stdout).Encode(v)
	default:
		return formatText(v)
	}
}

func formatText(v interface{}) error {
	switch data := v.(type) {
	case ContextInfo:
		return formatContextInfo(data)
	case []serverscom.SSHKey:
		return formatSSHKeys(data)
	default:
		return fmt.Errorf("unsupported type for text output: %T", v)
	}
}

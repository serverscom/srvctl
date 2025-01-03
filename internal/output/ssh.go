package output

import (
	"fmt"
	"os"
	"text/tabwriter"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
)

func formatSSHKeys(keys []serverscom.SSHKey) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	fmt.Fprintln(w, "NAME\tFINGERPRINT\tLABELS\tCREATED\tUPDATED")

	var created, updated string
	for _, key := range keys {
		labels := formatLabels(key.Labels)

		created = key.Created.Format("2006-01-02 15:04:05")
		updated = key.Updated.Format("2006-01-02 15:04:05")

		if key.Created.IsZero() {
			created = "<none>"
		}
		if key.Updated.IsZero() {
			updated = "<none>"
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			key.Name,
			key.Fingerprint,
			labels,
			created,
			updated,
		)
	}

	return w.Flush()
}

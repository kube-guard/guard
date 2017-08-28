package cmds

import (
	"fmt"
	"os"

	"github.com/appscode/go-term"
	"github.com/appscode/log"
	"github.com/spf13/cobra"
	"k8s.io/client-go/util/cert"
)

func NewCmdInitCA() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "ca",
		Short:             "Init CA",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			cfg := cert.Config{
				CommonName: "ca",
			}

			store, err := NewCertStore()
			if err != nil {
				log.Fatalf("Failed to create certificate store. Reason: %v.", err)
			}
			if store.IsExists("ca") {
				if !term.Ask(fmt.Sprintf("CA certificate found at %s. Do you want to overwrite?", store.Location()), false) {
					os.Exit(1)
				}
			}

			key, err := cert.NewPrivateKey()
			if err != nil {
				log.Fatalf("Failed to generate private key. Reason: %v.", err)
			}
			cert, err := cert.NewSelfSignedCACert(cfg, key)
			if err != nil {
				log.Fatalf("Failed to generate self-signed certificate. Reason: %v.", err)
			}
			err = store.Write(store.Filename(cfg), cert, key)
			if err != nil {
				log.Fatalf("Failed to init ca. Reason: %v.", err)
			}
			term.Successln("Wrote ca certificates in ", store.Location())
		},
	}
	return cmd
}

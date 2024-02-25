package certificates

import (
	b64 "encoding/base64"
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
	"github.com/angelokurtis/kts-cli/pkg/bash"
)

func download(cmd *cobra.Command, args []string) {
	secrets, err := kubectl.ListTLSSecrets()
	check(err)

	secrets, err = secrets.SelectMany()
	check(err)

	for _, sec := range secrets.Items {
		err = save(sec)
		check(err)
	}
}

func save(secret *kubectl.Secret) error {
	for file, content := range secret.Data {
		filePath := fmt.Sprintf("./certificates/%s/%s", secret.Metadata.Namespace, secret.Metadata.Name)

		_, err := bash.Run("mkdir -p " + filePath)
		if err != nil {
			return err
		}

		d, err := b64.URLEncoding.DecodeString(content)
		if err != nil {
			return errors.WithStack(err)
		}

		if err = ioutil.WriteFile(filePath+"/"+file, d, 0o644); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

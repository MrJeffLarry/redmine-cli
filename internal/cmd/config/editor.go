package config

import (
	"errors"

	"github.com/MrJeffLarry/redmine-cli/internal/config"
	"github.com/MrJeffLarry/redmine-cli/internal/print"
	"github.com/spf13/cobra"
)

func getEditor(r *config.Red_t, cmd *cobra.Command) (string, error) {
	view := cmd.Flags().Arg(0)

	if len(view) <= 0 {
		print.Error("Missing editor exec, please specify what program to use to view")
		return "", errors.New("")
	}
	return view, nil
}

func setGlobalEditor(r *config.Red_t) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set [id]",
		Short: "set editor",
		Long:  "this will set editor globaly in ~/.red/config.json",
		Run: func(cmd *cobra.Command, args []string) {
			var id string
			var err error
			if id, err = getEditor(r, cmd); err != nil {
				return
			}

			r.SetEditor(id)
			if err := r.Save(); err != nil {
				print.Error("%s [%s]", "Could not save editor, please verify permissions", err)
				return
			}
			print.OK("Editor set to [%s]", id)
		},
	}
	return cmd
}

func globalEditor(r *config.Red_t) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "editor",
		Short: "set or get editor",
		Long:  "set or get what editor to use globaly in ~/.red/config.json",
	}
	cmd.AddCommand(setGlobalEditor(r))
	return cmd
}

package config

import (
	"strconv"

	"github.com/MrJeffLarry/redmine-cli/internal/config"
	"github.com/MrJeffLarry/redmine-cli/internal/print"
	"github.com/spf13/cobra"
)

func getCategoryID(r *config.Red_t, cmd *cobra.Command) (int, error) {
	var id int
	var err error

	if id, err = strconv.Atoi(cmd.Flags().Arg(0)); err != nil {
		print.Error("ID is not an valid number, please use `category set 1` for category id 1")
		return 0, err
	}
	return id, nil
}

func setLocalCategory(r *config.Red_t) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set [id]",
		Short: "set local category",
		Long:  "set a local category under .red/config.json",
		Run: func(cmd *cobra.Command, args []string) {
			var id int
			var err error
			if id, err = getCategoryID(r, cmd); err != nil {
				return
			}

			if err := r.SaveLocalCategory(id); err != nil {
				print.Error("%s [%s]", "Could not save category, please verify permissions", err)
				return
			}
			print.OK("Category locally set to ID #%d", id)
		},
	}
	return cmd
}

func setGlobalCategory(r *config.Red_t) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set [id]",
		Short: "set global category",
		Long:  "set a global category under ~/.red/config.json",
		Run: func(cmd *cobra.Command, args []string) {
			var id int
			var err error
			if id, err = getCategoryID(r, cmd); err != nil {
				return
			}

			r.SetCategoryID(id)
			if err := r.Save(); err != nil {
				print.Error("%s [%s]", "Could not save category, please verify permissions", err)
				return
			}
			print.OK("Category globally set to ID #%d", id)
		},
	}
	return cmd
}

func localCategory(r *config.Red_t) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "category",
		Short: "set or get local category",
		Long:  "set or get a local category under ~/.red/config.json",
	}
	cmd.AddCommand(setLocalCategory(r))
	return cmd
}

func globalCategory(r *config.Red_t) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "category",
		Short: "set or get global category",
		Long:  "set or get a global category under ~/.red/config.json",
	}
	cmd.AddCommand(setGlobalCategory(r))
	return cmd
}

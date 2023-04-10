package config

import (
	"github.com/MrJeffLarry/redmine-cli/internal/config"
	"github.com/spf13/cobra"
)

func configGlobal(r *config.Red_t) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "global",
		Short: "get or set global config",
		Long:  "get or set global config",
	}

	cmd.AddCommand(globalProject(r))
	cmd.AddCommand(globalEditor(r))
	cmd.AddCommand(globalPager(r))

	return cmd
}

func configLocal(r *config.Red_t) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "local",
		Short: "get or set local config",
		Long:  "get or set local config",
	}

	cmd.AddCommand(localProject(r))

	return cmd
}

func NewCmdConfig(r *config.Red_t) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "get or set config",
		Long:  "get or set local or global config",
	}

	cmd.AddCommand(globalProject(r))
	cmd.AddCommand(globalEditor(r))
	cmd.AddCommand(globalPager(r))

	cmd.AddCommand(configGlobal(r))
	cmd.AddCommand(configLocal(r))

	return cmd
}

// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package profile

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/larksuite/cli/errs"
	"github.com/larksuite/cli/internal/cmdutil"
	"github.com/larksuite/cli/internal/core"
	"github.com/larksuite/cli/internal/output"
)

// NewCmdProfileEnv creates the profile env subcommand.
func NewCmdProfileEnv(f *cmdutil.Factory) *cobra.Command {
	var environment string
	var lane string
	cmd := &cobra.Command{
		Use:   "env <name>",
		Short: "Bind a profile to prod, pre, or boe runtime endpoints",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return profileEnvRun(f, args[0], environment, lane)
		},
	}
	cmd.Flags().StringVar(&environment, "environment", "", "runtime environment: prod | pre | boe")
	cmd.Flags().StringVar(&lane, "lane", "", "pre/boe lane name, e.g. ppe_hzc_test")
	_ = cmd.MarkFlagRequired("environment")
	cmdutil.SetRisk(cmd, "write")
	return cmd
}

func profileEnvRun(f *cmdutil.Factory, name, environment, lane string) error {
	envName, err := core.NormalizeRuntimeEnvironment(environment)
	if err != nil {
		return errs.NewValidationError(errs.SubtypeInvalidArgument, "%v", err).
			WithCause(err).
			WithParam("--environment")
	}
	lane = strings.TrimSpace(lane)
	if envName == core.RuntimeEnvProd {
		lane = ""
	}

	multi, err := core.LoadOrNotConfigured()
	if err != nil {
		return err
	}
	idx := multi.FindAppIndex(name)
	if idx < 0 {
		return errs.NewValidationError(errs.SubtypeInvalidArgument, "profile %q not found, available profiles: %s", name, strings.Join(multi.ProfileNames(), ", "))
	}
	app := &multi.Apps[idx]
	if _, err := core.ResolveRuntimeEnvironment(app.Brand, envName, lane); err != nil {
		return errs.NewValidationError(errs.SubtypeInvalidArgument, "%v", err).
			WithCause(err).
			WithParam("--lane")
	}

	app.Environment = envName
	app.Lane = lane
	if envName == core.RuntimeEnvProd {
		app.Environment = ""
	}
	if err := core.SaveMultiAppConfig(multi); err != nil {
		return errs.NewInternalError(errs.SubtypeStorage, "failed to save config: %v", err).WithCause(err)
	}

	displayEnv := string(envName)
	if displayEnv == "" {
		displayEnv = string(core.RuntimeEnvProd)
	}
	if lane != "" {
		output.PrintSuccess(f.IOStreams.ErrOut, fmt.Sprintf("Profile %q now uses %s/%s", app.ProfileName(), displayEnv, lane))
		return nil
	}
	output.PrintSuccess(f.IOStreams.ErrOut, fmt.Sprintf("Profile %q now uses %s", app.ProfileName(), displayEnv))
	return nil
}

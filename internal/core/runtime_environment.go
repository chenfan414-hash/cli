// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package core

import (
	"fmt"
	"net/http"
	"strings"
)

// RuntimeEnvironmentName is a constrained environment selector stored on a
// profile. It deliberately avoids arbitrary URL overrides, because app secrets
// and access tokens must only be sent to first-party endpoint groups.
type RuntimeEnvironmentName string

const (
	RuntimeEnvProd RuntimeEnvironmentName = "prod"
	RuntimeEnvPre  RuntimeEnvironmentName = "pre"
	RuntimeEnvBOE  RuntimeEnvironmentName = "boe"
)

// RuntimeEnvironment holds the resolved endpoint group and lane headers for a
// profile.
type RuntimeEnvironment struct {
	Name      RuntimeEnvironmentName
	Brand     LarkBrand
	Lane      string
	Endpoints Endpoints
	Headers   http.Header
}

// NormalizeRuntimeEnvironment canonicalizes stored environment values. Empty
// values are backward compatible and mean prod.
func NormalizeRuntimeEnvironment(value string) (RuntimeEnvironmentName, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", string(RuntimeEnvProd):
		return RuntimeEnvProd, nil
	case string(RuntimeEnvPre):
		return RuntimeEnvPre, nil
	case string(RuntimeEnvBOE):
		return RuntimeEnvBOE, nil
	default:
		return "", fmt.Errorf("environment must be one of prod, pre, boe")
	}
}

// ResolveRuntimeEnvironment resolves endpoints and validated lane headers for
// a profile. Non-production environments require an explicit lane.
func ResolveRuntimeEnvironment(brand LarkBrand, env RuntimeEnvironmentName, lane string) (RuntimeEnvironment, error) {
	name, err := NormalizeRuntimeEnvironment(string(env))
	if err != nil {
		return RuntimeEnvironment{}, err
	}
	lane = strings.TrimSpace(lane)
	resolved := RuntimeEnvironment{
		Name:      name,
		Brand:     ParseBrand(string(brand)),
		Lane:      lane,
		Endpoints: ResolveEndpoints(brand),
	}
	switch name {
	case RuntimeEnvProd:
		if lane != "" {
			return RuntimeEnvironment{}, fmt.Errorf("lane is only supported for pre or boe profiles")
		}
	case RuntimeEnvPre:
		if lane == "" {
			return RuntimeEnvironment{}, fmt.Errorf("lane is required for pre profiles")
		}
		resolved.Endpoints = Endpoints{
			Open:     "https://open.feishu-pre.cn",
			Accounts: "https://accounts.feishu-pre.cn",
			MCP:      "https://mcp.feishu.cn",
			AppLink:  "https://applink.feishu-pre.net",
		}
		resolved.Headers = http.Header{}
		resolved.Headers.Set("x-use-ppe", "1")
	case RuntimeEnvBOE:
		if lane == "" {
			return RuntimeEnvironment{}, fmt.Errorf("lane is required for boe profiles")
		}
		resolved.Endpoints = Endpoints{
			Open:     "https://open.feishu-boe.net",
			Accounts: "https://accounts.feishu-boe.net",
			MCP:      "https://mcp.feishu.cn",
			AppLink:  "https://applink.feishu-boe.net",
		}
		resolved.Headers = http.Header{}
		resolved.Headers.Set("x-use-boe", "1")
	}
	if lane != "" {
		if resolved.Headers == nil {
			resolved.Headers = http.Header{}
		}
		resolved.Headers.Set("X-TT-ENV", lane)
	}
	return resolved, nil
}

// MustResolveRuntimeEnvironment is for internal call sites where config has
// already been validated. It falls back to prod if invoked with invalid data.
func MustResolveRuntimeEnvironment(brand LarkBrand, env RuntimeEnvironmentName, lane string) RuntimeEnvironment {
	resolved, err := ResolveRuntimeEnvironment(brand, env, lane)
	if err != nil {
		return RuntimeEnvironment{Name: RuntimeEnvProd, Brand: ParseBrand(string(brand)), Endpoints: ResolveEndpoints(brand)}
	}
	return resolved
}

// ResolveRuntimeEndpoints returns the endpoint group for a profile.
func ResolveRuntimeEndpoints(brand LarkBrand, env RuntimeEnvironmentName, lane string) Endpoints {
	return MustResolveRuntimeEnvironment(brand, env, lane).Endpoints
}

// RuntimeHeaders returns the extra headers required by a profile's runtime
// environment.
func RuntimeHeaders(env RuntimeEnvironmentName, lane string) http.Header {
	return MustResolveRuntimeEnvironment(BrandFeishu, env, lane).Headers
}

// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package core

import "testing"

func TestResolveRuntimeEnvironment(t *testing.T) {
	tests := []struct {
		name         string
		brand        LarkBrand
		environment  RuntimeEnvironmentName
		lane         string
		wantOpen     string
		wantAccounts string
		wantLane     string
		wantHeader   string
		wantRoute    map[string]string
	}{
		{
			name:         "prod feishu",
			brand:        BrandFeishu,
			environment:  RuntimeEnvProd,
			wantOpen:     "https://open.feishu.cn",
			wantAccounts: "https://accounts.feishu.cn",
		},
		{
			name:         "prod lark",
			brand:        BrandLark,
			environment:  RuntimeEnvProd,
			wantOpen:     "https://open.larksuite.com",
			wantAccounts: "https://accounts.larksuite.com",
		},
		{
			name:         "pre",
			brand:        BrandFeishu,
			environment:  RuntimeEnvPre,
			lane:         " ppe_hzc_test ",
			wantOpen:     "https://open.feishu-pre.cn",
			wantAccounts: "https://accounts.feishu-pre.cn",
			wantLane:     "ppe_hzc_test",
			wantHeader:   "ppe_hzc_test",
			wantRoute:    map[string]string{"x-use-ppe": "1"},
		},
		{
			name:         "boe",
			brand:        BrandFeishu,
			environment:  RuntimeEnvBOE,
			lane:         "boe_lane",
			wantOpen:     "https://open.feishu-boe.net",
			wantAccounts: "https://accounts.feishu-boe.net",
			wantLane:     "boe_lane",
			wantHeader:   "boe_lane",
			wantRoute:    map[string]string{"x-use-boe": "1"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ResolveRuntimeEnvironment(tc.brand, tc.environment, tc.lane)
			if err != nil {
				t.Fatalf("ResolveRuntimeEnvironment() error = %v", err)
			}
			if got.Name != tc.environment {
				t.Fatalf("Name = %q, want %q", got.Name, tc.environment)
			}
			if got.Brand != ParseBrand(string(tc.brand)) {
				t.Fatalf("Brand = %q, want %q", got.Brand, ParseBrand(string(tc.brand)))
			}
			if got.Lane != tc.wantLane {
				t.Fatalf("Lane = %q, want %q", got.Lane, tc.wantLane)
			}
			if got.Endpoints.Open != tc.wantOpen {
				t.Fatalf("Open = %q, want %q", got.Endpoints.Open, tc.wantOpen)
			}
			if got.Endpoints.Accounts != tc.wantAccounts {
				t.Fatalf("Accounts = %q, want %q", got.Endpoints.Accounts, tc.wantAccounts)
			}
			if got.Headers.Get("X-TT-ENV") != tc.wantHeader {
				t.Fatalf("X-TT-ENV = %q, want %q", got.Headers.Get("X-TT-ENV"), tc.wantHeader)
			}
			for name, want := range tc.wantRoute {
				if got.Headers.Get(name) != want {
					t.Fatalf("%s = %q, want %q", name, got.Headers.Get(name), want)
				}
			}
		})
	}
}

func TestResolveRuntimeEnvironmentValidation(t *testing.T) {
	tests := []struct {
		name        string
		environment RuntimeEnvironmentName
		lane        string
	}{
		{name: "prod rejects lane", environment: RuntimeEnvProd, lane: "ppe"},
		{name: "pre requires lane", environment: RuntimeEnvPre},
		{name: "boe requires lane", environment: RuntimeEnvBOE},
		{name: "unknown environment", environment: RuntimeEnvironmentName("dev"), lane: "ppe"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := ResolveRuntimeEnvironment(BrandFeishu, tc.environment, tc.lane); err == nil {
				t.Fatal("ResolveRuntimeEnvironment() error = nil, want validation error")
			}
		})
	}
}

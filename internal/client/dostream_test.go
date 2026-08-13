// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package client_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"

	"github.com/larksuite/cli/errs"
	"github.com/larksuite/cli/internal/cmdutil"
	"github.com/larksuite/cli/internal/core"
	"github.com/larksuite/cli/internal/httpmock"
)

func TestDoStream_HTTPErrorIncludesLogID(t *testing.T) {
	t.Setenv("LARKSUITE_CLI_CONFIG_DIR", t.TempDir())

	config := &core.CliConfig{AppID: "test-app", AppSecret: "test-secret", Brand: core.BrandFeishu}
	factory, _, _, reg := cmdutil.TestFactory(t, config)
	reg.Register(&httpmock.Stub{
		Method:  http.MethodGet,
		URL:     "/open-apis/drive/v1/medias/file_token/download",
		Status:  http.StatusForbidden,
		RawBody: []byte("forbidden"),
		Headers: http.Header{
			larkcore.HttpHeaderKeyLogId: []string{"202605270003"},
		},
	})

	client, err := factory.NewAPIClientWithConfig(config)
	if err != nil {
		t.Fatalf("NewAPIClientWithConfig() error = %v", err)
	}

	_, err = client.DoStream(context.Background(), &larkcore.ApiReq{
		HttpMethod: http.MethodGet,
		ApiPath:    "/open-apis/drive/v1/medias/file_token/download",
	}, core.AsBot)
	var netErr *errs.NetworkError
	if !errors.As(err, &netErr) {
		t.Fatalf("expected *errs.NetworkError, got %T %v", err, err)
	}
	if netErr.LogID != "202605270003" {
		t.Fatalf("LogID = %q, want %q", netErr.LogID, "202605270003")
	}
}

func TestDoStream_RuntimeEnvironmentEndpointAndLaneHeader(t *testing.T) {
	tests := []struct {
		name       string
		env        core.RuntimeEnvironmentName
		lane       string
		wantURL    string
		wantHeader string
	}{
		{
			name:       "pre",
			env:        core.RuntimeEnvPre,
			lane:       "ppe_hzc_test",
			wantURL:    "https://open.feishu-pre.cn/open-apis/drive/v1/medias/file_token/download",
			wantHeader: "x-use-ppe",
		},
		{
			name:       "boe",
			env:        core.RuntimeEnvBOE,
			lane:       "boe_hzc_test",
			wantURL:    "https://open.feishu-boe.net/open-apis/drive/v1/medias/file_token/download",
			wantHeader: "x-use-boe",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("LARKSUITE_CLI_CONFIG_DIR", t.TempDir())

			config := &core.CliConfig{
				AppID:       "test-app",
				AppSecret:   "test-secret",
				Brand:       core.BrandFeishu,
				Environment: tc.env,
				Lane:        tc.lane,
			}
			factory, _, _, reg := cmdutil.TestFactory(t, config)
			stub := &httpmock.Stub{
				Method: http.MethodGet,
				URL:    tc.wantURL,
				Body:   map[string]interface{}{"code": 0, "msg": "ok", "data": map[string]interface{}{}},
			}
			reg.Register(stub)

			client, err := factory.NewAPIClientWithConfig(config)
			if err != nil {
				t.Fatalf("NewAPIClientWithConfig() error = %v", err)
			}
			resp, err := client.DoStream(context.Background(), &larkcore.ApiReq{
				HttpMethod: http.MethodGet,
				ApiPath:    "/open-apis/drive/v1/medias/file_token/download",
			}, core.AsBot)
			if err != nil {
				t.Fatalf("DoStream() error = %v", err)
			}
			resp.Body.Close()

			if got := stub.CapturedHeaders.Get("X-TT-ENV"); got != tc.lane {
				t.Fatalf("X-TT-ENV = %q, want %q", got, tc.lane)
			}
			if got := stub.CapturedHeaders.Get(tc.wantHeader); got != "1" {
				t.Fatalf("%s = %q, want 1", tc.wantHeader, got)
			}
		})
	}
}

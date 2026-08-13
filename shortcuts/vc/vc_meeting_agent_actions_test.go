// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package vc

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"testing"

	"github.com/larksuite/cli/errs"
	"github.com/larksuite/cli/internal/cmdutil"
	"github.com/larksuite/cli/internal/httpmock"
)

func TestMeetingStart_Execute_BodyActionStart(t *testing.T) {
	f, stdout, _, reg := cmdutil.TestFactory(t, defaultConfig())
	stub := &httpmock.Stub{
		Method: "POST",
		URL:    meetingBotStartPath,
		Body: map[string]interface{}{
			"code": 0, "msg": "ok",
			"data": map[string]interface{}{
				"meeting": map[string]interface{}{
					"id":         "7628568141510692381",
					"meeting_no": "123456789",
				},
			},
		},
	}
	reg.Register(stub)

	err := mountAndRun(t, VCMeetingStart, []string{
		"+meeting-start", "--as", "bot", "--meeting-number", "123456789", "--format", "json",
	}, f, stdout)
	if err != nil {
		t.Fatalf("mountAndRun() error = %v", err)
	}

	var req map[string]interface{}
	if err := json.Unmarshal(stub.CapturedBody, &req); err != nil {
		t.Fatalf("unmarshal request body: %v", err)
	}
	if req["join_type"].(float64) != 1 {
		t.Fatalf("join_type = %v, want 1", req["join_type"])
	}
	if req["action"].(float64) != 2 {
		t.Fatalf("action = %v, want 2", req["action"])
	}
	identify, _ := req["join_identify"].(map[string]interface{})
	if identify["meeting_no"] != "123456789" {
		t.Fatalf("join_identify.meeting_no = %v, want 123456789", identify["meeting_no"])
	}
}

func TestMeetingInvite_Execute_SelectedScopeWireMapping(t *testing.T) {
	f, stdout, _, reg := cmdutil.TestFactory(t, defaultConfig())
	stub := &httpmock.Stub{
		Method: "POST",
		URL:    meetingBotInvitePath + "?user_id_type=union_id",
		Body:   map[string]interface{}{"code": 0, "msg": "ok", "data": map[string]interface{}{}},
	}
	reg.Register(stub)

	err := mountAndRun(t, VCMeetingInvite, []string{
		"+meeting-invite", "--as", "bot",
		"--meeting-id", "7628568141510692381",
		"--scope", "selected",
		"--invitee-id-type", "union_id",
		"--invitee-ids", "onion_1,onion_2",
		"--format", "json",
	}, f, stdout)
	if err != nil {
		t.Fatalf("mountAndRun() error = %v", err)
	}

	var req map[string]interface{}
	if err := json.Unmarshal(stub.CapturedBody, &req); err != nil {
		t.Fatalf("unmarshal request body: %v", err)
	}
	if req["meeting_id"] != "7628568141510692381" {
		t.Fatalf("meeting_id = %v, want 7628568141510692381", req["meeting_id"])
	}
	if req["invite_type"].(float64) != 2 {
		t.Fatalf("invite_type = %v, want 2", req["invite_type"])
	}
	if _, ok := req["scope"]; ok {
		t.Fatalf("body must not include scope, got %#v", req["scope"])
	}
	invitees, ok := req["invitees"].([]interface{})
	if !ok || len(invitees) != 2 {
		t.Fatalf("invitees = %#v, want two invitees", req["invitees"])
	}
	first, _ := invitees[0].(map[string]interface{})
	if first["id"] != "onion_1" || first["user_type"].(float64) != 1 {
		t.Fatalf("first invitee = %#v, want id onion_1 user_type 1", first)
	}
}

func TestMeetingInvite_Execute_AllScopeWireMapping(t *testing.T) {
	f, stdout, _, reg := cmdutil.TestFactory(t, defaultConfig())
	stub := &httpmock.Stub{
		Method: "POST",
		URL:    meetingBotInvitePath,
		Body:   map[string]interface{}{"code": 0, "msg": "ok", "data": map[string]interface{}{}},
	}
	reg.Register(stub)

	err := mountAndRun(t, VCMeetingInvite, []string{
		"+meeting-invite", "--as", "bot",
		"--meeting-id", "7628568141510692381",
		"--scope", "all",
		"--format", "json",
	}, f, stdout)
	if err != nil {
		t.Fatalf("mountAndRun() error = %v", err)
	}

	var req map[string]interface{}
	if err := json.Unmarshal(stub.CapturedBody, &req); err != nil {
		t.Fatalf("unmarshal request body: %v", err)
	}
	if req["invite_type"].(float64) != 1 {
		t.Fatalf("invite_type = %v, want 1", req["invite_type"])
	}
	if _, ok := req["invitees"]; ok {
		t.Fatalf("all scope must not include invitees, got %#v", req["invitees"])
	}
	if _, ok := req["scope"]; ok {
		t.Fatalf("body must not include scope, got %#v", req["scope"])
	}
}

func TestMeetingInvite_Validation(t *testing.T) {
	tooMany := make([]string, 201)
	for i := range tooMany {
		tooMany[i] = "ou_user_" + strconv.Itoa(i)
	}
	tests := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "selected missing id type",
			args: []string{"+meeting-invite", "--as", "bot", "--meeting-id", "7628568141510692381", "--scope", "selected", "--invitee-ids", "ou_1"},
			want: "--invitee-id-type",
		},
		{
			name: "selected missing ids",
			args: []string{"+meeting-invite", "--as", "bot", "--meeting-id", "7628568141510692381", "--scope", "selected", "--invitee-id-type", "open_id"},
			want: "--invitee-ids",
		},
		{
			name: "all with ids",
			args: []string{"+meeting-invite", "--as", "bot", "--meeting-id", "7628568141510692381", "--scope", "all", "--invitee-ids", "ou_1"},
			want: "--scope",
		},
		{
			name: "too many selected ids",
			args: []string{"+meeting-invite", "--as", "bot", "--meeting-id", "7628568141510692381", "--scope", "selected", "--invitee-id-type", "open_id", "--invitee-ids", strings.Join(tooMany, ",")},
			want: "--invitee-ids",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f, _, _, _ := cmdutil.TestFactory(t, defaultConfig())
			err := mountAndRun(t, VCMeetingInvite, tc.args, f, nil)
			if err == nil {
				t.Fatal("mountAndRun() error = nil, want validation error")
			}
			var valErr *errs.ValidationError
			if !errors.As(err, &valErr) {
				t.Fatalf("error = %T %v, want *errs.ValidationError", err, err)
			}
			if valErr.Param != tc.want {
				t.Fatalf("Param = %q, want %q", valErr.Param, tc.want)
			}
		})
	}
}

func TestMeetingEnd_Execute_Body(t *testing.T) {
	f, stdout, _, reg := cmdutil.TestFactory(t, defaultConfig())
	stub := &httpmock.Stub{
		Method: "POST",
		URL:    meetingBotEndPath,
		Body:   map[string]interface{}{"code": 0, "msg": "ok", "data": map[string]interface{}{}},
	}
	reg.Register(stub)

	err := mountAndRun(t, VCMeetingEnd, []string{
		"+meeting-end", "--as", "bot", "--meeting-id", "7628568141510692381", "--format", "json",
	}, f, stdout)
	if err != nil {
		t.Fatalf("mountAndRun() error = %v", err)
	}

	var req map[string]interface{}
	if err := json.Unmarshal(stub.CapturedBody, &req); err != nil {
		t.Fatalf("unmarshal request body: %v", err)
	}
	if req["meeting_id"] != "7628568141510692381" {
		t.Fatalf("meeting_id = %v, want 7628568141510692381", req["meeting_id"])
	}
}

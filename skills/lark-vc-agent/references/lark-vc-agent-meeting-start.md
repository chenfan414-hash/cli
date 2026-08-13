# vc +meeting-start

通过 9 位会议号让应用机器人启动并加入一场会议。这是一次**写操作**，会实际让应用机器人入会。

本 skill 对应 shortcut：`lark-cli vc +meeting-start`（调用 `POST /open-apis/vc/v1/bots/join`，请求体包含 `action: 2`）。

## 命令

```bash
lark-cli vc +meeting-start --as bot --meeting-number 123456789
lark-cli vc +meeting-start --as bot --meeting-number 123456789 --format json
lark-cli vc +meeting-start --as bot --meeting-number 123456789 --dry-run
```

## 参数

| 参数 | 必填 | 说明 |
|------|------|------|
| `--meeting-number <no>` | 是 | 会议号，必须为 **9 位纯数字** |
| `--dry-run` | 否 | 预览 API 调用，不实际启动或加入会议 |

## Wire 合同

```json
{
  "join_type": 1,
  "join_identify": {
    "meeting_no": "123456789"
  },
  "action": 2
}
```

仅使用公开 OpenAPI `POST /open-apis/vc/v1/bots/join`。不 fallback BAM、OGW 或 internal RPC。

## 核心约束

- 使用应用身份 `--as bot`。
- `--meeting-number` 只接受 9 位纯数字；不要传会议链接整串，也不要传长数字 `meeting_id`。
- 这是写操作，会让应用机器人实际进入会议并对参会人可见。
- 如果只是查询会中事件，先用 `+meeting-list-active` 获取长数字 `meeting_id`，不要调用本命令。

## 参考

- [lark-vc-agent-meeting-invite](lark-vc-agent-meeting-invite.md) — 会中邀请
- [lark-vc-agent-meeting-end](lark-vc-agent-meeting-end.md) — 结束会议
- [lark-vc-agent-meeting-join](lark-vc-agent-meeting-join.md) — 普通入会

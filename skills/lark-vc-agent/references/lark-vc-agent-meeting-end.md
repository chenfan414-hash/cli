# vc +meeting-end

通过应用机器人结束一场会议。这是一次**写操作**，会影响所有参会人。

本 skill 对应 shortcut：`lark-cli vc +meeting-end`（调用 `POST /open-apis/vc/v1/bots/end`）。

## 命令

```bash
lark-cli vc +meeting-end --as bot --meeting-id <meeting_id>
lark-cli vc +meeting-end --as bot --meeting-id <meeting_id> --format json
lark-cli vc +meeting-end --as bot --meeting-id <meeting_id> --dry-run
```

## 参数

| 参数 | 必填 | 说明 |
|------|------|------|
| `--meeting-id <id>` | 是 | 长数字会议 ID，不是 9 位会议号 |
| `--dry-run` | 否 | 预览 API 调用，不实际结束会议 |

## Wire 合同

```json
{
  "meeting_id": "<meeting_id>"
}
```

仅使用公开 OpenAPI `POST /open-apis/vc/v1/bots/end`。不 fallback BAM、OGW 或 internal RPC。

## 核心约束

- 使用应用身份 `--as bot`。
- `--meeting-id` 必须是长数字会议 ID；9 位会议号会被拒绝。
- 这是影响整场会议的写操作，只在用户明确要求结束会议时调用。

## 参考

- [lark-vc-agent-meeting-start](lark-vc-agent-meeting-start.md) — 启动并加入会议
- [lark-vc-agent-meeting-invite](lark-vc-agent-meeting-invite.md) — 会中邀请
- [lark-vc-agent-meeting-leave](lark-vc-agent-meeting-leave.md) — 仅让机器人离会

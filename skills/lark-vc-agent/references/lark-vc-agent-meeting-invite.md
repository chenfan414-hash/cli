# vc +meeting-invite

通过应用机器人在会中邀请成员。这是一次**写操作**。

本 skill 对应 shortcut：`lark-cli vc +meeting-invite`（调用 `POST /open-apis/vc/v1/bots/invite`）。

## 命令

```bash
# 邀请全部建议成员
lark-cli vc +meeting-invite --as bot --meeting-id <meeting_id> --scope all

# 邀请指定成员
lark-cli vc +meeting-invite --as bot \
  --meeting-id <meeting_id> \
  --scope selected \
  --invitee-id-type open_id \
  --invitee-ids ou_xxx,ou_yyy

# 预览 API 调用
lark-cli vc +meeting-invite --as bot --meeting-id <meeting_id> --scope all --dry-run
```

## 参数

| 参数 | 必填 | 说明 |
|------|------|------|
| `--meeting-id <id>` | 是 | 长数字会议 ID，不是 9 位会议号 |
| `--scope all\|selected` | 是 | `all` 邀请全部建议成员；`selected` 邀请指定成员 |
| `--invitee-id-type open_id\|union_id\|user_id` | selected 必填 | 指定 `--invitee-ids` 的 ID 类型 |
| `--invitee-ids <ids>` | selected 必填 | 指定成员 ID，支持逗号分隔或重复传入，最多 200 个 |

## Wire 合同

`--scope all`：

```json
{
  "meeting_id": "<meeting_id>",
  "invite_type": 1
}
```

不发送 `invitees`，不发送 `scope`。

`--scope selected --invitee-id-type open_id --invitee-ids ou_xxx,ou_yyy`：

```text
POST /open-apis/vc/v1/bots/invite?user_id_type=open_id
```

```json
{
  "meeting_id": "<meeting_id>",
  "invite_type": 2,
  "invitees": [
    {"id": "ou_xxx", "user_type": 1},
    {"id": "ou_yyy", "user_type": 1}
  ]
}
```

不发送 `scope=selected`。仅使用公开 OpenAPI，不 fallback BAM、OGW 或 internal RPC。

## 校验

- `selected` 必须显式传 `--invitee-id-type` 和 `--invitee-ids`。
- `all` 不能传 `--invitee-id-type` 或 `--invitee-ids`。
- `--invitee-ids` 最多 200 个；重复和空值会在本地归一化后去重/忽略。
- `--meeting-id` 是长数字会议 ID，9 位会议号会被拒绝。

## 参考

- [lark-vc-agent-meeting-start](lark-vc-agent-meeting-start.md) — 启动并加入会议
- [lark-vc-agent-meeting-end](lark-vc-agent-meeting-end.md) — 结束会议

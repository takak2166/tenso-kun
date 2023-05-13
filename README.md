# tenso-kun

LINEから他のチャットツール（Webhook）にメッセージを転送してくれる君

## Quick Start
- .envファイルを作成
```.env
CHANNEL_ACCESS_TOKEN="<Your LINE Channel Access Token>"
CHANNEL_SECRET=<Your LINE Channel Secret>
WEBHOOK_LINK=<Webhook Link to Transfer>
```
- dockerコンテナを起動
```bash
$ docker-compose up -d
```
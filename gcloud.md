# gcloud

## notice

- docker container から 別の docker container に接続する場合、`localhost` ではなく、そのコンテナの名前で接続する必要がある。
  - https://zenn.dev/ryo_t/articles/3be7a5ca39d496
- gcloudコマンドを使ってdeployすると`INSTANCE_UNIX_SOCKET`の接続確認ができないので、ブラウザからデプロイしたほうが原因切り分けしやすい。

## sql

### create user

1. create office
2. create hashed_password `go run ./cmd user password "plain_password"`

```sql
insert into users (id, office_id, name, password, role)
VALUES (1, 1, 'admin-reiya', 'hashed_password', 'admin');
```
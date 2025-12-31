# OCI Checker

Oracle Cloud Infrastructure (OCI) の ARM インスタンス (Always Free A1.Flex) を自動で作成・チェックする Go 製ツールです。

## 特徴

- **リソース節約**: ローカルでビルドし、サーバー上ではバイナリを直接実行。
- **KISS 認証**: サーバー上では Instance Principal を使用するため、秘密鍵の転送や管理が不要です。
- **簡単デプロイ**: `just deploy` 一発でサーバーへ転送。

## 必要条件

- [Go](https://go.dev/) (1.25+)
- [just](https://github.com/casey/just)
- [mise](https://mise.jdx.dev/) (推奨)

## セットアップ

1. **OCI コンソールでの準備** (Instance Principal を使う場合):
   - インスタンスを含む **動的グループ** を作成します。
   - そのグループに対して、ターゲットのコンパートメント（またはテナンシー）での `instance-family` および `compute-capacity-reports` の管理権限を与える **ポリシー** を作成します。
   ```bash
   # 動的グループ内のインスタンスが、テナンシー内のどこにあるリソースでも操作できるようにする設定
   Allow dynamic-group YourGroupName to manage instance-family in tenancy
   Allow dynamic-group YourGroupName to manage compute-capacity-reports in tenancy
   Allow dynamic-group YourGroupName to use virtual-network-family in tenancy
   Allow dynamic-group YourGroupName to use volume-family in tenancy
   ```

2. **ローカルでの設定**:
   - `.env` を作成し、リソースの OCID 等を設定します。
   ```bash
   cp .env.example .env
   ```

## 環境変数

認証情報は自動的に取得されるため、`.env` にはリソース情報のみを記述します。

### インスタンス設定（必須）

| 変数名 | 説明 |
|--------|------|
| `OCI_COMPARTMENT_ID` | コンパートメント（またはテナンシー）のOCID |
| `OCI_SUBNET_ID` | サブネットのOCID |
| `OCI_IMAGE_ID` | イメージのOCID |
| `OCI_AVAILABILITY_DOMAIN` | 可用性ドメイン (例: `UlBA:AP-TOKYO-1-AD-1`) |
| `OCI_SSH_PUBLIC_KEY` | インスタンスに登録する SSH 公開鍵 |
| `OCI_DISPLAY_NAME` | 作成するインスタンスの表示名 |

### 動作設定（オプション）

| 変数名 | デフォルト | 説明 |
|--------|----------|------|
| `OCPUS` | 4 | OCPU数 |
| `MEMORY_IN_GBS` | 24 | メモリ (GB) |
| `RETRY_DELAY` | 60 | リトライ間隔（秒） |
| `CHECK_ONLY` | false | キャパシティ確認のみ実行 |

## 使い方

```bash
# ビルド & 転送
just deploy

# ローカル実行 (ローカルの ~/.oci/config を使用)
just run

# サーバー上のログ確認
just logs
```

## ライセンス

MIT

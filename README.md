# OCI A1 Instance Availability Checker (Go Version)

Oracle Cloud Infrastructure (OCI) の Free Tier で人気の高い `VM.Standard.A1.Flex` (Ampere ARM) インスタンスを自動的に取得するための監視ツールです。
Go 版は、シングルバイナリによる軽量な実行と、SDK による堅牢なエラーハンドリングを特徴としています。

## 概要

このツールは OCI Go SDK を使用して、指定した Availability Domain で A1 インスタンスの作成を定期的に試行します。
「容量不足 (Out of host capacity)」エラーが発生した場合は指定した間帰で再試行し、インスタンスが作成されるか、あるいは制限（Limit Exceeded）に達するまで動作し続けます。

## セットアップ

1.  OCI CLI 設定ファイル (`~/.oci/config` および API キー) がホストマシンにあることを確認してください。
2.  `.env` ファイルを作成し、後述の環境変数を設定します。
3.  Docker Compose で起動します。
    ```bash
    docker compose up --build -d
    ```

## 環境変数 (Environment Variables)

`.env` ファイルに設定が必要な変数の一覧です。

### 必須設定 (Required)
| 変数名 | 説明 | 例 |
| :--- | :--- | :--- |
| `COMPARTMENT_ID` | インスタンスを作成するコンパートメントの OCID | `ocid1.compartment.oc1..xxxx` |
| `SUBNET_ID` | インスタンスを配置するサブネットの OCID | `ocid1.subnet.oc1.ap-tokyo-1.xxxx` |
| `IMAGE_ID` | 使用する OS イメージの OCID | `ocid1.image.oc1.ap-tokyo-1.xxxx` |
| `AVAILABILITY_DOMAIN` | 対象の可用性ドメイン名 | `UlBA:AP-TOKYO-1-AD-1` |
| `SSH_PUBLIC_KEY` | インスタンスに登録する SSH 公開鍵の内容 | `ssh-rsa AAAAB3...` |
| `OCPUS` | 割り当てる OCPU 数 (A1.Flex の場合) | `4` |
| `MEMORY_IN_GBS` | 割り当てるメモリサイズ (GB) | `24` |

### 任意設定 (Optional)
| 変数名 | 説明 | デフォルト値 |
| :--- | :--- | :--- |
| `DISPLAY_NAME` | インスタンスの表示名 | `OCI-A1-Instance-Go` |
| `RETRY_DELAY` | 再試行までの待ち時間 (秒) | `60` |
| `CHECK_ONLY` | `true` に設定すると、作成を行わず空き状況の確認 (Peek) のみを行います | `false` |
| `PEEK_BEFORE_LAUNCH` | `true` に設定すると、作成試行の前に空き状況を確認します | `false` |

> [!NOTE]
> A1.Flex インスタンスの空き状況確認（Peek）は、Free Tier においては必ずしも正確ではない場合があります。確実に取得するためには、デフォルト設定（作成を直接試行する）での運用を推奨します。

## 運用とログ

ログを確認して、動作状況を監視できます。
```bash
docker compose logs -f
```

インスタンスが正常に作成されると、OCID がログに出力され、プログラムは自動的に終了します。

# 外部 OpenAPI サービス統合実装方針

## 概要

`internal/repository/announcement.go`で、別の OpenAPI スキーマによって定義されたサービスを呼び出す実装方針をまとめます。

**重要**: 外部サービスは Google Cloud Run で動作しており、サービス間認証が必要です。

## アーキテクチャ

現在のプロジェクトは以下のアーキテクチャパターンを使用しています：

```
Handler → Service → Repository → External API
```

リポジトリ層で外部 API を呼び出すことで、ビジネスロジックと外部依存を分離します。

## 実装方針

### 1. OpenAPI クライアントコード生成（推奨）

#### 1.1 クライアントコード生成の設定

外部サービスの OpenAPI スキーマからクライアントコードを生成します。

**手順：**

1. 外部サービスの OpenAPI スキーマファイルを配置

   ```
   openapi/
     external/
       announcement-service.yaml  # 外部サービスのOpenAPIスキーマ
   ```

2. クライアント生成用の設定ファイルを作成

   ```
   openapi/
     external/
       config.yaml
   ```

   ```yaml
   package: externalclient
   generate:
     client: true
     models: true
   output: generated/external/client.gen.go
   ```

3. コード生成コマンドを実行
   ```bash
   oapi-codegen -config openapi/external/config.yaml openapi/external/announcement-service.yaml
   ```

#### 1.2 リポジトリ実装

生成されたクライアントを使用してリポジトリを実装します。Google Cloud Run のサービス間認証を含む完全な実装例：

```go
package repository

import (
    "context"
    "fmt"
    "net/http"
    "time"
    
    "google.golang.org/api/idtoken"
    "github.com/fun-dotto/app-bff-api/generated/external"
    "github.com/fun-dotto/app-bff-api/internal/config"
    "github.com/fun-dotto/app-bff-api/internal/domain"
)

type ExternalAnnouncementRepository struct {
    client         *external.ClientWithResponses
    config         *config.ExternalServiceConfig
}

func NewExternalAnnouncementRepository(cfg *config.ExternalServiceConfig) (*ExternalAnnouncementRepository, error) {
    if cfg.AnnouncementServiceURL == "" {
        return nil, fmt.Errorf("ANNOUNCEMENT_SERVICE_URL is required")
    }
    if cfg.TargetAudience == "" {
        return nil, fmt.Errorf("ANNOUNCEMENT_SERVICE_TARGET_AUDIENCE is required")
    }
    
    // IDトークンを使用してHTTPリクエストに認証情報を付与するRequestEditorを作成
    requestEditor := func(ctx context.Context, req *http.Request) error {
        tokenSource, err := idtoken.NewTokenSource(ctx, cfg.TargetAudience)
        if err != nil {
            return fmt.Errorf("failed to create token source: %w", err)
        }
        
        token, err := tokenSource.Token()
        if err != nil {
            return fmt.Errorf("failed to get token: %w", err)
        }
        
        req.Header.Set("Authorization", "Bearer "+token.AccessToken)
        return nil
    }
    
    client, err := external.NewClientWithResponses(
        cfg.AnnouncementServiceURL,
        external.WithRequestEditorFn(requestEditor),
    )
    if err != nil {
        return nil, fmt.Errorf("failed to create client: %w", err)
    }
    
    return &ExternalAnnouncementRepository{
        client: client,
        config: cfg,
    }, nil
}

func (r *ExternalAnnouncementRepository) GetAnnouncements() ([]domain.Announcement, error) {
    ctx, cancel := context.WithTimeout(context.Background(), r.config.Timeout)
    defer cancel()
    
    resp, err := r.client.GetAnnouncementsWithResponse(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to call external API: %w", err)
    }

    if resp.StatusCode() != 200 {
        return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
    }

    // 外部APIのレスポンスをドメインモデルに変換
    announcements := make([]domain.Announcement, len(resp.JSON200))
    for i, a := range resp.JSON200 {
        announcements[i] = convertToDomain(a)
    }

    return announcements, nil
}

func convertToDomain(a external.Announcement) domain.Announcement {
    // 外部APIのモデルをドメインモデルに変換
    // 日付のパースなどが必要な場合がある
    date, _ := time.Parse(time.RFC3339, a.Date)
    return domain.Announcement{
        ID:    a.Id,
        Title: a.Title,
        Date:  date,
        URL:   a.Url,
    }
}
```

### 2. 設定管理

外部サービスのエンドポイント URL や認証情報を設定で管理します。

#### 2.1 環境変数による設定

Google Cloud Run のサービス間認証に対応した設定管理：

```go
package config

import (
    "os"
    "time"
)

type ExternalServiceConfig struct {
    AnnouncementServiceURL string
    TargetAudience         string // Google Cloud Run サービスのURL（IDトークンのaudience）
    Timeout                time.Duration
    // サービスアカウントの認証情報は環境変数またはメタデータサーバーから自動取得
}

func LoadExternalServiceConfig() *ExternalServiceConfig {
    return &ExternalServiceConfig{
        AnnouncementServiceURL: getEnv("ANNOUNCEMENT_SERVICE_URL", ""),
        TargetAudience:         getEnv("ANNOUNCEMENT_SERVICE_TARGET_AUDIENCE", ""),
        Timeout:                getDurationEnv("ANNOUNCEMENT_SERVICE_TIMEOUT", 30*time.Second),
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
    if value := os.Getenv(key); value != "" {
        if d, err := time.ParseDuration(value); err == nil {
            return d
        }
    }
    return defaultValue
}
```

**環境変数の設定例：**

```bash
# Google Cloud Run サービスのURL
export ANNOUNCEMENT_SERVICE_URL=https://announcement-service-xxxxx.run.app

# IDトークンのaudience（通常はサービスのURLと同じ）
export ANNOUNCEMENT_SERVICE_TARGET_AUDIENCE=https://announcement-service-xxxxx.run.app

# タイムアウト（オプション）
export ANNOUNCEMENT_SERVICE_TIMEOUT=30s
```

#### 2.2 Google Cloud Run サービス間認証の実装

Google Cloud Run のサービス間認証では、Google Cloud Identity Token（IDトークン）を使用します。

**必要な依存関係：**

```bash
go get google.golang.org/api/idtoken
```

**認証実装：**

```go
package repository

import (
    "context"
    "fmt"
    "net/http"
    
    "google.golang.org/api/idtoken"
    "github.com/fun-dotto/app-bff-api/generated/external"
    "github.com/fun-dotto/app-bff-api/internal/config"
)

type ExternalAnnouncementRepository struct {
    client        *external.ClientWithResponses
    targetAudience string
}

func NewExternalAnnouncementRepository(cfg *config.ExternalServiceConfig) (*ExternalAnnouncementRepository, error) {
    if cfg.AnnouncementServiceURL == "" {
        return nil, fmt.Errorf("ANNOUNCEMENT_SERVICE_URL is required")
    }
    if cfg.TargetAudience == "" {
        return nil, fmt.Errorf("ANNOUNCEMENT_SERVICE_TARGET_AUDIENCE is required")
    }
    
    // IDトークンを使用してHTTPリクエストに認証情報を付与するRequestEditorを作成
    requestEditor := func(ctx context.Context, req *http.Request) error {
        // IDトークンを取得
        tokenSource, err := idtoken.NewTokenSource(ctx, cfg.TargetAudience)
        if err != nil {
            return fmt.Errorf("failed to create token source: %w", err)
        }
        
        token, err := tokenSource.Token()
        if err != nil {
            return fmt.Errorf("failed to get token: %w", err)
        }
        
        // AuthorizationヘッダーにIDトークンを設定
        req.Header.Set("Authorization", "Bearer "+token.AccessToken)
        return nil
    }
    
    client, err := external.NewClientWithResponses(
        cfg.AnnouncementServiceURL,
        external.WithRequestEditorFn(requestEditor),
    )
    if err != nil {
        return nil, fmt.Errorf("failed to create client: %w", err)
    }
    
    return &ExternalAnnouncementRepository{
        client:         client,
        targetAudience: cfg.TargetAudience,
    }, nil
}
```

**認証の仕組み：**

1. **IDトークンの取得**: `idtoken.NewTokenSource()` を使用して、Google Cloud のメタデータサーバーまたはサービスアカウントの認証情報からIDトークンを取得します。
2. **トークンの設定**: 取得したIDトークンを `Authorization: Bearer <token>` ヘッダーに設定します。
3. **自動更新**: Google Cloud SDKはトークンの有効期限を管理し、必要に応じて自動的に更新します。

**ローカル開発環境での認証：**

ローカル開発環境では、Application Default Credentials (ADC) を使用します：

```bash
# gcloud CLIで認証情報を設定
gcloud auth application-default login

# または、サービスアカウントのキーファイルを指定
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/service-account-key.json
```

**Cloud Run での認証：**

Cloud Run で実行する場合、サービスアカウントを指定してデプロイします：

```bash
gcloud run deploy app-bff-api \
  --service-account=app-bff-api@PROJECT_ID.iam.gserviceaccount.com \
  --allow-unauthenticated
```

呼び出し先の Cloud Run サービスで、このサービスアカウントに `run.invoker` ロールを付与します：

```bash
gcloud run services add-iam-policy-binding announcement-service \
  --member="serviceAccount:app-bff-api@PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/run.invoker" \
  --region=asia-northeast1
```

### 3. エラーハンドリング

外部 API 呼び出し時のエラーを適切に処理します。

```go
func (r *ExternalAnnouncementRepository) GetAnnouncements() ([]domain.Announcement, error) {
    ctx, cancel := context.WithTimeout(context.Background(), r.config.Timeout)
    defer cancel()

    resp, err := r.client.GetAnnouncementsWithResponse(ctx)
    if err != nil {
        // ネットワークエラーやタイムアウト
        return nil, fmt.Errorf("failed to call external API: %w", err)
    }

    switch resp.StatusCode() {
    case 200:
        // 成功
    case 401:
        return nil, fmt.Errorf("authentication failed")
    case 404:
        return nil, fmt.Errorf("resource not found")
    case 500:
        return nil, fmt.Errorf("external service error")
    default:
        return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
    }

    // レスポンス処理...
}
```

### 4. リトライとサーキットブレーカー

外部サービスの可用性を向上させるため、リトライやサーキットブレーカーの実装を検討します。

#### 4.1 リトライ実装例

```go
func (r *ExternalAnnouncementRepository) GetAnnouncements() ([]domain.Announcement, error) {
    var lastErr error
    maxRetries := 3

    for i := 0; i < maxRetries; i++ {
        announcements, err := r.getAnnouncementsOnce()
        if err == nil {
            return announcements, nil
        }

        lastErr = err
        // 指数バックオフ
        time.Sleep(time.Duration(i+1) * time.Second)
    }

    return nil, fmt.Errorf("failed after %d retries: %w", maxRetries, lastErr)
}
```

### 5. テスト容易性

インターフェースを定義し、モック実装を容易にします。

#### 5.1 インターフェースの定義

```go
// internal/service/announcement.go に既に定義されている
type AnouncementRepository interface {
    GetAnnouncements() ([]domain.Announcement, error)
}
```

#### 5.2 テスト用のモック実装

```go
// internal/repository/announcement_test.go
type MockExternalAnnouncementRepository struct {
    announcements []domain.Announcement
    err           error
}

func (m *MockExternalAnnouncementRepository) GetAnnouncements() ([]domain.Announcement, error) {
    if m.err != nil {
        return nil, m.err
    }
    return m.announcements, nil
}
```

### 6. 実装手順

1. **依存関係の追加**

   ```bash
   go get google.golang.org/api/idtoken
   ```

2. **外部サービスの OpenAPI スキーマを取得**

   - 外部サービスから OpenAPI スキーマ（YAML/JSON）を取得
   - `openapi/external/`ディレクトリに配置

3. **クライアントコード生成設定を作成**

   - `openapi/external/config.yaml`を作成
   - コード生成コマンドを実行

4. **設定管理を実装**

   - `internal/config/`ディレクトリを作成
   - 環境変数から設定を読み込む機能を実装
   - Google Cloud Run の認証設定（`TargetAudience`）を含める

5. **リポジトリ実装**

   - `internal/repository/external_announcement.go`を作成
   - 生成されたクライアントを使用して実装
   - `idtoken.NewTokenSource()` を使用してIDトークンを取得
   - RequestEditorで各リクエストに認証ヘッダーを設定

6. **依存性注入の更新**

   - `cmd/server/main.go`でリポジトリの初期化を更新
   - 環境変数から設定を読み込み

7. **Google Cloud Run の権限設定**

   - 呼び出し元のサービスアカウントに `run.invoker` ロールを付与
   - Cloud Run サービスのデプロイ時にサービスアカウントを指定

8. **環境変数の設定**

   - `ANNOUNCEMENT_SERVICE_URL`: 外部サービスのURL
   - `ANNOUNCEMENT_SERVICE_TARGET_AUDIENCE`: IDトークンのaudience（通常はサービスのURLと同じ）
   - `ANNOUNCEMENT_SERVICE_TIMEOUT`: タイムアウト（オプション）

9. **テスト実装**
   - ユニットテストを作成（モックを使用）
   - 統合テストを検討（必要に応じて）
   - ローカル開発環境では `gcloud auth application-default login` を実行

## 代替案：手動 HTTP クライアント実装

OpenAPI スキーマが利用できない、または軽量な実装が必要な場合は、標準の`net/http`パッケージを使用して手動実装することも可能です。

```go
type HTTPAnnouncementRepository struct {
    baseURL    string
    httpClient *http.Client
    apiKey     string
}

func (r *HTTPAnnouncementRepository) GetAnnouncements() ([]domain.Announcement, error) {
    req, err := http.NewRequest("GET", r.baseURL+"/announcements", nil)
    if err != nil {
        return nil, err
    }

    req.Header.Set("Authorization", "Bearer "+r.apiKey)

    resp, err := r.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    // レスポンスをパースしてドメインモデルに変換
    // ...
}
```

## 推奨事項

1. **OpenAPI クライアント生成を優先**

   - 型安全性が高い
   - スキーマ変更時の検出が容易
   - メンテナンスコストが低い

2. **設定の外部化**

   - 環境変数や設定ファイルで管理
   - 本番環境と開発環境で切り替え可能に

3. **エラーハンドリングの徹底**

   - 外部サービスのエラーを適切に処理
   - ログ出力を実装

4. **監視とロギング**

   - 外部 API 呼び出しのログを記録
   - メトリクスを収集（呼び出し回数、エラー率など）

5. **タイムアウトの設定**
   - 外部 API 呼び出しにタイムアウトを設定
   - デフォルト値は 30 秒程度を推奨

## Google Cloud Run サービス間認証の詳細

### 認証フロー

1. **IDトークンの取得**: `idtoken.NewTokenSource()` が Google Cloud のメタデータサーバー（Cloud Run 環境）または Application Default Credentials（ローカル環境）から認証情報を取得します。
2. **トークンの生成**: 指定された `targetAudience`（呼び出し先サービスのURL）に対してIDトークンを生成します。
3. **リクエストへの付与**: 各HTTPリクエストの `Authorization` ヘッダーにIDトークンを設定します。
4. **検証**: 呼び出し先の Cloud Run サービスがIDトークンを検証し、適切な権限があることを確認します。

### 必要な権限設定

呼び出し元のサービスアカウントに、呼び出し先の Cloud Run サービスに対する `run.invoker` ロールが必要です。

```bash
# サービスアカウントにロールを付与
gcloud run services add-iam-policy-binding ANNOUNCEMENT_SERVICE_NAME \
  --member="serviceAccount:CALLER_SERVICE_ACCOUNT@PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/run.invoker" \
  --region=REGION
```

### トラブルシューティング

**認証エラーが発生する場合：**

1. **サービスアカウントの確認**
   ```bash
   # Cloud Run サービスのサービスアカウントを確認
   gcloud run services describe SERVICE_NAME --region=REGION --format="value(spec.template.spec.serviceAccountName)"
   ```

2. **IAMポリシーの確認**
   ```bash
   # 呼び出し先サービスのIAMポリシーを確認
   gcloud run services get-iam-policy SERVICE_NAME --region=REGION
   ```

3. **ローカル環境での認証情報の確認**
   ```bash
   # Application Default Credentials が設定されているか確認
   gcloud auth application-default print-access-token
   ```

## 参考リソース

- [oapi-codegen Documentation](https://github.com/deepmap/oapi-codegen)
- [Go HTTP Client Best Practices](https://www.alexedwards.net/blog/how-to-make-http-requests-in-go)
- [Google Cloud Identity Token Authentication](https://cloud.google.com/run/docs/authenticating/service-to-service)
- [google.golang.org/api/idtoken Package](https://pkg.go.dev/google.golang.org/api/idtoken)
- [Cloud Run Service-to-Service Authentication](https://cloud.google.com/run/docs/authenticating/service-to-service)

# 外部 OpenAPI サービス統合実装方針

## 概要

`internal/repository/announcement.go`で、別の OpenAPI スキーマによって定義されたサービスを呼び出す実装方針をまとめます。

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

生成されたクライアントを使用してリポジトリを実装します。

```go
package repository

import (
    "context"
    "github.com/fun-dotto/app-bff-api/generated/external"
    "github.com/fun-dotto/app-bff-api/internal/domain"
)

type ExternalAnnouncementRepository struct {
    client *external.ClientWithResponses
}

func NewExternalAnnouncementRepository(baseURL string) (*ExternalAnnouncementRepository, error) {
    client, err := external.NewClientWithResponses(baseURL)
    if err != nil {
        return nil, err
    }
    return &ExternalAnnouncementRepository{client: client}, nil
}

func (r *ExternalAnnouncementRepository) GetAnnouncements() ([]domain.Announcement, error) {
    ctx := context.Background()
    resp, err := r.client.GetAnnouncementsWithResponse(ctx)
    if err != nil {
        return nil, err
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
    return domain.Announcement{
        ID:    a.Id,
        Title: a.Title,
        Date:  parseDate(a.Date),
        URL:   a.Url,
    }
}
```

### 2. 設定管理

外部サービスのエンドポイント URL や認証情報を設定で管理します。

#### 2.1 環境変数による設定

```go
package config

import (
    "os"
)

type ExternalServiceConfig struct {
    AnnouncementServiceURL string
    APIKey                 string
    Timeout                time.Duration
}

func LoadExternalServiceConfig() *ExternalServiceConfig {
    return &ExternalServiceConfig{
        AnnouncementServiceURL: getEnv("ANNOUNCEMENT_SERVICE_URL", "https://api.example.com"),
        APIKey:                 getEnv("ANNOUNCEMENT_SERVICE_API_KEY", ""),
        Timeout:                getDurationEnv("ANNOUNCEMENT_SERVICE_TIMEOUT", 30*time.Second),
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
```

#### 2.2 認証の実装

API キーや Bearer トークンを使用する場合：

```go
func NewExternalAnnouncementRepository(config *config.ExternalServiceConfig) (*ExternalAnnouncementRepository, error) {
    client, err := external.NewClientWithResponses(
        config.AnnouncementServiceURL,
        external.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
            req.Header.Set("Authorization", "Bearer "+config.APIKey)
            return nil
        }),
    )
    if err != nil {
        return nil, err
    }
    return &ExternalAnnouncementRepository{client: client}, nil
}
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

1. **外部サービスの OpenAPI スキーマを取得**

   - 外部サービスから OpenAPI スキーマ（YAML/JSON）を取得
   - `openapi/external/`ディレクトリに配置

2. **クライアントコード生成設定を作成**

   - `openapi/external/config.yaml`を作成
   - コード生成コマンドを実行

3. **設定管理を実装**

   - `internal/config/`ディレクトリを作成
   - 環境変数から設定を読み込む機能を実装

4. **リポジトリ実装**

   - `internal/repository/external_announcement.go`を作成
   - 生成されたクライアントを使用して実装

5. **依存性注入の更新**

   - `cmd/server/main.go`でリポジトリの初期化を更新
   - 環境変数から設定を読み込み

6. **テスト実装**
   - ユニットテストを作成
   - 統合テストを検討（必要に応じて）

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

## 参考リソース

- [oapi-codegen Documentation](https://github.com/deepmap/oapi-codegen)
- [Go HTTP Client Best Practices](https://www.alexedwards.net/blog/how-to-make-http-requests-in-go)

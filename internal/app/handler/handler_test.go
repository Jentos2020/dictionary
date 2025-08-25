package handler_test

import (
	"context"
	"encoding/json"
	"fmt"
	"leetgo/config"
	"leetgo/internal/app/controller"
	"leetgo/internal/app/handler"
	"leetgo/internal/app/store"
	"leetgo/internal/app/store/pg"
	"leetgo/internal/entity"
	"leetgo/internal/logger"
	"log/slog"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/fasthttp/websocket"
	"github.com/stretchr/testify/require"
)

func TestWSHandlerIntegration(t *testing.T) {
	ctx := context.Background()
	log := slog.New(logger.CustomFormat(os.Stdout))

	os.Setenv("CONFIG_PATH", "/home/evgeniy/go/leetgo/config/config.yaml")

	cfg, err := config.New()
	require.NoError(t, err)

	cfg.App.Port = cfg.App.TestPort
	cfg.Dicts = "/home/evgeniy/go/leetgo/data/test/"
	cfg.Schema = "public"
	cfg.MigrationsPath = "/home/evgeniy/go/leetgo/migrations"

	dsn := cfg.GetDsn()
	db, err := pg.New(
		dsn,
		cfg.Schema,
		cfg.MaxConn,
		cfg.MaxIdle,
		log,
	)
	require.NoError(t, err)

	err = store.MigrateUp(ctx, cfg, log)
	require.NoError(t, err)

	res := db.Db.Exec("DROP TABLE IF EXISTS public.russian_test")
	require.NoError(t, res.Error)

	err = db.WriteDictsToDb(ctx, cfg, "russian_test")
	require.NoError(t, res.Error)

	c := controller.New(ctx, cfg, db, log)

	err = c.FillTrieWithWords(ctx, "russian_test")
	require.NoError(t, err)

	app := handler.New(c)

	go func() {
		if err := app.Listen(":8888"); err != nil {
			t.Logf("Server stopped: %v", err)
		}
	}()
	time.Sleep(200 * time.Millisecond)

	u := url.URL{Scheme: "ws", Host: "localhost:8888", Path: "/ws/search"}
	wsConn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	require.NoError(t, err)
	defer wsConn.Close()

	req := entity.SearchRequest{Prefix: "абстракт"}
	reqJSON, _ := json.Marshal(req)
	err = wsConn.WriteMessage(websocket.TextMessage, reqJSON)
	require.NoError(t, err)

	_, msg, err := wsConn.ReadMessage()
	if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
		t.Logf("WebSocket connection closed normally")
	} else {
		require.NoError(t, err)
	}

	var resp entity.SearchResponse
	err = json.Unmarshal(msg, &resp)
	require.NoError(t, err)

	for _, w := range resp.Words {
		fmt.Printf("Found \"%s\"\n", w.Data)
		require.True(t, strings.HasPrefix(w.Data, req.Prefix),
			fmt.Sprintf("Word %s does not start with prefix %s", w.Data, req.Prefix))
	}
}

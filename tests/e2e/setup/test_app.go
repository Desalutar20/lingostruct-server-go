package setup

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"path"
	"runtime"
	"testing"
	"time"

	"github.com/Desalutar20/lingostruct-server-go/internal/app"
	"github.com/Desalutar20/lingostruct-server-go/internal/config"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type TestApp struct {
	pool       *pgxpool.Pool
	redis      *redis.Client
	httpClient *http.Client
	addr       string
	config     *config.ApplicationConfig
}

func init() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("cannot get caller info")
	}

	if err := godotenv.Load(path.Join(filename, "../../../../.env.test")); err != nil {
		panic(err)
	}
}

func Run(t *testing.T, fn func(testApp *TestApp)) {
	ctx := context.Background()

	config := config.New()
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	log.SetOutput(io.Discard)

	config.Database.Name = fmt.Sprintf("test-%s", uuid.NewString())

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}

	redis, cleanupRedis := setupRedis(ctx, &config.Redis)
	postgres, cleanupPostgres := setupDatabase(ctx, &config.Database)

	app := app.New(ctx, config, logger, &listener)

	t.Cleanup(func() {
		app.Close(ctx)
		listener.Close()
		cleanupRedis()
		cleanupPostgres()
	})

	go app.Run()

	fn(&TestApp{
		pool:  postgres,
		redis: redis,
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
		addr:   fmt.Sprintf("http://127.0.0.1:%d", listener.Addr().(*net.TCPAddr).Port),
		config: &config.Application,
	})
}

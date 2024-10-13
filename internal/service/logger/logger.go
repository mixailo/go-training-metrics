package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Log будет доступен всему коду как синглтон.
// Никакой код навыка, кроме функции Initialize, не должен модифицировать эту переменную.
// По умолчанию установлен no-op-логер, который не выводит никаких сообщений.
var Log *zap.Logger = zap.NewNop()

// Initialize инициализирует синглтон логера с необходимым уровнем логирования.
func Initialize(level string) error {
	// преобразуем текстовый уровень логирования в zap.AtomicLevel
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}
	// создаём новую конфигурацию логера
	cfg := zap.NewProductionConfig()
	// устанавливаем уровень
	cfg.Level = lvl
	// создаём логер на основе конфигурации
	zl, err := cfg.Build()
	if err != nil {
		return err
	}
	// устанавливаем синглтон
	Log = zl
	return nil
}

type ResponseData struct {
	httpStatus int
	size       int
	elapsed    time.Duration
}

type LoggingResponseWriter struct {
	http.ResponseWriter
	ResponseData ResponseData
}

func (rw *LoggingResponseWriter) Write(b []byte) (int, error) {
	// записываем ответ, используя оригинальный http.ResponseWriter
	size, err := rw.ResponseWriter.Write(b)
	rw.ResponseData.size += size // захватываем размер
	return size, err
}

func (rw *LoggingResponseWriter) WriteHeader(statusCode int) {
	// записываем код статуса, используя оригинальный http.ResponseWriter
	rw.ResponseWriter.WriteHeader(statusCode)
	rw.ResponseData.httpStatus = statusCode // захватываем код статуса
}

// RequestResponseLogger — middleware-логер для входящих HTTP-запросов.
func RequestResponseLogger(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		Log.Info("got incoming HTTP request",
			zap.String("method", r.Method),
			zap.String("URI", r.RequestURI),
		)

		rw := &LoggingResponseWriter{
			ResponseWriter: w,
			ResponseData: ResponseData{
				httpStatus: 0,
				size:       0,
			},
		}

		t1 := time.Now()

		defer func() {
			Log.Info("completed HTTP request",
				zap.Int("code", rw.ResponseData.httpStatus),
				zap.Duration("duration", time.Since(t1)),
			)
		}()

		next.ServeHTTP(rw, r)
	}
	return http.HandlerFunc(fn)
}

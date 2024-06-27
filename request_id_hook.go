package logging

import (
	"github.com/binamico/requestid"
	"github.com/sirupsen/logrus"
)

// RequestIDHook представляет собой хук для logrus, который добавляет идентификатор запроса к записям журнала.
// Идентификатор запроса извлекается из контекста и добавляется в данные записи.
type RequestIDHook struct {
	levels []logrus.Level // Уровни логирования, на которых хук активен
}

// NewRequestIDHook создает новый экземпляр RequestIDHook.
// Принимает уровни логирования, на которых хук должен срабатывать, и возвращает хук.
func NewRequestIDHook(levels ...logrus.Level) *RequestIDHook {
	return &RequestIDHook{
		levels: levels,
	}
}

// Заставляем компилятор проверять, что RequestIDHook реализует интерфейс logrus.Hook.
var _ logrus.Hook = (*RequestIDHook)(nil)

// Fire вызывается logrus при создании записи журнала.
// Добавляет идентификатор запроса к записи, если он присутствует в контексте.
func (hook *RequestIDHook) Fire(entry *logrus.Entry) error {
	if entry == nil {
		return nil // Возвращаем nil, если запись пустая
	}
	ctx := entry.Context
	if ctx == nil {
		return nil // Возвращаем nil, если контекст отсутствует
	}

	// Извлекаем идентификатор запроса из контекста
	requestID, ok := requestid.EjectRequestID(ctx)
	if !ok {
		return nil // Возвращаем nil, если идентификатор запроса не найден
	}
	entry.Data["requestId"] = requestID // Добавляем идентификатор запроса в данные записи
	return nil
}

// Levels возвращает уровни логирования, на которых хук должен срабатывать.
// Эти уровни были заданы при создании хука.
func (hook *RequestIDHook) Levels() []logrus.Level {
	return hook.levels
}

package log

type Config struct {
	// Level уровень логирования
	Level string `env:"LOGGER_LEVEL,default=info"`
	// Timestamp отображать время вызова события
	DisableTimestamp bool `env:"LOGGER_DISABLE_TIMESTAMP"`
	// Caller отображать откуда был вызван логгер
	Caller bool `env:"LOGGER_ENABLE_CALLER"`
	// Console использовать человеко-читаемый вывод в консоли(с подсветкой, форматированием), в противном случае - json
	Console bool `env:"LOGGER_ENABLE_CONSOLE"`
}

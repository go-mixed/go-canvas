package misc

// Logger 定义 RichText 可注入的日志接口。
// Logger defines the injectable logging interface used by RichText.
//
// 标准库 *log.Logger 直接满足该接口（含 Printf）。
// Standard library *log.Logger satisfies this interface directly.
type Logger interface {
	Printf(format string, v ...any)
}

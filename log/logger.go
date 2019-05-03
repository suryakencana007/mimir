/*  logger.go
*
* @Author:             Nanang Suryadi
* @Date:               February 12, 2019
* @Last Modified by:   @suryakencana007
* @Last Modified time: 2019-02-12 11:38
 */

package log

type logging interface {
	Debug(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Fatal(msg string, fields ...interface{})
	Panic(msg string, fields ...interface{})

	// for type fields
	Field(key string, value interface{}) interface{}
}

var logger logging

func Debug(msg string, fields ...interface{}) {
	logger.Debug(msg, fields...)
}

// Info add log entry with info level
func Info(msg string, fields ...interface{}) {
	logger.Info(msg, fields...)
}

// Warn add log entry with warn level
func Warn(msg string, fields ...interface{}) {
	logger.Warn(msg, fields...)
}

// Error add log entry with error level
func Error(msg string, fields ...interface{}) {
	logger.Error(msg, fields...)
}

// Fatal add log entry with fatal level
func Fatal(msg string, fields ...interface{}) {
	logger.Fatal(msg, fields...)
}

// Panic add log entry with panic level
func Panic(msg string, fields ...interface{}) {
	logger.Panic(msg, fields...)
}

func Field(key string, value interface{}) interface{} {
	return logger.Field(key, value)
}

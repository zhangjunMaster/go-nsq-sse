/**
 * logger instance method
 * By calling an instance method,
 * you do not need to repeatedly generate an instance
 */

package logger

import "log"

var loggerInstance = NewRealStLogger(0)

func LOG_DEBUG(format string, args ...interface{}) {
	log.Println(format)
	loggerInstance.DEBUG(format, args...)
}

func LOG_INFO(format string, args ...interface{}) {
	loggerInstance.INFO(format, args...)
}

func LOG_WARNING(format string, args ...interface{}) {
	loggerInstance.WARNING(format, args...)
}

func LOG_ERROR(format string, args ...interface{}) {
	loggerInstance.ERROR(format, args...)
}

func LOG_CRITIC(format string, args ...interface{}) {
	loggerInstance.CRITIC(format, args...)
}

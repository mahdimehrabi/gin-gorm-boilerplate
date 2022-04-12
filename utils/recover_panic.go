package utils

import "boilerplate/core/infrastructure"

// RecoverPanic recovers panic in the application
func RecoverPanic(logger infrastructure.Logger) func() {
	return func() {
		if r := recover(); r != nil {
			logger.Zap.Info("☠️ panic recovered: ", r)
		}
	}
}

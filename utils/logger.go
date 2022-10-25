package utils

import "go.uber.org/zap"

var Logger, _ = zap.NewDevelopment(nil)

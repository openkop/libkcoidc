/*
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Copyright 2020 Kopano and its licensors
 */

package main

/*
typedef void (*kcoidc_cb_func_log_s) (char*);

void bridge_kcoidc_log_cb_func_log_s(kcoidc_cb_func_log_s f, char* s);
*/
import "C"
import (
	"fmt"
	"log"
	"os"
	"strings"

	"stash.kopano.io/kc/libkcoidc"
)

type callbackLogger struct {
	cb C.kcoidc_cb_func_log_s
}

func (logger *callbackLogger) Printf(format string, a ...interface{}) {
	s := strings.TrimRight(fmt.Sprintf(format, a...), "\n")
	C.bridge_kcoidc_log_cb_func_log_s(logger.cb, C.CString(s))
}

func getSimpleLogger(prefix string) kcoidc.Logger {
	return log.New(os.Stdout, prefix, 0)
}

func getCLogger(cb C.kcoidc_cb_func_log_s) kcoidc.Logger {
	return &callbackLogger{
		cb: cb,
	}
}

var defaultLogger kcoidc.Logger

func getDefaultDebugLogger() kcoidc.Logger {
	if defaultLogger == nil {
		defaultLogger = getSimpleLogger("[kcoidc-c debug] ")
	}

	return defaultLogger
}

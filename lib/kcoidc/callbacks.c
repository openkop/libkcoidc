/*
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Copyright 2020 Kopano and its licensors
 */

#include "kcoidc_callbacks.h"

void bridge_kcoidc_log_cb_func_log_s(kcoidc_cb_func_log_s f, char* s)
{
	return f(s);
}

void bridge_kcoidc_watch_cb_func_updated(kcoidc_cb_func_watch f)
{
	return f();
}

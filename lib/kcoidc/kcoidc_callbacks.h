/*
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Copyright 2020 Kopano and its licensors
 */

#ifndef KCOIDC_CALLBACKS_H
#define KCOIDC_CALLBACKS_H

typedef void (*kcoidc_cb_func_log_s) (char*);
typedef void (*kcoidc_cb_func_watch) ();

void bridge_kcoidc_log_cb_func_log_s(kcoidc_cb_func_log_s f, char* s);
void bridge_kcoidc_watch_cb_func_updated(kcoidc_cb_func_watch f);

#endif /* !KCOIDC_CALLBACKS_H */

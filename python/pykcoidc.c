/*
 * Copyright 2018 Kopano and its licensors
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License, version 3
 * or later, as published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

#include <Python.h>
#include <kcoidc.h>

#if PY_MAJOR_VERSION >= 3
#define PY3K
#endif

#if KCOIDC_VERSION >= 10100
#define WITH_REQUIRE_SCOPE
#endif

static PyObject *PyKCOIDCError;

static PyObject *
pykcoidc_initialize(PyObject *self, PyObject *args)
{
	const char *iss_s;
	int res;

	if (!PyArg_ParseTuple(args, "s", &iss_s))
		return NULL;

	Py_BEGIN_ALLOW_THREADS;
	res = kcoidc_initialize(iss_s);
	Py_END_ALLOW_THREADS;

	if (res != 0) {
		PyErr_SetObject(PyKCOIDCError, PyLong_FromLong(res));
		return NULL;
	}

	return PyLong_FromLong(res);
}

static PyObject *
pykcoidc_wait_until_ready(PyObject *self, PyObject *args)
{
	unsigned long long timeout;
	int res;

	if (!PyArg_ParseTuple(args, "K", &timeout))
		return NULL;

	Py_BEGIN_ALLOW_THREADS;
	res = kcoidc_wait_until_ready(timeout);
	Py_END_ALLOW_THREADS;

	if (res != 0) {
		PyErr_SetObject(PyKCOIDCError, PyLong_FromLong(res));
		return NULL;
	}

	return PyLong_FromLong(res);
}

static PyObject *
pykcoidc_insecure_skip_verify(PyObject *self, PyObject *args)
{
	int insecure;
	int res;

	if (!PyArg_ParseTuple(args, "i", &insecure))
		return NULL;

	Py_BEGIN_ALLOW_THREADS;
	res = kcoidc_insecure_skip_verify(insecure);
	Py_END_ALLOW_THREADS;

	if (res != 0) {
		PyErr_SetObject(PyKCOIDCError, PyLong_FromLong(res));
		return NULL;
	}

	return PyLong_FromLong(res);
}

static PyObject *
pykcoidc_validate_token_s(PyObject *self, PyObject *args)
{
	PyObject *res = NULL;
	const char *token_s;
	struct kcoidc_validate_token_s_return token_result;

	if (!PyArg_ParseTuple(args, "s", &token_s))
		return NULL;

	Py_BEGIN_ALLOW_THREADS;
	token_result = kcoidc_validate_token_s(token_s);
	Py_END_ALLOW_THREADS;

	if (token_result.r1 != 0) {
		PyErr_SetObject(PyKCOIDCError, PyLong_FromLong(token_result.r1));
	} else {
		res = Py_BuildValue("zizz", token_result.r0, token_result.r2, token_result.r3, token_result.r4);
	}

	// Free the strings passed from the library.
	free(token_result.r0);
	free(token_result.r3);
	free(token_result.r4);

	return res;
}

#ifdef WITH_REQUIRE_SCOPE
static PyObject *
pykcoidc_validate_token_and_require_scope_s(PyObject *self, PyObject *args)
{
	PyObject *res = NULL;
	const char *token_s;
	const char *required_scope_s;
	struct kcoidc_validate_token_and_require_scope_s_return token_result;

	if (!PyArg_ParseTuple(args, "ss", &token_s, &required_scope_s))
		return NULL;

	Py_BEGIN_ALLOW_THREADS;
	token_result = kcoidc_validate_token_and_require_scope_s(token_s, required_scope_s);
	Py_END_ALLOW_THREADS;

	if (token_result.r1 != 0) {
		PyErr_SetObject(PyKCOIDCError, PyLong_FromLong(token_result.r1));
	} else {
		res = Py_BuildValue("zizz", token_result.r0, token_result.r2, token_result.r3, token_result.r4);
	}

	// Free the strings passed from the library.
	free(token_result.r0);
	free(token_result.r3);
	free(token_result.r4);

	return res;
}
#endif

static PyObject *
pykcoidc_fetch_userinfo_with_accesstoken_s(PyObject *self, PyObject *args)
{
	PyObject *res = NULL;
	const char *token_s;
	struct kcoidc_fetch_userinfo_with_accesstoken_s_return userinfo_result;

	if (!PyArg_ParseTuple(args, "s", &token_s))
		return NULL;

	Py_BEGIN_ALLOW_THREADS;
	userinfo_result = kcoidc_fetch_userinfo_with_accesstoken_s(token_s);
	Py_END_ALLOW_THREADS;

	if (userinfo_result.r1 != 0) {
		PyErr_SetObject(PyKCOIDCError, PyLong_FromLong(userinfo_result.r1));
	} else {
		res = Py_BuildValue("z", userinfo_result.r0);
	}

	// Free the strings passed from the library.
	free(userinfo_result.r0);

	return res;
}

static PyObject *
pykcoidc_uninitialize(PyObject *self, PyObject *args)
{
	int res;

	if (!PyArg_ParseTuple(args, ""))
		return NULL;

	Py_BEGIN_ALLOW_THREADS;
	res = kcoidc_uninitialize();
	Py_END_ALLOW_THREADS;

	if (res != 0) {
		PyErr_SetObject(PyKCOIDCError, PyLong_FromLong(res));
		return NULL;
	}

	return PyLong_FromLong(res);
}

static PyMethodDef MyMethods[] = {
	{"initialize", pykcoidc_initialize, METH_VARARGS, "Initialize ODIC."},
	{"wait_until_ready", pykcoidc_wait_until_ready, METH_VARARGS, "Wait until ODIC is ready or until timeout."},
	{"insecure_skip_verify", pykcoidc_insecure_skip_verify, METH_VARARGS, "Set insecure skip verify flag."},
	{"validate_token_s", pykcoidc_validate_token_s, METH_VARARGS, "Validate token and return authenticted user ID."},
#ifdef WITH_REQUIRE_SCOPE
	{"validate_token_and_require_scope_s", pykcoidc_validate_token_and_require_scope_s, METH_VARARGS, "Validate token and scope and return authenticated user ID."},
#endif
	{"fetch_userinfo_with_accesstoken_s", pykcoidc_fetch_userinfo_with_accesstoken_s, METH_VARARGS, "Fetch userinfo with access token."},
	{"uninitialize",  pykcoidc_uninitialize, METH_VARARGS, "Uninitialize ODIC."},
	{NULL, NULL, 0, NULL} /* Sentinel */
};

#ifdef PY3K
static struct PyModuleDef myModule = {
	PyModuleDef_HEAD_INIT,
	"_pykcoidc",
	NULL,
	-1,
	MyMethods
};
PyMODINIT_FUNC
PyInit__pykcoidc(void)
{
	PyObject *m;

	m = PyModule_Create(&myModule);
	if (m == NULL)
		return NULL;

	PyKCOIDCError = PyErr_NewException("_pykcoidc.Error", NULL, NULL);
	Py_INCREF(PyKCOIDCError);
	PyModule_AddObject(m, "Error", PyKCOIDCError);

	return m;
}
#else // PY3K
void init_pykcoidc(void)
{
	PyObject *m;

	m = Py_InitModule3("_pykcoidc", MyMethods, NULL);
	if (m == NULL)
		return;

	PyKCOIDCError = PyErr_NewException("_pykcoidc.Error", NULL, NULL);
	Py_INCREF(PyKCOIDCError);
	PyModule_AddObject(m, "Error", PyKCOIDCError);

	return;
}
#endif // PY3K

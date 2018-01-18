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

#include <stdio.h>
#include <stdlib.h>
#include <time.h>

#include "kcoidc.h"

int main(int argc, char** argv)
{
	clock_t begin = clock();
	char* token_s = argv[1];
	struct kcoidc_validate_token_s_return res = kcoidc_validate_token_s(token_s);
	clock_t end = clock();
	double time_spent = (double)(end - begin) / CLOCKS_PER_SEC;

	printf("\n");
	printf("Token subject : %s -> %s\n", res.r0, res.r1 ? "valid" : "invalid");
	printf("Time spent    : %8fs\n", time_spent);

	free(res.r0);

	if (!res.r1) {
		return -1;
	}
	return 0;
}

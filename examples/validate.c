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
	int res;
	int res2;
	clock_t begin = clock();
	clock_t end;
	double time_spent;

	char* iss_s = argv[1];
	char* token_s = argv[2];
	struct kcoidc_validate_token_s_return valid;

	// Initialize first.
	res = kcoidc_initialize(iss_s);
	if (res != 0) {
		printf("> Error: initialize failed: 0x%x\n", res);
		goto exit;
	}
	// Wait until oidc validation becomes ready.
	res = kcoidc_wait_untill_ready(10);
	if (res != 0) {
		printf("> Error: failed to get ready in time: 0x%x\n", res);
		goto exit;
	}

	// validate token passed from commandline.
	valid = kcoidc_validate_token_s(token_s);
	end = clock();
	time_spent = (double)(end - begin) / CLOCKS_PER_SEC;

	// Show the result.
	printf("> Token subject : %s -> %s\n", valid.r0, valid.r1 == 0 ? "valid" : "invalid");
	printf("> Time spent    : %8fs\n", time_spent);

	// Free the subjects memory.
	free(valid.r0);

	// Handle validation result.
	res = valid.r1;
	printf("> Result code   : 0x%x\n", res);

	// Remember to uninitialize on success as well.
	res2 = kcoidc_uninitialize();
	if (res2 != 0) {
		printf("> Error: failed to uninitialize: 0x%x\n", res2);
	}

exit:
	if (res != 0) {
		return -1;
	}
	return 0;
}

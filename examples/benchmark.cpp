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

#include <chrono>
#include <iostream>
#include <string>
#include <sstream>
#include <thread>
#include <mutex>
#include <vector>

namespace kcoidc {
	extern "C" {
		#include "kcoidc.h"
	}
}

std::mutex log_mutex;
void log(std::string msg)
{
	std::lock_guard<std::mutex> lock(log_mutex);
	std::cout << "> " << msg;
}

void bench_validateTokenS(int id, int count, std::string token_s)
{
	std::stringstream msg;
	msg << "Info : thread " << id << " started ..." << std::endl;
	log(msg.str());

	unsigned success = 0;
	unsigned failed = 0;
	for (auto c = 0; c < count; ++c) {
		auto result = kcoidc::kcoidc_validate_token_s(&token_s[0u]);
		if (result.r1 != 0) {
			// Error.
			failed++;
			msg.str("");
			msg << "Error: validation failed: " << "0x" << std::hex << result.r1  << std::endl;
			log(msg.str());
		} else {
			success++;
		};
	}

	msg.str("");
	msg << "Info : thread " << id << " done:" << success << " failed:" << failed << std::endl;
	log(msg.str());
}

int handleError(std::string msg, int code)
{
	std::cerr << "> Error: " << msg << std::endl;
	return -1;
}

int main(int argc, char** argv)
{
	std::vector<std::thread> threads;
	std::string iss_s = argv[1];
	std::string token_s = argv[2];
	int res;

	if ((res = kcoidc::kcoidc_insecure_skip_verify(1)) != 0) {
		return handleError("insecure_skip_verify failed", res);
	};
	if ((res = kcoidc::kcoidc_initialize(&iss_s[0u])) != 0) {
		return handleError("initialize failed", res);
	};
	if ((res = kcoidc::kcoidc_wait_untill_ready(10)) != 0) {
		return handleError("failed to get ready in time", res);
	};

	int concurentThreadsSupported = std::thread::hardware_concurrency();
	int count = 100000;
	std::cout << "> Info : using " << concurentThreadsSupported << " threads with " << count << " runs per thread" << std::endl;
	auto  begin_time = std::chrono::system_clock::now();
	for (auto i = 1; i <= concurentThreadsSupported; ++i) {
		threads.push_back(std::thread(bench_validateTokenS, i, count, token_s));
	}
	for (auto& th : threads) {
		th.join();
	}
	auto end_time = std::chrono::system_clock::now();
	auto duration = std::chrono::duration_cast<std::chrono::milliseconds>(
		end_time - begin_time
	);
	auto rate = (count * concurentThreadsSupported) / (float(duration.count())/1000);
	std::cout << "> Time : " << float(duration.count())/1000 << "s" << std::endl;
	std::cout << "> Rate : " << rate << " op/s" << std::endl;

	if ((res = kcoidc::kcoidc_uninitialize()) != 0) {
		return handleError("failed to uninitialize", res);
	};

	return 0;
}

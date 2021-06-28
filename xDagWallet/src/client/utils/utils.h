//
//  utils.h
//  xdag
//
//  Copyright © 2018 xdag contributors.
//

#ifndef XDAG_UTILS_HEADER_H
#define XDAG_UTILS_HEADER_H

#include <stdio.h>
#include <stdint.h>
#include <pthread.h>
#include "../types.h"

#ifdef _WIN32
#define DELIMITER "\\"
#else
#define DELIMITER "/"
#endif

#ifdef __cplusplus
extern "C" {
#endif

	extern uint64_t get_timestamp(void);

	extern void xdag_init_path(char *base);
	extern FILE* xdag_open_file(const char *path, const char *mode);
	extern void xdag_close_file(FILE *f);
	extern int xdag_file_exists(const char *path);
	extern int xdag_mkdir(const char *path);

	long double log_difficulty2hashrate(long double log_diff);
	void xdag_str_toupper(char *str);
	void xdag_str_tolower(char *str);
	char *xdag_basename(char *path);
	char *xdag_filename(char *_filename);

	// convert xdag_time_t to string representation
	// minimal length of string buffer `buf` should be 60
	void xdag_time_to_string(xdag_time_t time, char *buf);

	// convert time_t to string representation
	// minimal length of string buffer `buf` should be 50
	void time_to_string(time_t time, char* buf);
	size_t validate_remark(const char *str);
	size_t validate_ascii_safe(const char *str, size_t maxsize);

#ifdef __cplusplus
};
#endif

#endif /* utils_h */

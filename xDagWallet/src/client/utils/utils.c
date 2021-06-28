//
//  utils.c
//  xdag
//
//  Copyright © 2018 xdag contributors.
//

#include "utils.h"
#include <stdlib.h>
#include <string.h>
#include <ctype.h>
#include <sys/stat.h>
#include <time.h>

#if defined (__MACOS__) || defined (__APPLE__)
#include <libgen.h>
#define PATH_MAX 4096
#include <sys/time.h>
#include <unistd.h>
#elif defined(_WIN32) || defined(_WIN64)
#include <direct.h>
#include <shlwapi.h>
//// #include "../../win/unistd.h"
#else
#include <libgen.h>
#include <linux/limits.h>
#endif

#include "log.h"
#include "../system.h"
#include "math.h"

uint64_t get_timestamp(void)
{
	struct timeval tp;

	gettimeofday(&tp, 0);

	return (uint64_t)(unsigned long)tp.tv_sec << 10 | ((tp.tv_usec << 10) / 1000000);
}

static char g_xdag_current_path[4096] = {0};

void xdag_init_path(char *path)
{
#ifdef _WIN32
	char szPath[MAX_PATH];
	char szBuffer[MAX_PATH];
	char *pszFile;

	GetModuleFileName(NULL, (LPTCH)szPath, sizeof(szPath) / sizeof(*szPath));
	GetFullPathName((LPTSTR)szPath, sizeof(szBuffer) / sizeof(*szBuffer), (LPTSTR)szBuffer, (LPTSTR*)&pszFile);
	*pszFile = 0;

	strcpy(g_xdag_current_path, szBuffer);
#else
	char pathcopy[PATH_MAX] = {0};
	strcpy(pathcopy, path);
	char *prefix = dirname(pathcopy);
	if (*prefix != '/' && *prefix != '\\') {
		char buf[PATH_MAX] = {0};
		getcwd(buf, PATH_MAX);
		sprintf(g_xdag_current_path, "%s/%s", buf, prefix);
	} else {
		sprintf(g_xdag_current_path, "%s", prefix);
	}
#if defined (__MACOS__) || defined (__APPLE__)
	free(prefix);
#endif
#endif

	const size_t pathLen = strlen(g_xdag_current_path);
	if (pathLen == 0 || g_xdag_current_path[pathLen - 1] != *DELIMITER) {
		g_xdag_current_path[pathLen] = *DELIMITER;
		g_xdag_current_path[pathLen + 1] = 0;
	}
}

FILE* xdag_open_file(const char *path, const char *mode)
{
	char abspath[1024] = {0};
	sprintf(abspath, "%s%s", g_xdag_current_path, path);
	FILE* f = fopen(abspath, mode);
	return f;
}

void xdag_close_file(FILE *f)
{
	fclose(f);
}

int xdag_file_exists(const char *path)
{
	char abspath[1024] = {0};
	sprintf(abspath, "%s%s", g_xdag_current_path, path);
	struct stat st;
	return !stat(abspath, &st);
}

int xdag_mkdir(const char *path)
{
	char abspath[1024] = {0};
	sprintf(abspath, "%s%s", g_xdag_current_path, path);

#if defined(_WIN32)
	return _mkdir(abspath);
#else 
	return mkdir(abspath, 0770);
#endif	
}

long double log_difficulty2hashrate(long double log_diff)
{
	return ldexpl(expl(log_diff), -58)*(0.65);
}

void xdag_str_toupper(char *str)
{
	while(*str) {
		*str = toupper((unsigned char)*str);
		str++;
	}
}

void xdag_str_tolower(char *str)
{
	while(*str) {
		*str = tolower((unsigned char)*str);
		str++;
	}
}

char *xdag_basename(char *path)
{
#if defined(_WIN32)
	char *ptr;
	while((ptr = strchr(path, '/')) || (ptr = strchr(path, '\\'))) {
		path = ptr + 1;
	}
	return strdup(path);
#else
	return strdup(basename(path));
#endif
}

char *xdag_filename(char *_filename)
{
	char *filename = xdag_basename(_filename);
	char *ext = strchr(filename, '.');

	if(ext) {
		*ext = 0;
	}

	return filename;
}

// convert time to string representation
// minimal length of string buffer `buf` should be 60
void xdag_time_to_string(xdag_time_t time, char *buf)
{
	struct tm tm;
	char tmp[64];
	time_t t = time >> 10;
	localtime_r(&t, &tm);
	strftime(tmp, 60, "%Y-%m-%d %H:%M:%S", &tm);
	sprintf(buf, "%s.%03d", tmp, (int)((time & 0x3ff) * 1000) >> 10);
}

// convert time to string representation
// minimal length of string buffer `buf` should be 50
void time_to_string(time_t time, char* buf)
{
	struct tm tm;
	localtime_r(&time, &tm);
	strftime(buf, 50, "%Y-%m-%d %H:%M:%S", &tm);
}

size_t validate_remark(const char *str)
{
	return validate_ascii_safe(str, 33);// sizeof(xdag_remark_t) + 1
}

size_t validate_ascii_safe(const char *str, size_t maxsize)
{
	if (str == NULL) {
		return 0;
	}

	const char* start = str;
	const char* stop = str + maxsize;

	for (; str < stop; ++str) {
		if (*str < 32 || *str > 126) {
			if (*str == '\0') {
				return str - start;
			}
			return 0;
		}
	}

	return 0;
}

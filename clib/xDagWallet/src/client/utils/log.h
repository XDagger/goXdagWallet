/* logging, T13.670-T13.788 $DVS:time$ */

#ifndef XDAG_LOG_H
#define XDAG_LOG_H

enum xdag_debug_levels {
	XDAG_NOERROR,
	XDAG_FATAL,
	XDAG_CRITICAL,
	XDAG_INTERNAL,
	XDAG_ERROR,
	XDAG_WARNING,
	XDAG_MESSAGE,
	XDAG_INFO,
	XDAG_DEBUG,
	XDAG_TRACE,
	XDAG_MAX_LEVEL
};

#ifdef __cplusplus
extern "C" {
#endif
	
	extern int xdag_log(int level, int err, char* file, int line, const char *format, ...);

	extern char *xdag_log_array(const void *arr, unsigned size);

	extern int xdag_log_init(void);

	// sets the maximum error level for output to the log, returns the previous level (0 - do not log anything, 9 - all)
	extern int xdag_set_log_level(int level);
    
#ifdef __cplusplus
};
#endif

#define xdag_log_hash(hash) xdag_log_array(hash, sizeof(xdag_hash_t))

#define xdag_fatal(err, ...) xdag_log(XDAG_FATAL, err, __FILE__, __LINE__, __VA_ARGS__)
#define xdag_crit(err, ...)  xdag_log(XDAG_CRITICAL, err, __FILE__, __LINE__, __VA_ARGS__)
#define xdag_err(err, ...)   xdag_log(XDAG_ERROR , err, __FILE__, __LINE__, __VA_ARGS__)
#define xdag_warn(...)  xdag_log(XDAG_WARNING , 0, __FILE__, __LINE__, __VA_ARGS__)
#define xdag_mess(...)  xdag_log(XDAG_MESSAGE , 0, __FILE__, __LINE__, __VA_ARGS__)
#define xdag_info(...)  xdag_log(XDAG_INFO, 0, __FILE__, __LINE__, __VA_ARGS__)
#ifndef NDEBUG
#define xdag_debug(...) xdag_log(XDAG_DEBUG, 0, __FILE__, __LINE__, __VA_ARGS__)
#else
#define xdag_debug(...)
#endif

#endif

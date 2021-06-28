#ifndef XDAG_COMMANDS_H
#define XDAG_COMMANDS_H

#include <time.h>
#include "block.h"
#include "errno.h"

#define XDAG_COMMAND_MAX	0x100
#define XFER_MAX_IN		11

#ifdef __cplusplus
extern "C" {
#endif

	extern xdag_error_no processAccountCommand(char **out);
	extern xdag_error_no processAddressCommand(char **out);
	extern xdag_error_no processBalanceCommand(char **out);
	extern xdag_error_no processLevelCommand(const char *level, char **out);
	extern xdag_error_no processXferCommand(const char *address, const char *amount, const char *remark, char **out);
	extern xdag_error_no processStateCommand(char **out);
	extern xdag_error_no processExitCommand(void);
	extern xdag_error_no processHelpCommand(char **out);

	extern xdag_error_no xdag_do_xfer(const char *amount, const char *address, const char *remark, char **out);

	extern void xdag_log_xfer(xdag_hash_t from, xdag_hash_t to, xdag_amount_t amount);
    
#ifdef __cplusplus
};
#endif

#endif // !XDAG_COMMANDS_H

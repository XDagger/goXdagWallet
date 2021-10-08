//
//  common.h
//  xDagWallet
//
//  Created by Rui Xie on 7/11/18.
//  Copyright © 2018 xrdavies. All rights reserved.
//

#ifndef common_h
#define common_h

#include <time.h>
#include "block.h"
#include "system.h"
#include "errno.h"
#include "events.h"
#include "xdag_wrapper.h"

#define COINNAME "XDAG"

// This is for timeval redefinition issue for pthread
#define HAVE_STRUCT_TIMESPEC

// This is to disable security warning from Visual C++
#ifdef _MSC_VER
#define _CRT_SECURE_NO_WARNINGS
#endif

enum xdag_states
{
#define xdag_state(n,s) XDAG_STATE_##n ,
#include "state.h"
#undef xdag_state
};

extern struct xdag_stats
{
	xdag_diff_t difficulty, max_difficulty;
	uint64_t nblocks, total_nblocks;
	uint64_t nmain, total_nmain;
	uint32_t nhosts, total_nhosts, reserved1, reserved2;
} g_xdag_stats;

#define HASHRATE_LAST_MAX_TIME	(64 * 4) // numbers of main blocks in about 4H, to calculate the pool and network mean hashrate

extern struct xdag_ext_stats
{
	xdag_diff_t hashrate_total[HASHRATE_LAST_MAX_TIME];
	xdag_diff_t hashrate_ours[HASHRATE_LAST_MAX_TIME];
	xdag_time_t hashrate_last_time;
	uint64_t nnoref;
	uint64_t nhashes;
	double hashrate_s;
	uint32_t nwaitsync;
} g_xdag_extstats;

#define xdag_amount2xdag(amount) ((unsigned)((amount) >> 32))
#define xdag_amount2cheato(amount) ((unsigned)(((uint64_t)(unsigned)(amount) * 1000000000) >> 32))
#define ARG_EQUAL(a,b,c) strcmp(c, "") == 0 ? strcmp(a, b) == 0 : (strcmp(a, b) == 0 || strcmp(a, c) == 0)

#ifdef __cplusplus
extern "C" {
#endif
	//Default type of the block header
	//Test network and main network have different types of the block headers, so blocks from different networks are incompatible
	extern enum xdag_field_type g_block_header_type;

	/* the program state */
	extern int g_xdag_state;

	/* is there command 'run' */
	extern int g_xdag_run;

	/* 1 - the program works in a test network */
	extern int g_xdag_testnet;

	/* time of last transfer */
	extern time_t g_xdag_xfer_last;

	// convert cheato to xdag
	extern long double amount2xdags(xdag_amount_t amount);

	// contert xdag to cheato
	extern xdag_amount_t xdags2amount(const char *str);

	extern enum xdag_states xdag_get_state(void);

	extern void xdag_set_state(enum xdag_states state);

	extern const char *xdag_get_state_str(void);


	/* see dnet_user_crypt_action */
	extern int xdag_user_crypt_action(unsigned *data, unsigned long long data_id, unsigned size, int action);

	extern int xdag_wrapper_log(int level, xdag_error_no err, const char *msg);
	extern int xdag_wrapper_event(xdag_event_id event_id, xdag_error_no err, const char *msg);
    extern int xdag_wrapper_interact(xdag_event_id event_id, xdag_wrapper_msg *data);
#ifdef __cplusplus
}
#endif

#endif /* common_h */

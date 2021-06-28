//
//  common.c
//  xDagWallet
//
//  Created by Rui Xie on 7/11/18.
//  Copyright © 2018 xrdavies. All rights reserved.
//

#include "common.h"
#include <math.h>
#include "dnet_crypt.h"

enum xdag_field_type g_block_header_type = XDAG_FIELD_HEAD;
int g_xdag_state = XDAG_STATE_INIT;
int g_xdag_testnet = 0;
time_t g_xdag_xfer_last = 0;
struct xdag_stats g_xdag_stats;
struct xdag_ext_stats g_xdag_extstats;


xdag_amount_t xdags2amount(const char *str)
{
	long double sum;
	if(sscanf(str, "%Lf", &sum) != 1 || sum <= 0) {
		return 0;
	}
	long double flr = floorl(sum);
	xdag_amount_t res = (xdag_amount_t)flr << 32;
	sum -= flr;
	sum = ldexpl(sum, 32);
	flr = ceill(sum);
	return res + (xdag_amount_t)flr;
}

// convert xdag_amount_t to long double
long double amount2xdags(xdag_amount_t amount)
{
	return xdag_amount2xdag(amount) + (long double)xdag_amount2cheato(amount) / 1000000000;
}

enum xdag_states xdag_get_state(void)
{
	return g_xdag_state;
}

void xdag_set_state(enum xdag_states state)
{
	if(g_xdag_state != state) {
		g_xdag_state = state;
		xdag_wrapper_event(event_id_state_change, 0, xdag_get_state_str());
	}
}

const char *xdag_get_state_str()
{
	static const char *states[] = {
#define xdag_state(n,s) s ,
#include "state.h"
#undef xdag_state
	};
	return states[g_xdag_state];
}

/* see dnet_user_crypt_action */
int xdag_user_crypt_action(unsigned *data, unsigned long long data_id, unsigned size, int action)
{
	return dnet_user_crypt_action(data, data_id, size, action);
}

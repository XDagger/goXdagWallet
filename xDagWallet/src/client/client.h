#ifndef XDAG_MINER_H
#define XDAG_MINER_H

#include <stdio.h>
#include "block.h"
#include <pthread.h>

typedef struct  {
	struct xdag_field task[2], lastfield, minhash, nonce;
	xdag_time_t task_time;
	void *ctx0, *ctx;
} xdag_pool_task_t;

/* connecting the miner to pool pool_arg - pool parameters ip:port, testnet - 1 means testnet, 0 means mainnet*/
typedef struct {
    char pool_arg[256];
    int testnet;
} xdag_thread_param_t;

extern int g_xdag_client_running;

#ifdef __cplusplus
extern "C" {
#endif
	
	extern pthread_t g_client_thread;
	extern struct dfslib_crypt *g_crypt;

	/* client main thread */
	extern void *xdag_client_thread(void *arg);

	/* send block to network via pool */
	extern int xdag_send_block_via_pool(struct xdag_block *block);

	extern xdag_pool_task_t g_xdag_pool_task[2];
	extern uint64_t g_xdag_pool_task_index; /* global variables are instantiated with 0 */

#ifdef __cplusplus
};
#endif
		
#endif

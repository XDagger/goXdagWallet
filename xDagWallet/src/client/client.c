#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <fcntl.h>
#include <errno.h>
#include <pthread.h>

#if defined(_WIN32) || defined(_WIN64)

#include "../win/unistd.h"
#include "../win/winsockx.h"
#include <Winsock2.h>
#include <ws2tcpip.h>
#include <windows.h>
// need link with Ws2_32.lib
#pragma comment(lib, "Ws2_32.lib")

#else
#include <unistd.h>
#include <sys/socket.h>
#include <sys/ioctl.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <netdb.h>
#include "system.h"

#endif
#include "../dus/dfslib_crypt.h"
#include "../dus/crc.h"
#include "crypt.h"
#include "wallet.h"
#include "address.h"
#include "block.h"
#include "common.h"
#include "client.h"
#include "storage.h"
#include "utils/log.h"
#include "commands.h"
#include "dnet_crypt.h"
#include "./utils/utils.h"
#include "version.h"

//#define __stand_alone_lib__ // if run as stand alone library

#if defined(_WIN32) || defined(_WIN64)
#if defined(_WIN64)
#define poll WSAPoll
#else
#define poll(a, b, c) ((a)->revents = (a)->events, (b))
#endif
#else
#include <poll.h>
#endif

#define DATA_SIZE          (sizeof(struct xdag_field) / sizeof(uint32_t))
#define BLOCK_HEADER_WORD  0x3fca9e2bu

time_t g_xdag_last_received = 0;
pthread_t g_client_thread;

xdag_pool_task_t g_xdag_pool_task[2];
uint64_t g_xdag_pool_task_index;

struct dfslib_crypt *g_crypt;

#define MINERS_PWD             "minersgonnamine"
#define SECTOR0_BASE           0x1947f3acu
#define SECTOR0_OFFSET         0x82e9d1b5u
#define SEND_PERIOD            5                                  /* share period of sending shares */

struct miner {
	struct xdag_field id;
	uint64_t nfield_in;
	uint64_t nfield_out;
};

static struct miner g_local_miner;

static int g_socket = -1;

static int crypt_start(void)
{
	struct dfslib_string str;
	uint32_t sector0[128];
	int i;

	g_crypt = malloc(sizeof(struct dfslib_crypt));
	if(!g_crypt) return -1;
	dfslib_crypt_set_password(g_crypt, dfslib_utf8_string(&str, MINERS_PWD, strlen(MINERS_PWD)));

	for(i = 0; i < 128; ++i) {
		sector0[i] = SECTOR0_BASE + i * SECTOR0_OFFSET;
	}

	for(i = 0; i < 128; ++i) {
		dfslib_crypt_set_sector0(g_crypt, sector0);
		dfslib_encrypt_sector(g_crypt, sector0, SECTOR0_BASE + i * SECTOR0_OFFSET);
	}

	return 0;
}

static int can_send_share(time_t current_time, time_t task_time, time_t share_time)
{
	int can_send = (current_time - share_time >= SEND_PERIOD) && (current_time - task_time <= 64) && (share_time < task_time);
	return can_send;
}

static int send_to_pool(struct xdag_field *fld, int nfld)
{
	struct xdag_field f[XDAG_BLOCK_FIELDS];
	xdag_hash_t h;
	struct miner *m = &g_local_miner;
	int todo = nfld * sizeof(struct xdag_field), done = 0;

	if(g_socket < 0) {
		return -1;
	}

	memcpy(f, fld, todo);

	if(nfld == XDAG_BLOCK_FIELDS) {
		f[0].transport_header = 0;

		xdag_hash(f, sizeof(struct xdag_block), h);

		f[0].transport_header = BLOCK_HEADER_WORD;

		uint32_t crc = crc_of_array((uint8_t*)f, sizeof(struct xdag_block));

		f[0].transport_header |= (uint64_t)crc << 32;
	}

	for(int i = 0; i < nfld; ++i) {
		dfslib_encrypt_array(g_crypt, (uint32_t*)(f + i), DATA_SIZE, m->nfield_out++);
	}

	while(todo) {
		struct pollfd p;

		p.fd = g_socket;
		p.events = POLLOUT;

		if(!poll(&p, 1, 1000)) continue;

		if(p.revents & (POLLHUP | POLLERR)) {
			return -1;
		}

		if(!(p.revents & POLLOUT)) continue;

		int res = (int)write(g_socket, (uint8_t*)f + done, todo);
		if(res <= 0) {
			return -1;
		}

		done += res;
		todo -= res;
	}

	if(nfld == XDAG_BLOCK_FIELDS) {
		xdag_info("Sent  : %016llx%016llx%016llx%016llx t=%llx res=%d",
			h[3], h[2], h[1], h[0], fld[0].time, 0);
	}

	return 0;
}

static int client_init(void)
{
	memset(&g_xdag_stats, 0, sizeof(g_xdag_stats));
	memset(&g_xdag_extstats, 0, sizeof(g_xdag_extstats));

    if(g_xdag_testnet) {
        g_block_header_type = XDAG_FIELD_HEAD_TEST; //block header has the different type in the test network
    }

	xdag_mess("Starting xdag, version %s", XDAG_VERSION);
	xdag_mess("Starting dnet transport...");

    if (dnet_crypt_init(DNET_VERSION)) {
		sleep(3);
		xdag_wrapper_event(event_id_err_exit, error_pwd_incorrect, "Password incorrect.\n");
		return -1;
	}

	if (xdag_log_init()) {
		xdag_wrapper_event(event_id_err_exit, error_init_log, "Init log failed.\n");
		return -1;
	}

	xdag_mess("Initializing cryptography...");
	if (xdag_crypt_init(1)) {
		xdag_wrapper_event(event_id_err_exit, error_init_crypto, "Init crypto failed.\n");
		return -1;
	}

	xdag_mess("Reading wallet...");
	if (xdag_wallet_init()) {
		xdag_wrapper_event(event_id_err_exit, error_init_wallet, "Init wallet failed.\n");
		return -1;
	}

	xdag_mess("Initializing addresses...");
	if (xdag_address_init()) {
		xdag_wrapper_event(event_id_err_exit, error_init_address, "Init wallet failed.\n");
		return -1;
	}

	xdag_mess("Starting blocks engine...");
	if (xdag_blocks_start()) {
		xdag_wrapper_event(event_id_err_exit, error_init_block, "load blocks failed.\n");
		return -1;
	}

	//	if(is_rpc) {
	//		xdag_mess("Initializing RPC service...");
	//		if(!!xdag_rpc_service_init(rpc_port)) return -1;
	//	}

	for(int i = 0; i < 2; ++i) {
		g_xdag_pool_task[i].ctx0 = malloc(xdag_hash_ctx_size());
		g_xdag_pool_task[i].ctx = malloc(xdag_hash_ctx_size());

		if(!g_xdag_pool_task[i].ctx0 || !g_xdag_pool_task[i].ctx) {
			xdag_wrapper_event(event_id_err_exit, error_init_task, "Init task failed.\n");
			return -1;
		}
	}

	if(crypt_start()) {
		xdag_wrapper_event(event_id_err_exit, error_start_crypto, "Crypt start failed.\n");
		return -1;
	}

	return 0;
}

static void client_thread_cleanup()
{
	xdag_debug("client thread clean up called ");
    xdag_wallet_finish();
    xdag_storage_finish();
    xdag_blocks_finish();
	xdag_debug(" work thread clean up finished ");
}

void *xdag_client_thread(void *arg)
{
    int oldcancelstate;
    int oldcanceltype;
    pthread_setcancelstate(PTHREAD_CANCEL_ENABLE, &oldcancelstate);
    pthread_setcanceltype(PTHREAD_CANCEL_DEFERRED, &oldcanceltype);
    
    pthread_cleanup_push(client_thread_cleanup, NULL);
    
    xdag_error_no err_no = error_none;
    char *err_mess = NULL;
    
    int err = pthread_detach(pthread_self());
    if(err != 0) {
        err_no = error_unknown;
        err_mess = "Detach xdag_client_thread failed.";
        goto end;
    }
    
    xdag_thread_param_t *param = (xdag_thread_param_t *)arg;
	if(!param) {
        err_no = error_missing_param;
        err_mess = "Missing parameters.";
        goto end;
	}
    
	char pool_param[256];
	strcpy(pool_param, param->pool_arg);
    
    g_xdag_testnet = param->testnet;
    
    xdag_mess("testnet %d", g_xdag_testnet);

	if(!!client_init()) {
        pthread_exit((void *)event_id_err_exit);
        return 0;
    } else {
        xdag_wrapper_event(event_id_init_done, error_none, "");
    }

	xdag_mess("Initialize miner...");

	struct xdag_block b;
	struct xdag_field data[2];
	xdag_hash_t hash;
	xdag_time_t t;

	struct sockaddr_in peeraddr;
	char *lasts;
	int res = 0, reuseaddr = 1;
	struct linger linger_opt = { 1, 0 }; // Linger active, timeout 0

	xdag_mess("Entering main cycle...");
	char pool_arg[0x100];

begin:
	strcpy(pool_arg, pool_param);
	memset(&g_local_miner, 0, sizeof(struct miner));
	xdag_get_our_block(g_local_miner.id.data);

	struct miner *m = &g_local_miner;
	m->nfield_in = m->nfield_out = 0;

	memcpy(hash, g_local_miner.id.data, sizeof(xdag_hash_t));

	int ndata = 0;
	int maxndata = sizeof(struct xdag_field);
	time_t share_time = 0;
	time_t task_time = 0;

	const int64_t pos = xdag_get_block_pos(hash, &t);
	if(pos < 0) {
        err_no = error_block_not_found;
        err_mess = "Cann't find the block";
        goto end;
	}
    
	struct xdag_block *blk = xdag_storage_load(hash, t, pos, &b);
	if(!blk) {
        err_no = error_storage_load_faild;
        err_mess = "Cann't load the block";
        goto end;
	}
	if(blk != &b) {
		memcpy(&b, blk, sizeof(struct xdag_block));
	}
    
	// Create a socket
#if defined(_WIN32) || defined(_WIN64)
	WSADATA mainSdata;
	if (WSAStartup(2.2, &mainSdata) == 0)
	{
		g_socket = WSASocket(AF_INET, SOCK_STREAM, IPPROTO_TCP, NULL, 0, NULL);
	}
	else
	{
		g_socket = INVALID_SOCKET;
	}
#else
	g_socket = socket(AF_INET, SOCK_STREAM, IPPROTO_TCP);
#endif
	
	if(g_socket == INVALID_SOCKET) {
        err_no = error_socket_create;
        err_mess = "Cann't create a socket.";
        goto end;
	}
    
	if(fcntl(g_socket, F_SETFD, FD_CLOEXEC) == -1) {
		xdag_err(error_socket_create, "pool  : Cann't set FD_CLOEXEC flag on socket %d, %s\n", g_socket, strerror(errno));
	}
    
	// Fill in the address of server
	memset(&peeraddr, 0, sizeof(peeraddr));
	peeraddr.sin_family = AF_INET;

	// Resolve the server address (convert from symbolic name to IP number)
	const char *s = strtok_r(pool_arg, " \t\r\n:", &lasts);
	if(!s) {
        err_no = error_missing_param;
        err_mess = "Host is not given.";
        goto end;
	}
    
	if(!strcmp(s, "any")) {
		peeraddr.sin_addr.s_addr = htonl(INADDR_ANY);
	} else if(
#if defined(_WIN64)
		!inet_pton(AF_INET, s, &peeraddr.sin_addr)
#elif defined(_WIN32)
		!inet_pton_32(AF_INET, s, &peeraddr.sin_addr)
#else
		!inet_aton(s, &peeraddr.sin_addr)
#endif
		) {
		struct hostent *host = gethostbyname(s);
		if(host == NULL || host->h_addr_list[0] == NULL) {
            err_no = error_socket_resolve_host;
            err_mess = "Cann't resolve host.";
            goto end;
		}
		// Write resolved IP address of a server to the address structure
		memmove(&peeraddr.sin_addr.s_addr, host->h_addr_list[0], 4);
	}

	// Resolve port
	s = strtok_r(0, " \t\r\n:", &lasts);
	if(!s) {
        err_no = error_missing_param;
        err_mess = "Port is not given.";
        goto end;
	}
    
	peeraddr.sin_port = htons(atoi(s));

	// Set the "LINGER" timeout to zero, to close the listen socket
	// immediately at program termination.
	setsockopt(g_socket, SOL_SOCKET, SO_LINGER, (char*)&linger_opt, sizeof(linger_opt));
	setsockopt(g_socket, SOL_SOCKET, SO_REUSEADDR, (char*)&reuseaddr, sizeof(int));
    
    xdag_set_state(g_xdag_testnet ? XDAG_STATE_TTST : XDAG_STATE_TRYP);
    
	// Now, connect to a pool
	res = connect(g_socket, (struct sockaddr*)&peeraddr, sizeof(peeraddr));
	if(res) {
		xdag_err(error_socket_connect, "Cann't connect to the pool");
		xdag_set_state(g_xdag_testnet ? XDAG_STATE_TTST : XDAG_STATE_TRYP);
		goto err;
	}

	if(send_to_pool(b.field, XDAG_BLOCK_FIELDS) < 0) {
		xdag_err(error_socket_closed, "Socket is closed");
		xdag_set_state(g_xdag_testnet ? XDAG_STATE_TTST : XDAG_STATE_TRYP);
		goto err;
	}
    
	for(;;) {
		if(get_timestamp() - t > 1024) {
			t = get_timestamp();
			if (xdag_get_state() == XDAG_STATE_REST) {
				xdag_err(error_block_load_failed, "Block reset!!!");
			} else {
				if (t > (g_xdag_last_received << 10) && t - (g_xdag_last_received << 10) > 3 * MAIN_CHAIN_PERIOD) {
					xdag_set_state(g_xdag_testnet ? XDAG_STATE_TTST : XDAG_STATE_TRYP);
				} else {
					if (t - (g_xdag_xfer_last << 10) <= 2 * MAIN_CHAIN_PERIOD + 4) {
						xdag_set_state(XDAG_STATE_XFER);
					} else {
						xdag_set_state(g_xdag_testnet ? XDAG_STATE_PTST : XDAG_STATE_POOL);
					}
				}
			}
		}

		struct pollfd p;

		if(g_socket < 0) {
			xdag_err(error_socket_closed, "socket is closed");
			goto err;
		}

		p.fd = g_socket;
		time_t current_time = time(0);
		p.events = POLLIN | (can_send_share(current_time, task_time, share_time) ? POLLOUT : 0);
        
		if(!poll(&p, 1, 0)) {
			sleep(1);
			continue;
		}
        
		if(p.revents & POLLIN) {
			res = (int)read(g_socket, (uint8_t*)data + ndata, maxndata - ndata);
			if(res < 0) {
				xdag_err(error_socket_read, "read error on socket");
				goto err;
			}
			ndata += res;
			if(ndata == maxndata) {
				struct xdag_field *last = data + (ndata / sizeof(struct xdag_field) - 1);

				dfslib_uncrypt_array(g_crypt, (uint32_t*)last->data, DATA_SIZE, m->nfield_in++);
				xdag_info("My Hash  : %016llx%016llx%016llx%016llx", hash[3], hash[2], hash[1], hash[0]);
				xdag_info("Received Hash  : %016llx%016llx%016llx%016llx", last->data[3], last->data[2], last->data[1], last->data[0]);

				if(!memcmp(last->data, hash, sizeof(xdag_hashlow_t))) {
					xdag_set_balance(hash, last->amount);
					g_xdag_last_received = current_time;
					ndata = 0;

					maxndata = sizeof(struct xdag_field);
				} else if(maxndata == 2 * sizeof(struct xdag_field)) {
					const uint64_t task_index = g_xdag_pool_task_index + 1;
                    xdag_pool_task_t *task = &g_xdag_pool_task[task_index & 1];

					task->task_time = xdag_main_time();
					xdag_hash_set_state(task->ctx, data[0].data,
						sizeof(struct xdag_block) - 2 * sizeof(struct xdag_field));
					xdag_hash_update(task->ctx, data[1].data, sizeof(struct xdag_field));
					xdag_hash_update(task->ctx, hash, sizeof(xdag_hashlow_t));

					dnet_generate_random_array(task->nonce.data, sizeof(xdag_hash_t));

					memcpy(task->nonce.data, hash, sizeof(xdag_hashlow_t));
					memcpy(task->lastfield.data, task->nonce.data, sizeof(xdag_hash_t));

					xdag_hash_final(task->ctx, &task->nonce.amount, sizeof(uint64_t), task->minhash.data);

					g_xdag_pool_task_index = task_index;
					task_time = time(0);

					xdag_info("Task  : t=%llx N=%llu", task->task_time << 16 | 0xffff, task_index);

					ndata = 0;
					maxndata = sizeof(struct xdag_field);
				} else {
					maxndata = 2 * sizeof(struct xdag_field);
				}
			}
		} else if(p.revents & POLLOUT) {
			const uint64_t task_index = g_xdag_pool_task_index;
            xdag_pool_task_t *task = &g_xdag_pool_task[task_index & 1];
			uint64_t *h = task->minhash.data;

			share_time = time(0);
			res = send_to_pool(&task->lastfield, 1);

			xdag_info("Share : %016llx%016llx%016llx%016llx t=%llx res=%d",
				h[3], h[2], h[1], h[0], task->task_time << 16 | 0xffff, res);

			if(res) {
				xdag_err(error_socket_write, "write error on socket");
				goto err;
			}
        } else {
            if(p.revents & POLLHUP) {
                xdag_err(error_socket_hangup, "socket hangup");
                goto err;
            }
            
            if(p.revents & POLLERR) {
                xdag_err(error_socket_err, "socket error");
                goto err;
            }
        }
	}

err:
	if(g_socket != INVALID_SOCKET) {
		close(g_socket);
		g_socket = INVALID_SOCKET;
	}
	sleep(5);
	goto begin;

end:
	if(g_socket != INVALID_SOCKET) {
		close(g_socket);
		g_socket = INVALID_SOCKET;
	}
    
    if (err_no != error_none || err_mess != NULL) {
        xdag_wrapper_event(event_id_err_exit, err_no, err_mess);
    }
    
    pthread_cleanup_pop(0);
    pthread_exit(0);
    return 0;
}

/* send block to network via pool */
int xdag_send_block_via_pool(struct xdag_block *b) {
	if(g_socket < 0) return -1;
	int ret = send_to_pool(b->field, XDAG_BLOCK_FIELDS);
	return ret;
}

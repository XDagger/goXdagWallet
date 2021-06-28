//
//  errno.h
//  xDagWallet
//
//  Created by Rui Xie on 7/12/18.
//  Copyright © 2018 xrdavies. All rights reserved.
//

#ifndef errno_h
#define errno_h

#include <stdio.h>
typedef enum {
	// no error
	error_none						= 0x0000,

	// general errors
	error_unknown					= 0x0001,
	error_malloc					= 0x0002,
	error_fatal						= 0x0003,
	error_invalid_command			= 0x0004,

	// init errors
	error_pwd_inconsistent			= 0x1001,
	error_pwd_incorrect				= 0x1002,
	error_create_dnet_key_failed	= 0x1003,
	error_write_wallet_key_failed	= 0x1004,
	error_add_wallet_key_failed		= 0x1005,
	error_init_task_failed			= 0x1006,
	error_init_crypt_failed			= 0x1007,
	error_missing_param				= 0x1008,
	error_not_ready					= 0x1009,
	error_incorrect_level			= 0x100A,
	error_init_log					= 0x100B,
	error_init_crypto				= 0x100C,
	error_init_wallet				= 0x100D,
	error_init_address				= 0x100E,
	error_init_block				= 0x1010,
	error_init_task					= 0x1011,
	error_start_crypto				= 0x1012,

	// storage
	error_storage_load_faild		= 0x2001,
	error_storage_sum_corrupted		= 0x2002,
	error_storage_create_file		= 0x2003,
	error_storage_write_file		= 0x2004,
	error_storage_corrupted			= 0x2005,

	// xfer
	error_xfer_nothing				= 0x3001,
	error_xfer_too_small			= 0x3002,
	error_xfer_incorrect_address	= 0x3003,
	error_xfer_no_address			= 0x3004,
	error_xfer_no_amount			= 0x3005,
	error_xfer_not_ready			= 0x3006,
	error_xfer_make_failed			= 0x3007,
	error_xfer_no_remark			= 0x3008,
	error_xfer_incorrect_remark		= 0x3009,

	// block
	error_block_create				= 0x4001,
	error_block_not_found			= 0x4002,
	error_block_load_failed			= 0x4003,

	// network
	error_socket_create				= 0x5001,
	error_socket_host				= 0x5002,
	error_socket_port				= 0x5003,
	error_socket_connect			= 0x5004,
	error_socket_hangup				= 0x5005,
	error_socket_err				= 0x5006,
	error_socket_read				= 0x5007,
	error_socket_write				= 0x5008,
	error_socket_timeout			= 0x5009,
	error_socket_resolve_host		= 0x5010,
	error_socket_closed				= 0x5011

} xdag_error_no;

#endif /* errno_h */

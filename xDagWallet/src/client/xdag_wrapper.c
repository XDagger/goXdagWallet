//
//  xdag_wrapper.c
//  xDagWallet
//
//  Created by Rui Xie on 7/12/18.
//  Copyright © 2018 xrdavies. All rights reserved.
//

#include "xdag_wrapper.h"
#include <stdlib.h>
#include <string.h>
#include "common.h"
#include "commands.h"

#if !defined(NDBUG)
#define assert(p)    if(!(p)){fprintf(stderr,\
"Assertion failed: %s, file %s, line %d\n",\
#p, __FILE__, __LINE__);abort();}
#else
#define assert(p)
#endif


xdag_log_callback_t g_wrapper_log_callback = NULL;
xdag_event_callback_t g_wrapper_event_callback = NULL;
static void* g_thisObj = NULL;

int xdag_set_log_callback(xdag_log_callback_t callback)
{
	g_wrapper_log_callback = callback;
	return 0;
}

int xdag_set_event_callback(xdag_event_callback_t callback)
{
	g_wrapper_event_callback = callback;
	return 0;
}

int xdag_set_password_callback(xdag_password_callback_t callback)
{
	return xdag_user_crypt_action((uint32_t *)(void *)callback, 0, 0, 6);
}

int xdag_wrapper_init(void* thisObj, xdag_password_callback_t password, xdag_event_callback_t event)
{
	if(thisObj) g_thisObj = thisObj;
//    if(password) xdag_set_password_callback(password);
	if(event) xdag_set_event_callback(event);

	return 0;
}

int xdag_wrapper_init_client(const char *args)
{
	return 0;
}

int xdag_wrapper_xfer(const char *amount, const char *to, const char *remark)
{
	char *result = NULL;
	int err = processXferCommand(amount, to, remark, &result);

	if(err != error_none) {
		xdag_wrapper_event(event_id_promot, err, result);
	} else {
		xdag_wrapper_event(event_id_xfer_done, 0, result);
	}

	if(result) {
		free(result);
	}
	return err;
}

int xdag_wrapper_account(void)
{
	char *result = NULL;
	int err = processAccountCommand(&result);

	if(err != error_none) {
		xdag_wrapper_event(event_id_promot, err, result);
	} else {
		xdag_wrapper_event(event_id_account_done, 0, result);
	}

	if(result) {
		free(result);
	}
	return err;
}

int xdag_wrapper_address(void)
{
	char *result = NULL;
	int err = processAddressCommand(&result);

	if(err != error_none) {
		xdag_wrapper_event(event_id_promot, err, result);
	} else {
		xdag_wrapper_event(event_id_address_done, 0, result);
	}

	if(result) {
		free(result);
	}
	return err;
}

int xdag_wrapper_balance(void)
{
	char *result = NULL;
	int err = processBalanceCommand(&result);

	if(err != error_none) {
		xdag_wrapper_event(event_id_promot, err, result);
	} else {
		xdag_wrapper_event(event_id_balance_done, 0, result);
	}

	if(result) {
		free(result);
	}
	return err;
}

int xdag_wrapper_level(const char *level)
{
	char *result = NULL;
	int err = processLevelCommand(level, &result);
	if(err != error_none) {
		xdag_wrapper_event(event_id_promot, err, result);
	} else {
		xdag_wrapper_event(event_id_level_done, 0, result);
	}

	if(result) {
		free(result);
	}

	return err;
}

int xdag_wrapper_state(void)
{
	char *result = NULL;
	int err = processStateCommand(&result);

	if(err != error_none) {
		xdag_wrapper_event(event_id_promot, err, result);
	} else {
		xdag_wrapper_event(event_id_state_done, 0, result);
	}

	if(result) {
		free(result);
	}
	return err;
}

int xdag_wrapper_exit(void)
{
	return processExitCommand();
}

int xdag_wrapper_help(void)
{
    char *result = NULL;
    
    processHelpCommand(&result);
    
    if(result) {
        xdag_wrapper_event(event_id_promot, error_none, result);
        free(result);
    }
    return 0;
}


int xdag_wrapper_log(int level, xdag_error_no err, const char *data)
{
    if (err == error_none) {
        xdag_wrapper_event(event_id_log, err, data);
    } else {
        xdag_wrapper_event(event_id_err, err, data);
    }
    
	return 0;
}

int xdag_wrapper_interact(xdag_event_id event_id, xdag_wrapper_msg *data)
{
    if(!g_wrapper_event_callback) {
        assert(0);
    } else {
        xdag_event *evt = calloc(1, sizeof(xdag_event));
        evt->event_id = event_id;
        evt->error_no = error_none;
        evt->event_data = (void *)data;
        
        if (g_wrapper_event_callback) {
            (*g_wrapper_event_callback)(g_thisObj, evt);
        }
        
        free(evt);
    }
	return 0;
}

int xdag_wrapper_event(xdag_event_id event_id, xdag_error_no err,  const char *msg)
{
	if(!g_wrapper_event_callback) {
		assert(0);
	} else {
		xdag_event *evt = calloc(1, sizeof(xdag_event));
		evt->event_id = event_id;
		evt->error_no = err;
        evt->event_data = msg == NULL? strdup(""): strdup(msg);

		if (g_wrapper_event_callback) {
			(*g_wrapper_event_callback)(g_thisObj, evt);
		}

		if(evt->event_data) {
			free(evt->event_data);
		}
		free(evt);
	}

	return 0;
}




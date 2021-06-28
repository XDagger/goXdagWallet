//
//  xdag_wrapper.h
//  xDagWallet
//
//  Created by Rui Xie on 7/12/18.
//  Copyright © 2018 xrdavies. All rights reserved.
//

#ifndef xdag_wrapper_h
#define xdag_wrapper_h

#include <stdio.h>
#include "errno.h"
#include "events.h"

#define XDAG_WRAPPER_MSG_LEN 256
typedef struct {
    char msg[XDAG_WRAPPER_MSG_LEN];
} xdag_wrapper_msg;

#ifdef __cplusplus
extern "C" {
#endif

    typedef int(*xdag_log_callback_t)(int, xdag_error_no, char *);
    typedef int(*xdag_event_callback_t)(void *, xdag_event *) ;
    typedef int(*xdag_password_callback_t)(const char *prompt, char *buf, unsigned size);


    extern int xdag_wrapper_init(void* thisObj, xdag_password_callback_t password, xdag_event_callback_t event);
    extern int xdag_wrapper_init_client(const char *args);

    extern int xdag_wrapper_xfer(const char *amount, const char *to, const char *remark);
    extern int xdag_wrapper_account(void);
    extern int xdag_wrapper_address(void);
    extern int xdag_wrapper_balance(void);
    extern int xdag_wrapper_level(const char *level);
    extern int xdag_wrapper_state(void);
    extern int xdag_wrapper_exit(void);
    extern int xdag_wrapper_help(void);

    extern int xdag_set_event_callback(xdag_event_callback_t callback);
    extern int xdag_wrapper_exit(void);

#ifdef __cplusplus
}
#endif

#endif /* xdag_wrapper_h */

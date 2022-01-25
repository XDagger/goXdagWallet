//
// Created by swordlet on 2021/3/24.
//

#ifndef DAGWALLET_XDAG_RUNTIME_H
#define DAGWALLET_XDAG_RUNTIME_H


#define _TIMESPEC_DEFINED

#include "xDagWallet/src/client/common.h"
#include "xDagWallet/src/client/commands.h"
#include "xDagWallet/src/client/client.h"
#include "xDagWallet/src/client/events.h"
#include "xDagWallet/src/client/utils/utils.h"
#include "xDagWallet/src/client/address.h"
#include "xDagWallet/src/client/dnet_crypt.h"
#include "xDagWallet/src/client/xdag_wrapper.h"

////---- Duplicated from dnet_crypt.c ----
#define KEYFILE	    "dnet_key.dat"
struct dnet_keys {
    struct dnet_key priv;
    struct dnet_key pub;
};
////------------------------------------


////---- Duplicated from commands.c ----
struct account_callback_data {
    char out[128];
    int count;
};

struct xfer_callback_data {
    struct xdag_field fields[XFER_MAX_IN + 1];
    int keys[XFER_MAX_IN + 1];
    xdag_amount_t todo, done, remains;
    int fieldsCount, keysCount, outsig;
    xdag_hash_t transactionBlockHash;
};
typedef int(*event_callback)(void*, xdag_event *);
typedef int(*password_callback)(const char *prompt, char *buf, unsigned size);
////------------------------------------
#ifdef __cplusplus
extern "C" {
#endif
//extern int xdag_event_callback(void* thisObj, xdag_event *event);


////---- Exporting functions ----
extern int xdag_init_wrap(int argc, char **argv, const char * pool_address, int testnet);
extern int xdag_set_password_callback_wrap(password_callback callback);
extern int xdag_set_event_callback_wrap(event_callback callback);
extern int xdag_get_state_wrap(void);
extern int xdag_get_balance_wrap(void);
extern int xdag_get_address_wrap(void);
extern int xdag_exit_wrap(void);


extern int xdag_transfer_wrap(const char* toAddress, const char* amountString, const char* remarkString);
extern int xdag_is_valid_wallet_address(const char* address);
extern int xdag_dnet_crpt_found();
extern int xdag_is_valid_remark(const char* remark);

#ifdef __cplusplus
};
#endif
#endif //DAGWALLET_XDAG_RUNTIME_H

//
// Created by swordlet on 2021/3/24.
//
#if defined(_WIN32) || defined(_WIN64)
#else
#include <unistd.h>
#endif
#include <cstdlib>
#include <cstring>
#include "xdag_runtime.h"

////------------------------------------
extern pthread_t g_client_thread;
static int g_client_init_done = 0;
////---- Exporting functions wrapping functions ----
int xdag_init_wrap(int argc, char **argv, const char *pool_address, int testnet)
{
    xdag_init_path(argv[0]);

    ////const char *pool_arg = "de1.xdag.org:13654";

    ////xdag_set_event_callback(&xdag_event_callback);

    xdag_thread_param_t param;
    strncpy(param.pool_arg, pool_address, 255);
    param.testnet = testnet;

    int err = pthread_create(&g_client_thread, 0, xdag_client_thread, (void *)&param);
    if (err != 0)
    {
        printf("create client_thread failed, error : %s\n", strerror(err));
        return -1;
    }
    while (!g_client_init_done)
    {
        sleep(1);
    }

    return 0;
}

int xdag_set_event_callback_wrap(event_callback callback)
{
    return xdag_set_event_callback(callback);
}

int xdag_get_state_wrap(void)
{
    xdag_wrapper_state();
    return 0;
}

int xdag_get_balance_wrap(void)
{
    xdag_wrapper_balance();
    return 0;
}

int xdag_get_address_wrap(void)
{
    xdag_wrapper_address();
    return 0;
}

int xdag_exit_wrap(void)
{
    xdag_wrapper_exit();
    return pthread_cancel(g_client_thread);
}

// int xdag_event_callback(void* thisObj, xdag_event *event)
//{
//     if (!event) {
//         return -1;
//     }
//
//     switch (event->event_id) {
//         case event_id_init_done:
//         {
//             g_client_init_done = 1;
//             break;
//         }
//         case event_id_log:
//         {
//
//             break;
//         }
//
//         case event_id_interact:
//         {
//
//             break;
//         }
//
//             //		case event_id_err:
//             //		{
//             //			fprintf(stdout, "error : %x, msg : %s\n", event->error_no, event->event_data);
//             //			fflush(stdout);
//             //			break;
//             //		}
//
//         case event_id_err_exit:
//         {
//
//             break;
//         }
//
//         case event_id_account_done:
//         {
//
//             break;
//         }
//
//         case event_id_address_done:
//         {
//
//             break;
//         }
//
//         case event_id_balance_done:
//         {
//
//             break;
//         }
//
//         case event_id_xfer_done:
//         {
//
//             break;
//         }
//
//         case event_id_level_done:
//         {
//
//             break;
//         }
//
//         case event_id_state_done:
//         {
//
//             break;
//         }
//
//         case event_id_exit_done:
//         {
//
//             break;
//         }
//
//         case event_id_passwd:
//         {
//             break;
//         }
//
//         case event_id_set_passwd:
//         {
//             break;
//         }
//
//         case event_id_set_passwd_again:
//         {
//             break;
//         }
//
//         case event_id_random_key:
//         {
//             break;
//         }
//
//         case event_id_state_change:
//         {
//
//             break;
//         }
//
//         default:
//         {
//
//             break;
//         }
//     }
//     return 0;
// }

int xdag_set_password_callback_wrap(password_callback callback)
{
    //// return xdag_set_password_callback(callback);
    return xdag_user_crypt_action((uint32_t *)(callback), 0, 0, 6);
}

int xdag_transfer_wrap(const char *toAddress, const char *amountString, const char *remarkString)
{
    char *address = (char *)malloc(strlen(toAddress) + 1);
    char *amount = (char *)malloc(strlen(amountString) + 1);
    char *remark = (char *)malloc(strlen(remarkString) + 1);

    strcpy(address, toAddress);
    address[strlen(toAddress)] = 0;
    strcpy(amount, amountString);
    amount[strlen(amountString)] = 0;
    strcpy(remark, remarkString);
    remark[strlen(remarkString)] = 0;

    char *result = NULL;
    int err = processXferCommand(amount, address, remark, &result);

    if (err != error_none)
    {
        xdag_wrapper_event(event_id_log, (xdag_error_no)err, result);
    }
    else
    {
        xdag_wrapper_event(event_id_xfer_done, (xdag_error_no)0, result);
    }

    free(address);
    free(amount);
    free(remark);

    return err;
}

int xdag_is_valid_wallet_address(const char *address)
{
    struct xfer_callback_data xfer;
    if (xdag_address2hash(address, xfer.fields[XFER_MAX_IN].hash) == 0)
    {
        return 0;
    }
    else
    {
        return -1;
    }
}

int xdag_dnet_crpt_found()
{
    FILE *f = NULL;
    struct dnet_keys *keys = (struct dnet_keys *)malloc(sizeof(struct dnet_keys));

    int is_found = -1;
    f = xdag_open_file(KEYFILE, "rb");

    if (f)
    {
        if (fread(keys, sizeof(struct dnet_keys), 1, f) == 1)
        {
            is_found = 0;
        }
        xdag_close_file(f);
    }

    free(keys);
    return is_found;
}

int xdag_is_valid_remark(const char *remark)
{
    size_t s = validate_remark(remark);
    if (s < 1 || s > 33)
    {
        return -1;
    }
    else
    {
        return 0;
    }
}

void *xdag_get_default_key()
{
    return xdag_default_key();
}

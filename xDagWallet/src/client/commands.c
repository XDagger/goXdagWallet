#include "commands.h"
#include <string.h>
#include <math.h>
#include <stdlib.h>
#include <ctype.h>
#include "common.h"
#include "address.h"
#include "wallet.h"
#include "utils/log.h"
#include "utils/utils.h"
#include "client.h"
#include "crypt.h"
#include "client.h"
#include "storage.h"
#include "errno.h"

#if !defined(_WIN32) && !defined(_WIN64)
#include <unistd.h>
#endif

#define Nfields(d) (2 + d->hasRemark + d->fieldsCount + 3 * d->keysCount + 2 * d->outsig)

struct account_callback_data {
	char out[128];
	int count;
};

struct xfer_callback_data {
	struct xdag_field fields[XFER_MAX_IN + 1];
	int keys[XFER_MAX_IN + 1];
	xdag_amount_t todo, done, remains;
	int fieldsCount, keysCount, outsig, hasRemark;
	xdag_hash_t transactionBlockHash;
	xdag_remark_t remark;
};

// Function declarations
int account_callback(void *data, xdag_hash_t hash, xdag_amount_t amount, xdag_time_t time, int n_our_key);
int xfer_callback(void *data, xdag_hash_t hash, xdag_amount_t amount, xdag_time_t time, int n_our_key);

int account_callback(void *data, xdag_hash_t hash, xdag_amount_t amount, xdag_time_t time, int n_our_key)
{
	char address[33] = {0};
	struct account_callback_data *d = (struct account_callback_data *)data;
	if(!d->count--) {
		return -1;
	}
	xdag_hash2address(hash, address);

	if(xdag_get_state() < XDAG_STATE_XFER) {
		sprintf(d->out, "%s  key %d", address, n_our_key);
	} else {
		sprintf(d->out, "%s %20.9Lf  key %d", address, amount2xdags(amount), n_our_key);
	}
	return 0;
}

xdag_error_no processAccountCommand(char **out)
{
	struct account_callback_data d;
	d.count = 1;

	char tmp[128] = {0};
	if(xdag_get_state() < XDAG_STATE_XFER) {
		sprintf(tmp, "Not ready to show balances. Type 'state' command to see the reason.");
	}
	xdag_traverse_our_blocks(&d, &account_callback);

	*out = strdup(strcat(tmp, d.out));

	return error_none;
}

xdag_error_no processAddressCommand(char **out)
{
	struct account_callback_data d;
	d.count = 1;

	char tmp[33] = {0};

	xdag_hash_t hash;
	xdag_get_our_block(hash);
	xdag_hash2address(hash, tmp);

	*out = strdup(tmp);

	return error_none;
}

xdag_error_no processBalanceCommand(char **out)
{
	if(xdag_get_state() < XDAG_STATE_XFER) {
		*out = strdup("Not ready to show a balance. Type 'state' command to see the reason.");
		return error_not_ready;
	} else {
		xdag_amount_t balance;
		balance = xdag_get_balance(0);
		char result[128] = {0};
		sprintf(result, "%.9Lf", amount2xdags(balance));
		*out = strdup(result);

		return error_none;
	}
}

xdag_error_no processLevelCommand(const char *level, char **out)
{
	unsigned lv;
	if(!level) {
		char tmp[16];
		sprintf(tmp, "%d\n", xdag_set_log_level(-1));
		*out = strdup(tmp);
		return error_none;
	} else if(sscanf(level, "%u", &lv) != 1 || lv > XDAG_TRACE) {
		*out = strdup("Illegal level.\n");
		return error_incorrect_level;
	} else {
		xdag_set_log_level(lv);
		return error_none;
	}
}

xdag_error_no processXferCommand(const char *amount, const char *address, const char *remark, char **out)
{
	if(!amount) {
		*out = strdup("Xfer: amount not given.");
		return error_xfer_no_amount;
	}
	if(!address) {
		*out = strdup("Xfer: destination address not given.");
		return error_xfer_no_address;
	}
	if (!remark) {
		*out = strdup("Xfer: remark not given.");
		return error_xfer_no_remark;
	}
	if(xdag_user_crypt_action(0, 0, 0, 3)) {
		sleep(3);
		*out = strdup("Password incorrect.");
		return error_pwd_incorrect;
	} else {
		return xdag_do_xfer(amount, address, remark, out);
	}
}


xdag_error_no processStateCommand(char **out)
{
	*out = strdup(xdag_get_state_str());
	return error_none;
}

static int make_transaction_block(struct xfer_callback_data *xferData)
{
	char address[33];
	if(xferData->fieldsCount != XFER_MAX_IN) {
		memcpy(xferData->fields + xferData->fieldsCount, xferData->fields + XFER_MAX_IN, sizeof(xdag_hashlow_t));
	}
	xferData->fields[xferData->fieldsCount].amount = xferData->todo;

	if(xferData->hasRemark) {
		memcpy(xferData->fields + xferData->fieldsCount + xferData->hasRemark, xferData->remark, sizeof(xdag_remark_t));
	}

	int res = xdag_create_block(xferData->fields, xferData->fieldsCount, 1, xferData->hasRemark, 0, 0, xferData->transactionBlockHash);
	if(res) {
		xdag_hash2address(xferData->fields[xferData->fieldsCount].hash, address);
		xdag_err(error_block_create, "FAILED: to %s xfer %.9Lf %s, error %d",
			address, amount2xdags(xferData->todo), COINNAME, res);
		return error_block_create;
	}
	xferData->done += xferData->todo;
	xferData->todo = 0;
	xferData->fieldsCount = 0;
	xferData->keysCount = 0;
	xferData->outsig = 1;
	return 0;
}

xdag_error_no xdag_do_xfer(const char *amount, const char *address, const char *remark, char **out)
{
	char address_buf[33];
	char result[256] = {0};
	struct xfer_callback_data xfer;

	memset(&xfer, 0, sizeof(xfer));
	xfer.remains = xdags2amount(amount);
	if(!xfer.remains) {
		*out = strdup("Xfer: nothing to transfer.");
		return error_xfer_nothing;
	}

	if(xfer.remains > xdag_get_balance(0)) {
		*out = strdup("Xfer: balance too small.");
		return error_xfer_too_small;
	}

	if(xdag_address2hash(address, xfer.fields[XFER_MAX_IN].hash)) {
		*out = strdup("Xfer: incorrect address.");
		return error_xfer_incorrect_address;
	}
#if REMARK_ENABLED
	if (remark) {
		if (!validate_remark(remark)) {
			if (out) {
				fprintf(out, "Xfer: transaction remark exceeds max length 32 chars or is invalid ascii.\n");
			}
			return error_xfer_incorrect_remark;
		}
		else {
			memcpy(xfer.remark, remark, strlen(remark));
			xfer.hasRemark = 1;
		}
	}
#endif

	xdag_wallet_default_key(&xfer.keys[XFER_MAX_IN]);
	xfer.outsig = 1;
	xdag_set_state(XDAG_STATE_XFER);
	g_xdag_xfer_last = time(0);

	int err = xdag_traverse_our_blocks(&xfer, &xfer_callback);
	if(err != 0 && err != 1) {
		sprintf(result, "%.9Lf", amount2xdags(xfer.done));
		return (xdag_error_no)err;
	}

	xdag_hash2address(xfer.transactionBlockHash, address_buf);
	sprintf(result, "%s", address_buf);
	*out = strdup(result);

	return error_none;
}

int xfer_callback(void *data, xdag_hash_t hash, xdag_amount_t amount, xdag_time_t time, int n_our_key)
{
	struct xfer_callback_data *xferData = (struct xfer_callback_data*)data;
	xdag_amount_t todo = xferData->remains;
	int i;
	if(!amount) {
		return error_xfer_nothing;
	}
	for(i = 0; i < xferData->keysCount; ++i) {
		if(n_our_key == xferData->keys[i]) {
			break;
		}
	}
	if(i == xferData->keysCount) {
		xferData->keys[xferData->keysCount++] = n_our_key;
	}
	if(xferData->keys[XFER_MAX_IN] == n_our_key) {
		xferData->outsig = 0;
	}
	if(Nfields(xferData) > XDAG_BLOCK_FIELDS) {
		int err = make_transaction_block(xferData);
		if(err) {
			return err;
		}
		xferData->keys[xferData->keysCount++] = n_our_key;
		if(xferData->keys[XFER_MAX_IN] == n_our_key) {
			xferData->outsig = 0;
		}
	}
	if(amount < todo) {
		todo = amount;
	}
	memcpy(xferData->fields + xferData->fieldsCount, hash, sizeof(xdag_hashlow_t));
	xferData->fields[xferData->fieldsCount++].amount = todo;
	xferData->todo += todo;
	xferData->remains -= todo;
	xdag_log_xfer(hash, xferData->fields[XFER_MAX_IN].hash, todo);
	if(!xferData->remains || Nfields(xferData) == XDAG_BLOCK_FIELDS) {
		int err = make_transaction_block(xferData);
		if(err) {
			return err;
		}
		if(!xferData->remains) { // xfer done
			return 1;
		}
	}
	return 0;
}

void xdag_log_xfer(xdag_hash_t from, xdag_hash_t to, xdag_amount_t amount)
{
	char address_from[33], address_to[33];
	xdag_hash2address(from, address_from);
	xdag_hash2address(to, address_to);
	xdag_mess("Xfer : from %s to %s xfer %.9Lf %s", address_from, address_to, amount2xdags(amount), COINNAME);
}

xdag_error_no processExitCommand()
{
	xdag_wallet_finish();
	xdag_storage_finish();
	xdag_blocks_finish();

	return (xdag_error_no)-1;
}

xdag_error_no processHelpCommand(char **out)
{
	*out = strdup("Commands:\n"
		"  account             - print our address with their amounts\n"
		"  address             - print our address\n"
		"  balance             - print balance of the address A or total balance for all our addresses\n"
		"  level [N]           - print level of logging or set it to N (0 - nothing, ..., 9 - all)\n"
		"  state               - print the program state\n"
		"  xfer S A            - transfer S our xdag to the address A\n"
		"  exit                - exit this program\n"
		"  help                - print this help\n");

	return error_none;
}

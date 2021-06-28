//
//  terminal.c
//  xDagWallet
//
//  Created by Rui Xie on 7/23/18.
//  Copyright © 2018 xrdavies. All rights reserved.
//

#include "terminal.h"
#include "commands.h"
#include "common.h"
#include "xdag_wrapper.h"
#include <string.h>
#include <stdlib.h>

typedef int (*xdag_com_func_t)(char *);
typedef struct {
	char *name;				/* command name */
	xdag_com_func_t func;	/* command function */
} XDAG_COMMAND;


int read_command(char *cmd);
int xdag_command(char *cmd);

int xdag_com_account(char *);
int xdag_com_address(char *);
int xdag_com_balance(char *);
int xdag_com_level(char *);
int xdag_com_xfer(char *);
int xdag_com_state(char *);
int xdag_com_help(char *);
int xdag_com_exit(char *);

XDAG_COMMAND commands[] = {
	{ "account"    , xdag_com_account },
	{ "address"    , xdag_com_address },
	{ "balance"    , xdag_com_balance },
	{ "level"      , xdag_com_level },
	{ "xfer"       , xdag_com_xfer },
	{ "state"      , xdag_com_state },
	{ "exit"       , xdag_com_exit },
	{ "xfer"       , (xdag_com_func_t)NULL},
	{ "help"       , xdag_com_help},
	{ (char *)NULL , (xdag_com_func_t)NULL}
};

int xdag_com_account(char *args)
{
	return xdag_wrapper_account();
}

int xdag_com_address(char *args)
{
	return xdag_wrapper_address();
}

int xdag_com_balance(char *args)
{
	return xdag_wrapper_balance();
}

int xdag_com_xfer(char *args)
{
	char *amount = strtok_r(args, " \t\r\n", &args);
	char *address = strtok_r(0, " \t\r\n", &args);
	char *remark = strtok_r(0, " \t\r\n", &args);

	return xdag_wrapper_xfer(amount, address, remark);
}

int xdag_com_level(char *args)
{
	char *cmd = strtok_r(args, " \t\r\n", &args);
	return xdag_wrapper_level(cmd);
}

int xdag_com_state(char *args)
{
	return xdag_wrapper_state();
}

int xdag_com_exit(char * args)
{
	return xdag_wrapper_exit();
}

int xdag_com_help(char *args)
{
	return xdag_wrapper_help();
}

int xdag_command(char *cmd)
{
	char *nextParam;

	cmd = strtok_r(cmd, " \t\r\n", &nextParam);
	if(!cmd) return 0;

    for(int i = 0; commands[i].name; i++) {
        if(strcmp(cmd, commands[i].name) == 0) {
            return (*(commands[i].func))(nextParam);
        }
    }

    xdag_wrapper_event(event_id_promot, error_none, "Illegal command.\n");
	return 0;
}

int read_command(char *cmd)
{
	printf("xdag> ");
	fflush(stdout);
	fgets(cmd, XDAG_COMMAND_MAX, stdin);
	return 0;
}

void startCommandProcessing(void)
{
	char cmd[XDAG_COMMAND_MAX];
	xdag_wrapper_event(event_id_promot, error_none, "Type command, help for example.\n");

	for(;;) {
		read_command(cmd);
		if(strlen(cmd) > 0) {
			int ret = xdag_command(cmd);
			if(ret < 0) {
				break;
			}
		}
	}
}


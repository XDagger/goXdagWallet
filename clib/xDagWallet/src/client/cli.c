#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <ctype.h>
#include "version.h"
#include "cli.h"
#include "common.h"
#include "client.h"
#include "terminal.h"
#include "utils/utils.h"
#include "utils/log.h"
#include "xdag_wrapper.h"

#if defined(_WIN32) || defined(_WIN64)
#else
#include <sys/termios.h>
#endif

/*
 Command line wallet
 */

void printUsage(char* appName);
int log_callback(int level, xdag_error_no err, char *buffer);
int event_callback(void* thisObj, xdag_event *event);
int password_callback(const char *prompt, char *buf, unsigned len);

int input(xdag_event_id evt, xdag_wrapper_msg *msg);

static int g_client_init_done = 0;

int xdag_cli_init(int argc, char **argv)
{
    printf("XDAG client/server, version %s.\n", XDAG_VERSION);
    
    if (argc <= 1) {
        printUsage(argv[0]);
        return 0;
    }
    
    const char *pool_arg = 0;
    int testnet = 0;
    
	for (int i = 1; i < argc; ++i) {
		if (argv[i][0] != '-') {
			if ((!argv[i][1] || argv[i][2]) && strchr(argv[i], ':')) {
				pool_arg = argv[i];
			} else {
				printUsage(argv[0]);
				return 0;
			}
			continue;
		}
        
        if(ARG_EQUAL(argv[i], "-t", "--testnet")) { /* testnet */
            testnet = 1;
        } else if(ARG_EQUAL(argv[i], "-h", "--help")) { /* help */
			printUsage(argv[0]);
			return 0;
		} else {
			printUsage(argv[0]);
			return 0;
		}
	}
    
    printf("Init path...\n");
    xdag_init_path(argv[0]);
    printf("Set log callback...\n");
    xdag_wrapper_init(NULL, &password_callback, &event_callback);
    
	printf("Starting command line wallet...\n");
    xdag_thread_param_t param;
    strncpy(param.pool_arg, pool_arg, 255);
    param.testnet = testnet;
    
    int err = pthread_create(&g_client_thread, 0, xdag_client_thread, (void*)&param);
    if(err != 0) {
        printf("create xdag_client_thread failed, error : %s\n", strerror(err));
        return -1;
    }
    
    while (!g_client_init_done) {
        sleep(1);
    }
    startCommandProcessing();
    
	return 0;
}

#define XDAG_LOG_FILE "xdag.log"

int event_callback(void* thisObj, xdag_event *event)
{
	if(!event) {
		return -1;
	}

	switch (event->event_id) {
        case event_id_init_done:
        {
            g_client_init_done = 1;
            break;
        }
            
        case event_id_promot:
        {
            printf("%s\n", event->event_data);
            break;
        }
            
		case event_id_log:
		{
			FILE *f;
			char buf[64] = {0};
			sprintf(buf, XDAG_LOG_FILE);
			f = xdag_open_file(buf, "a");
			if (f) {
				fprintf(f, "%s\n", event->event_data);
			}

			xdag_close_file(f);
			break;
		}

		case event_id_interact:
		{
			fprintf(stdout, "%s\n", event->event_data);
			fflush(stdout);
			break;
		}

        case event_id_err:
        {
            fprintf(stdout, "error : %x, msg : %s\n", event->error_no, event->event_data);
            fflush(stdout);
            break;
        }

		case event_id_err_exit:
		{
			fprintf(stdout, "error : %x, msg : %s\n", event->error_no, event->event_data);
			fflush(stdout);
			xdag_wrapper_exit();
			pthread_cancel(g_client_thread);
			exit(1);
			break;
		}
            
        case event_id_exit:
        {
            xdag_wrapper_exit();
            pthread_cancel(g_client_thread);
            exit(1);
            break;
        }

		case event_id_account_done:
		{
			fprintf(stdout, "%s\n", event->event_data);
			fflush(stdout);
			break;
		}

		case event_id_address_done:
		{
			fprintf(stdout, "%s\n", event->event_data);
			fflush(stdout);
			break;
		}

		case event_id_balance_done:
		{
			fprintf(stdout, "%s\n", event->event_data);
			fflush(stdout);
			break;
		}

		case event_id_xfer_done:
		{
			fprintf(stdout, "%s\n", event->event_data);
			fflush(stdout);
			break;
		}

		case event_id_level_done:
		{
			fprintf(stdout, "%s\n", event->event_data);
			fflush(stdout);
			break;
		}

		case event_id_state_done:
		{
			fprintf(stdout, "%s\n", event->event_data);
			fflush(stdout);
			break;
		}

		case event_id_exit_done:
		{
			fprintf(stdout, "%s\n", event->event_data);
			fflush(stdout);
			break;
		}

        case event_id_passwd:
        case event_id_set_passwd:
        case event_id_set_passwd_again:
        case event_id_random_key:
        {
            input(event->event_id, (xdag_wrapper_msg *)event->event_data);
            break;
        }
            
		case event_id_state_change:
		{
			fprintf(stdout, "State changed %s\n", event->event_data);
			fflush(stdout);
			break;
		}

		default:
		{
			break;
		}
	}
	return 0;
}



int password_callback(const char *prompt, char *buf, unsigned len)
{
    /*
     Password
     Set password
     Re-type password
     Type random keys
     */
#if !defined(_WIN32) && !defined(_WIN64)
	struct termios t[1];
	int noecho = !!strstr(prompt, "assword");
	fprintf(stdout,"%s: ", prompt); fflush(stdout);

	if (noecho) {
		tcgetattr(0, t);
		t->c_lflag &= ~ECHO;
		tcsetattr(0, TCSANOW, t);
	}
	fgets(buf, len, stdin);
	if (noecho) {
		t->c_lflag |= ECHO;
		tcsetattr(0, TCSANOW, t);
		fprintf(stdout,"\n");
	}
	len = (int)strlen(buf);
	if (len && buf[len - 1] == '\n') buf[len - 1] = 0;
#endif

	return 0;
}

int input(xdag_event_id evt, xdag_wrapper_msg *msg)
{
#if !defined(_WIN32) && !defined(_WIN64)
    /*
     Password
     Set password
     Re-type password
     Type random keys
     */
    struct termios t[1];
    int noecho = 0;
    if (evt == event_id_passwd) {
        noecho = 1;
        fprintf(stdout,"%s: ", "Password"); fflush(stdout);
    } else if(evt == event_id_set_passwd) {
        noecho = 1;
        fprintf(stdout,"%s: ", "Set password"); fflush(stdout);
    } else if(evt == event_id_set_passwd_again) {
        noecho = 1;
        fprintf(stdout,"%s: ", "Re-type password"); fflush(stdout);
    } else if(evt == event_id_random_key) {
        noecho = 0;
        fprintf(stdout,"%s: ", "Type random keys"); fflush(stdout);
    }
    
    if (noecho) {
        tcgetattr(0, t);
        t->c_lflag &= ~ECHO;
        tcsetattr(0, TCSANOW, t);
    }
    
    fgets(msg->msg, XDAG_WRAPPER_MSG_LEN, stdin);
    
    if (noecho) {
        t->c_lflag |= ECHO;
        tcsetattr(0, TCSANOW, t);
        fprintf(stdout,"\n");
    }
    int len = (int)strlen(msg->msg);
    if (len && msg->msg[len - 1] == '\n') msg->msg[len - 1] = 0;
#endif
    
    return 0;
}

void printUsage(char* appName)
{
	printf("Usage: %s flags [pool_ip:port]\n"
		"If pool_ip:port argument is given, then the node operates as a miner.\n"
		"Flags:\n"
		"  -h             - print this help\n"
		"  -t             - connect to test net (default is main net)\n"
		"  -v N           - set loglevel to N\n"
		, appName);
}

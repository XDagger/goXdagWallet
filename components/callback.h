#include "../xDagWallet/src/client/events.h"

void init_event_callback();
extern int goEventCallback(void*, xdag_event *);

void init_password_callback();
extern int goPasswordCallback(const char *prompt, char *buf, unsigned size);
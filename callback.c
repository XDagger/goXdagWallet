#include "callback.h"
#include "./src/xdag_runtime.h"
void init_event_callback() {
    xdag_set_event_callback_wrap(goEventCallback);
}

void init_password_callback() {
    xdag_set_password_callback_wrap(goPasswordCallback);
}
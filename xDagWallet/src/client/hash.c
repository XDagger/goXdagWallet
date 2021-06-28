#include <string.h>
#ifdef SHA256_OPENSSL_MBLOCK
#endif
#include "sha256.h"
#include "hash.h"
#include "system.h"

#if defined(_WIN32) || defined(_WIN64)
//// #include <WinSock2.h>
#else
#include <arpa/inet.h>
#endif


void xdag_hash(void *data, size_t size, xdag_hash_t hash)
{
	SHA256REF_CTX ctx;

	sha256_init(&ctx);
	sha256_update(&ctx, data, size);
	sha256_final(&ctx, (uint8_t*)hash);
	sha256_init(&ctx);
	sha256_update(&ctx, (uint8_t*)hash, sizeof(xdag_hash_t));
	sha256_final(&ctx, (uint8_t*)hash);
}

unsigned xdag_hash_ctx_size(void)
{
	return sizeof(SHA256REF_CTX);
}

void xdag_hash_init(void *ctxv)
{
	SHA256REF_CTX *ctx = (SHA256REF_CTX*)ctxv;

	sha256_init(ctx);
}

void xdag_hash_update(void *ctxv, void *data, size_t size)
{
	SHA256REF_CTX *ctx = (SHA256REF_CTX*)ctxv;

	sha256_update(ctx, data, size);
}

void xdag_hash_final(void *ctxv, void *data, size_t size, xdag_hash_t hash)
{
	SHA256REF_CTX ctx;

	memcpy(&ctx, ctxv, sizeof(ctx));
	sha256_update(&ctx, (uint8_t*)data, size);
	sha256_final(&ctx, (uint8_t*)hash);
	sha256_init(&ctx);
	sha256_update(&ctx, (uint8_t*)hash, sizeof(xdag_hash_t));
	sha256_final(&ctx, (uint8_t*)hash);
}

uint64_t xdag_hash_final_multi(void *ctxv, uint64_t *nonce, int attempts, int step, xdag_hash_t hash)
{
	SHA256REF_CTX ctx;
	xdag_hash_t hash0;
	uint64_t min_nonce = 0;
	int i;

	for (i = 0; i < attempts; ++i) {
		memcpy(&ctx, ctxv, sizeof(ctx));
		sha256_update(&ctx, (uint8_t*)nonce, sizeof(uint64_t));
		sha256_final(&ctx, (uint8_t*)hash0);
		sha256_init(&ctx);
		sha256_update(&ctx, (uint8_t*)hash0, sizeof(xdag_hash_t));
		sha256_final(&ctx, (uint8_t*)hash0);

		if (!i || xdag_cmphash(hash0, hash) < 0) {
			memcpy(hash, hash0, sizeof(xdag_hash_t));
			min_nonce = *nonce;
		}

		*nonce += step;
	}

	return min_nonce;
}

void xdag_hash_get_state(void *ctxv, xdag_hash_t state)
{
	SHA256REF_CTX *ctx = (SHA256REF_CTX*)ctxv;

	memcpy(state, ctx->state, sizeof(xdag_hash_t));
}

void xdag_hash_set_state(void *ctxv, xdag_hash_t state, size_t size)
{
	SHA256REF_CTX *ctx = (SHA256REF_CTX*)ctxv;

	memcpy(ctx->state, state, sizeof(xdag_hash_t));
	ctx->datalen = 0;
	ctx->bitlen = (unsigned int) (size << 3);
	ctx->bitlenH = 0;
	ctx->md_len = SHA256_BLOCK_SIZE;
}

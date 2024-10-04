CREATE OR REPLACE FUNCTION evm_hex_encode_int AS (i) -> if(
    i = 0,
    '0x0',
    concatAssumeInjective('0x', trim(LEADING '0' FROM lower(hex(i))))
);
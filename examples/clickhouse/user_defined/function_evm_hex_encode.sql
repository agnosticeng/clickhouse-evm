CREATE FUNCTION evm_hex_encode AS s -> concatAssumeInjective('0x', lower(hex(s)))

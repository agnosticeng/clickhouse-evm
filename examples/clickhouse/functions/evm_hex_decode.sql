CREATE FUNCTION evm_hex_decode AS (s) -> unhex(substring(s, 3));
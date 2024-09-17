CREATE FUNCTION evm_hex_decode_uint256 AS s -> reinterpretAsUInt256(reverse(evm_hex_decode(replaceRegexpAll(s, concat('^[', regexpQuoteMeta('"'), ']+|[', regexpQuoteMeta('"'), ']+$'), ''))))

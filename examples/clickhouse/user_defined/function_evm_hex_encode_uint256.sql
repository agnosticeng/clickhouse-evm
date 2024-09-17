CREATE FUNCTION evm_hex_encode_uint256 AS i -> concatAssumeInjective('0x', replaceRegexpOne(lower(hex(i)), concat('^[', regexpQuoteMeta('0'), ']+'), ''))

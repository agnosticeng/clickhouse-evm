-- decode_string

select throwIf(evm_hex_decode('0x68656c6c6f') != 'hello')

;;

-- encode_string

select throwIf(evm_hex_encode('hello') != '0x68656c6c6f')

;;

-- decode_int

select throwIf(evm_hex_decode_int('0x7b', 'Int64') != 123)

;;

-- encode_int

select throwIf(evm_hex_encode_int(123) != '0x7b')
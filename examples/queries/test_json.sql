select 
    JSONExtract('{"key": "value"}', 'Tuple(String)'),
    JSON_VALUE('{"key": "value"}', '$.key')


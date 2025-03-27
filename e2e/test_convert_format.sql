-- yaml_to_json

select 
    throwIf(JSONExtract(v, 'root', 'key1', 'String') != 'val'),
    throwIf(JSONExtract(v, 'root', 'arr', 'Array(String)') != ['item1', 'item2'])
from (
    select convert_format('YAML', 'JSON', $heredoc$
        root:
            key1: val
            arr:
                - item1
                - item2
    $heredoc$) as v
)

;;

-- toml_to_json

select 
    throwIf(JSONExtract(v, 'root', 'key1', 'String') != 'val'),
    throwIf(JSONExtract(v, 'root', 'arr', 'Array(String)') != ['item1', 'item2'])
from (
    select convert_format('TOML', 'JSON', $heredoc$
        [root]
        key1 = "val"
        arr = ["item1", "item2"]
    $heredoc$) as v
)
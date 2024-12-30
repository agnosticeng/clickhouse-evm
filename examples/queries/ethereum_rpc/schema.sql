select
    arrayJoin(
        distinctJSONPathsAndTypes(
            JSONExtract(
                ethereum_rpc(
                    'trace_block',
                    [evm_hex_encode_int(number)],
                    'https://nood-beta.eu-02.flubeer.xyz/polygon/mainnet?token=eyJhbGciOiJFUzI1NiIsImtpZCI6InBfQ3hFZnBzODNrV3FlMnZUM2ZJU3Q5MmxjcS1odmJuTEIxampsV0w4elUiLCJ0eXAiOiJKV1QifQ.eyJhdWQiOiJpc3RpbyIsImV4cCI6MTczNTY4NjAwMCwiaWF0IjoxNzA0MjkxMzUwLCJpc3MiOiJub29kQGFnbm9zdGljLWluZnJhIiwianRpIjoiZDU2OWFlNzgxYWIxNTU2YzAxOGU3NDBmODU1NTE5NGVjYjNlYzlhYjMyMTRiMzE0MDBlYzJkNzlhOWUyYWQ5NCIsIm5iZiI6MTcwNDI5MTM1MCwic3ViIjoibm9vZCJ9.2A-UrtIyslQZgIQ8F6P6fDjrAqXYC1ppb5BsQXG7cC9Wv2HM8y5ZBqfZTgOs97Zio9dwBPszxJOqF3VlvjlwPQ'
                ),
                'value', 1,
                'JSON'
            )
        )
    )
from numbers(65351185, 10)




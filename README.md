## Using agnostic-clickhouse-udf with clickhouse-local

### Point clickhouse-local the the custom config file

```sh
clickhouse local --config-file ./examples/clickhouse/config.xml
```

### Running clickhouse-local with the custom config

```sh
clickhouse local --config examples/clickhouse/config.xml --path tmp/clickhouse
```

## Running chdb in python 

```sh
python3 -m venv venv
source venv/bin/activate
pip install -t requirements.txt
make
cp ./bin/agnostic-clickhouse-udf ./examples/clickhouse/user_defined
python examples/chdb/run_query.py ./examples/queries/num_transfer_event.sql
```
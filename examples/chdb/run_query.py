import os
import sys
from chdb import session as chs

default_udf_path = './examples/clickhouse/user_defined'

if __name__ == '__main__':
    q = open(sys.argv[1]).read()
    sess = chs.Session()

    for f in os.listdir(default_udf_path):
        if f.endswith(".sql"):
            content = open(os.path.join(default_udf_path, f)).read()
            sess.query(content)

    res = sess.query(q, 'Pretty', udf_path=default_udf_path)

    print(res)

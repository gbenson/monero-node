import os
import sys

from .application import Application


def main(redis_url=None, redis_url_env_var="RIGSTATUS_DB_URL"):
    if redis_url is None:
        redis_url = os.environ.get(redis_url_env_var)
    if not redis_url:
        print(f"rs: please export {redis_url_env_var}", file=sys.stderr)
        return 1
    Application(redis_url).run()

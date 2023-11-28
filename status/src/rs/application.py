import json
import os
import sys
import time

from datetime import datetime, timedelta, timezone
from functools import cached_property

from redis import Redis


class Application:
    def __init__(self, redis_url):
        self.redis_url = redis_url
        self.window = timedelta(seconds=65)
        self.poll_interval = timedelta(minutes=1)
        self._report_cache = {}

    @cached_property
    def db(self):
        return Redis.from_url(self.redis_url, decode_responses=True)

    def run(self):
        while True:
            self._step()

    def _step(self):
        step_start = datetime.now(timezone.utc)
        window_start = (step_start - self.window).timestamp()
        reports = self.get_report_range(window_start, "+inf")
        if reports:
            self.log_window(reports)
            step_start = reports[0].time
        else:
            self.log("[no miners]")
        sleep_until = step_start + self.poll_interval
        sleep_time = sleep_until - datetime.now(timezone.utc)
        sleep_secs = sleep_time.total_seconds()
        if sleep_time < self.poll_interval:
            sleep_secs += 1  # don't expect exact report timings?
        if sleep_secs > 0:
            time.sleep(sleep_secs)

    def get_report_range(self, start, limit):
        return [self.get_report(reporter, unixtime)
                for reporter, unixtime in self.db.zrangebyscore(
                        "reports", start, limit, withscores=True)]

    def get_report(self, reporter, unixtime):
        report = self._report_cache.get(reporter)
        if report is not None:
            if abs(unixtime - report.unixtime) < 1e-2:
                return report
        report = Report(
            reporter,
            unixtime,
            json.loads(self.db.get(reporter)),
        )
        self._report_cache[reporter] = report
        return report

    def log_window(self, reports):
        cutoff = datetime.now(timezone.utc) - timedelta(seconds=61)
        hashrate = sum(report.hashrate_10s
                       for report in reports
                       if report.time > cutoff)
        hashrate = f"{hashrate:5.0f} H/s"
        for report in reports:
            if report.is_logged:
                continue
            self.log(f"{report.summary} {hashrate}", time=report.time)
            report.is_logged = True

    def log(self, msg, *, time=None):
        if time is None:
            time = datetime.now(timezone.utc)
        print(str(time)[:23], msg)


class Report:
    def __init__(self, reporter_key, unixtime, report):
        self.key = reporter_key
        self.unixtime = unixtime
        self.report = report
        self.is_logged = False

    @cached_property
    def time(self):
        return datetime.fromtimestamp(self.unixtime, timezone.utc)

    @property
    def summary(self):
        return f"{self.worker_id:>34s}  {self.hashrate_summary}"

    @property
    def hashrate_summary(self):
        return f"{self.hashrate_10s:4.0f} H/s"

    @property
    def worker_id(self):
        return self.report["worker_id"]

    @property
    def hashrate_10s(self):
        return self.hashrate["total"][0]

    @property
    def hashrate_max(self):
        return self.hashrate["highest"]

    @property
    def hashrate(self):
        return self.miner_status["hashrate"]

    @property
    def miner_status(self):
        return self.report["miner_status"]


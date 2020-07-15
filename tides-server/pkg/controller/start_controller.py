from crontab import CronTab
import os
from config import BASE_DIR


def start_controller():
    my_cron = CronTab(user='root')
    polling_interval = 1
    for job in my_cron:
        if job.comment == 'monitor_usage' or job.comment == 'idle_monitor':
            my_cron.remove(job)
            my_cron.write()
    path = os.path.join(BASE_DIR, 'controller')
    job1 = my_cron.new(command='python ' + path + '/monitor.py', comment='monitor_usage')
    job1.minute.every(polling_interval)

    job2 = my_cron.new(command='python ' + path + '/idle_monitor.py', comment='idle_monitor')
    job2.minute.every(polling_interval)
    my_cron.write()

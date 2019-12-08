import psycopg2
import os
from core.settings import DATABASES, BASE_DIR
from crontab import CronTab


def monitor():
    db = DATABASES['default']['NAME']
    user = DATABASES['default']['USER']
    password = DATABASES['default']['PASSWORD']
    host = DATABASES['default']['HOST']
    port = DATABASES['default']['PORT']
    conn = psycopg2.connect(database=db, user=user, password=password, host=host, port=port)
    cur = conn.cursor()

    cur.execute(
        'SELECT host_address, host_name, username, password, polling_interval FROM resource_resource WHERE monitored '
        '= False')
    results = cur.fetchall()
    path = os.path.join(BASE_DIR, 'controller')
    my_cron = CronTab(user='dbaekajnr')
    # for job in my_cron
    #     if job.comment == 'dateinfo':
    #         my_cron.remove(job)
    #         my_cron.write()
    for result in results:
        command1 = 'python ' + path + '/query_usage.py -s ' + result[0] + ' -u ' + result[
            2] + ' -p ' + \
                   result[3] + ' -n ' + result[1] + ' --no-ssl\n'
        job1 = my_cron.new(command=command1, comment='usage-' + result[1])
        job1.minute.every(result[4])

        command2 = 'python ' + path + '/get_vm_usage_class.py -s ' + result[0] + ' -u ' + \
                   result[2] + ' -p ' + \
                   result[3] + ' --no-ssl\n'
        job2 = my_cron.new(command=command2, comment='vm_usage-' + result[1])
        job2.minute.every(result[4])

        cur.execute('UPDATE resource_resource SET monitored = True WHERE host_address = %s AND host_name = %s',
                    (result[0], result[1]))
        conn.commit()

    conn.commit()
    cur.close()
    conn.close()
    my_cron.write()

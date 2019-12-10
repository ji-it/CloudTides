import psycopg2
import os
import json
import ast
import requests

FULL_HOSTNAME = "http://localhost:8000"
DATABASES = {
    'default': {
        'NAME': 'tides2',
        'USER': 'postgres',
        'PASSWORD': 'password',  # created at the time of password setup
        'HOST': 'localhost',
        'PORT': '5432',
    }
}
BASE_DIR = "/Users/dbaekajnr/Projects/CloudTides/server"


def main():
    db = DATABASES['default']['NAME']
    user = DATABASES['default']['USER']
    password = DATABASES['default']['PASSWORD']
    host = DATABASES['default']['HOST']
    port = DATABASES['default']['PORT']
    conn = psycopg2.connect(database=db, user=user, password=password, host=host, port=port)
    cur = conn.cursor()
    path = os.path.join(BASE_DIR, 'controller')

    cur.execute(
        'SELECT id, host_address, username, password, policy_id, total_cpu, total_ram FROM resource_resource WHERE '
        'monitored = True')
    results = cur.fetchall()

    for result in results:
        resource_id = result[0]
        total_cpu, total_ram = result[5], result[6]
        cur.execute(
            'SELECT cpu, ram FROM usage_hostusage WHERE resource_id = ' + str(
                resource_id) + ' ORDER BY date_added DESC LIMIT 1')
        usage = cur.fetchone()

        cur.execute('SELECT idle_policy, threshold_policy, template_id, is_destroy FROM policy_policy WHERE id = %s',
                    str(result[4]))
        policy = cur.fetchone()
        deploy = False
        destroy = False
        idle_policy = eval(policy[0])  # idle policy
        busy_policy = eval(policy[1])  # busy policy

        cpu_usage, mem_usage = usage[0] / total_cpu, usage[1] / total_ram
        if 'cpu' not in idle_policy.keys():
            if mem_usage < idle_policy['ram']:
                deploy = True
        elif 'ram' not in idle_policy.keys():
            if cpu_usage < idle_policy['cpu']:
                deploy = True
        else:
            if cpu_usage < idle_policy['cpu'] and mem_usage < idle_policy['ram']:
                deploy = True

        if 'cpu' not in busy_policy.keys():
            if mem_usage > busy_policy['ram']:
                destroy = True
        elif 'ram' not in busy_policy.keys():
            if cpu_usage > busy_policy['cpu']:
                destroy = True
        else:
            if cpu_usage > busy_policy['cpu'] or mem_usage > busy_policy['ram']:
                destroy = True

        if deploy:
            cur.execute("UPDATE resource_resource SET status = 'idle' WHERE id = " + str(resource_id))
            conn.commit()
            cur.execute("SELECT name FROM template_template WHERE id = " + str(policy[2]))
            template_name = cur.fetchone()
            print('python ' + path + '/clone_vm.py -s ' + result[1] + ' -u ' + result[2] + ' -p ' + result[3] + \
                  ' --no-ssl --power-on --template ' + template_name[0])
            os.system(
                'python ' + path + '/clone_vm.py -s ' + result[1] + ' -u ' + result[2] + ' -p ' + result[3] + \
                ' --no-ssl --power-on --template ' + template_name[0])

        elif destroy:
            cur.execute("UPDATE resource_resource SET status = 'busy' WHERE id = " + str(resource_id))
            conn.commit()
            if policy[3] is False:
                ## Place code to suspend VMs here
                pass
            else:
                cur.execute("SELECT ip_address FROM resource_vm WHERE is_destroyed = False AND resource_id = " + str(
                    resource_id))
                vms = cur.fetchall()
                for vm in vms:
                    data = {}
                    data['ip_address'] = vm[0]
                    requests.post(FULL_HOSTNAME + "/api/resource/destroy_vm/",
                                  data=json.dumps(data))  # TODO: make non-loop
                    os.system(
                        'python ' + path + '/destroy_vm.py -s ' + result[1] + ' -u ' + result[2] + ' -p ' +
                        result[3] + \
                        ' -i ' + vm[0])

        else:
            cur.execute("UPDATE resource_resource SET status = 'normal' WHERE id = " + str(resource_id))
            conn.commit()


# start
if __name__ == "__main__":
    main()

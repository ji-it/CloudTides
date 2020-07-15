import psycopg2
import os
import json
from config import BASE_DIR, DATABASES, FULL_HOSTNAME
import requests



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
        'SELECT id, host_address, username, password, policy_ref, name, is_resource_pool FROM resources WHERE '
        'monitored = True AND is_active = True')
    results = cur.fetchall()

    for result in results:
        resource_id = result[0]
        # total_cpu, total_ram = result[5], result[6]
        cur.execute(
            'SELECT percent_cpu, percent_ram FROM resource_usages WHERE resource_ref = ' + str(
                resource_id))
        usage = cur.fetchone()

        cur.execute('SELECT idle_policy, threshold_policy, template_ref, is_destroy FROM policies WHERE id = %s',
                    str(result[4]))
        policy = cur.fetchone()
        deploy = False
        destroy = False
        idle_policy = eval(policy[0])  # idle policy
        busy_policy = eval(policy[1])  # busy policy

        cpu_usage, mem_usage = usage[0], usage[1]
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
            # cur.execute("UPDATE resources SET status = 'idle' WHERE id = " + str(resource_id))
            # conn.commit()
            data = {}
            data['ResourceID'] = resource_id
            data['Status'] = 'idle'
            data['Monitored'] = True
            headers = {'Content-type': 'application/json'}
            requests.put(FULL_HOSTNAME + "/v1/resource/update_status/",
                data=json.dumps(data), headers=headers)
            cur.execute("SELECT name FROM templates WHERE id = " + str(policy[2]))
            template_name = cur.fetchone()
            if result[6]:
                os.system(
                    'python3 ' + path + '/clone_vm.py -s ' + result[1] + ' -u ' + result[2] + ' -p ' + result[3] + \
                    ' --no-ssl --power-on --template ' + template_name[0] + ' --resource-pool ' + result[5] + \
                    ' -n ' + result[5])
            else:
                os.system(
                    'python3 ' + path + '/clone_vm.py -s ' + result[1] + ' -u ' + result[2] + ' -p ' + result[3] + \
                    ' --no-ssl --power-on --template ' + template_name[0] + ' --cluster-name ' + result[5] + \
                    ' -n ' + result[5])

        elif destroy:
            # cur.execute("UPDATE resources SET status = 'busy' WHERE id = " + str(resource_id))
            # conn.commit()
            data = {}
            data['ResourceID'] = resource_id
            data['Status'] = 'busy'
            data['Monitored'] = True
            headers = {'Content-type': 'application/json'}
            requests.put(FULL_HOSTNAME + "/v1/resource/update_status/",
                data=json.dumps(data), headers=headers)
            if policy[3] is False:
                ## Place code to suspend VMs here
                pass
            else:
                cur.execute("SELECT ip_address FROM vms WHERE is_destroyed = False AND resource_ref = " + str(
                    resource_id))
                vms = cur.fetchall()
                for vm in vms:
                    data = {}
                    data['IPAddress'] = vm[0]
                    headers = {'Content-type': 'application/json'}
                    requests.delete(FULL_HOSTNAME + "/v1/resource/destroy_vm/",
                                  data=json.dumps(data), headers=headers)  # TODO: make non-loop
                    os.system(
                        'python3 ' + path + '/destroy_vm.py -s ' + result[1] + ' -u ' + result[2] + ' -p ' +
                        result[3] + \
                        ' -i ' + vm[0])

        else:
            # cur.execute("UPDATE resources SET status = 'normal' WHERE id = " + str(resource_id))
            # conn.commit()
            data = {}
            data['ResourceID'] = resource_id
            data['Status'] = 'normal'
            data['Monitored'] = True
            headers = {'Content-type': 'application/json'}
            requests.put(FULL_HOSTNAME + "/v1/resource/update_status/",
                data=json.dumps(data), headers=headers)


# start
if __name__ == "__main__":
    main()

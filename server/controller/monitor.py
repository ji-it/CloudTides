import psycopg2
import os
from config import BASE_DIR, DATABASES


def main():
    db = DATABASES['default']['NAME']
    user = DATABASES['default']['USER']
    password = DATABASES['default']['PASSWORD']
    host = DATABASES['default']['HOST']
    port = DATABASES['default']['PORT']
    conn = psycopg2.connect(database=db, user=user, password=password, host=host, port=port)
    cur = conn.cursor()

    cur.execute(
        'SELECT host_address, host_name, username, password FROM resource_resource')
    results = cur.fetchall()
    path = os.path.join(BASE_DIR, 'controller')

    for result in results:
        os.system('python ' + path + '/query_usage.py -s ' + result[0] + ' -u ' + result[
            2] + ' -p ' + \
                  result[3] + ' -n ' + result[1] + ' --no-ssl\n')

        os.system('python ' + path + '/get_vm_usage_class.py -s ' + result[0] + ' -u ' + \
                  result[2] + ' -p ' + \
                  result[3] + ' --no-ssl\n')

        cur.execute('UPDATE resource_resource SET monitored = True WHERE host_address = %s AND host_name = %s',
                    (result[0], result[1]))
        conn.commit()

    conn.commit()
    cur.close()
    conn.close()


# start
if __name__ == "__main__":
    main()

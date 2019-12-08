import psycopg2
import os

def main():
    conn=psycopg2.connect(database="Tides",user="postgres",
            password="t6bB2T5KoQuPq6DrpWxJa3rYKVjIpOCtVSrKyBMB8PHcMShkidcQo8Kjn1lcXswB",host="10.11.16.83",port="30123")
    cur=conn.cursor()

    cur.execute('SELECT host_address, host_name, username, password, polling_interval FROM resource_resource WHERE monitored = False')
    results = cur.fetchall()
    #print(results)
    with open('/home/shen1997/ve450/shencron', 'a') as f:
        for result in results:
            string = '*/' + str(result[4]) + ' * * * * python ~/ve450/query_usage.py -s ' + result[0] + ' -u ' + result[2] + ' -p ' +\
                    result[3] + ' -n ' + result[1] + ' --no-ssl\n'
            f.write(string)
            string2 = '*/' + str(result[4]) + ' * * * * python ~/ve450/get_vm_usage_class.py -s ' + result[0] + ' -u ' + result[2] + ' -p ' +\
                    result[3] + ' --no-ssl\n'
            f.write(string2)
            cur.execute('UPDATE resource_resource SET monitored = True WHERE host_address = %s AND host_name = %s', (result[0], result[1]))
            conn.commit()
    
    conn.commit()
    cur.close()
    conn.close()
    os.system('crontab /home/shen1997/ve450/shencron')

# start
if __name__ == "__main__":
    main()
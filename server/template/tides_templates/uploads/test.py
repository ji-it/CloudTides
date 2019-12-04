import psycopg2
 
#创建连接对象
conn=psycopg2.connect(database="Test",user="postgres",
    password="t6bB2T5KoQuPq6DrpWxJa3rYKVjIpOCtVSrKyBMB8PHcMShkidcQo8Kjn1lcXswB",host="10.11.16.83",port="30123")
cur=conn.cursor() #创建指针对象
 

 
# 获取结果
cur.execute('SELECT ip_address FROM usage_vmusage')
results=cur.fetchall()
print (results)

# 关闭练级
conn.commit()
cur.close()
conn.close()

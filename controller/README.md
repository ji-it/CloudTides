<<<<<<< HEAD
## ControllerCreate a file used for `crontab`. ```*/2 * * * * python ~/ve450/monitor.py``````*/3 * * * * python ~/ve450/idle_monitor.py```The `*`s from left to right represent minute, hour, day, month, year respectively. The above setting means `monitor.py` will execute every 2 minutes, and `idle_monitor.py` will execute every 3 minutes.```monitor.py``` searches any unmonitored resource and start to monitor them by setting `crontab` for `query_usage.py` and `get_vm_usage_class.py`.`idle_monitor.py` monitors whether the monitored resources are idle, normal or busy, and deploys or destroys VMs according to the current usage.
=======
## Controller

Create a file used for `crontab`. 

```
*/2 * * * * python ~/ve450/monitor.py
```

```
*/3 * * * * python ~/ve450/idle_monitor.py
```

The `*`s from left to right represent minute, hour, day, month, year respectively. The above setting means `monitor.py` will execute every 2 minutes, and `idle_monitor.py` will execute every 3 minutes.

```monitor.py``` searches any unmonitored resource and start to monitor them by setting `crontab` for `query_usage.py` and `get_vm_usage_class.py`.

`idle_monitor.py` monitors whether the monitored resources are idle, normal or busy, and deploys or destroys VMs according to the current usage.
>>>>>>> 1084254c7cf92e2e28b6b99f38b2fb48443d156f

#!/bin/bash

sed -i "/PBS_SERVER/ c PBS_SERVER=pbspro"  /etc/pbs.conf
sed -i "/PBS_START_MOM/ c PBS_START_MOM=1"  /etc/pbs.conf
sed -i "/client/ c \$clienthost pbspro" /var/spool/pbs/mom_priv/config
sed -i "/PBS_SCP/ c PBS_SCP=/bin/pbs_scp"  /etc/pbs.conf

echo "#!/bin/bash
scp -r -v -p \$2 \$3" > /bin/pbs_scp
chmod a+x /bin/pbs_scp

/etc/init.d/pbs start

sleep 10

/opt/pbs/bin/qmgr -c "s s job_history_enable=1"
/opt/pbs/bin/qmgr -c "s s acl_roots=root"


#!/usr/bin/env bash
#
# Copyright (C) 2019 LambdaCal Inc.
#
# The database schema initial script
#

log()
{
  echo -e "$@"
}

logerr()
{
  echo -e "Error: \033[1;31m$@\033[0m"
}

mysql_usr=$1
mysql_pwd=$2
mysql_port=$3
mysql_db_name=$4

if [ "${mysql_usr}" == "" ] || [ "${mysql_pwd}" == "" ];then
  while :
  do
    log "Please input MySQL username:"
    while :
    do
      read  mysql_usr
      if [ x"$mysql_usr" = "x" ];then
        logerr "Username cannot be null.Please input again: "
      else
        break 1
      fi
    done
    log "Password:"
    while :
    do
      read -s mysql_pwd
      if [ x"$mysql_pwd" = "x" ];then
        logerr "Password cannot be null.Please input again:"
      else
        break 1
      fi
    done
    log "Port(Press Enter if use \"3306\"):"
    read mysql_port
    if [ x"$mysql_port" = "x" ];then
      mysql_port="3306"
    fi
    log "Remote IP(Press Enter if use \"localhost\"):"
    read mysql_ip
    if [ x"$mysql_ip" = "x" ];then
      mysql_ip="localhost"
    fi
    log "Database name(Press Enter if use \"TiCP\"):"
    read mysql_db_name
    if [ x"$mysql_db_name" = "x" ];then
      mysql_db_name="TiCP"
    fi
    break 1
  done
else
  mysql_ip="localhost"
fi

# Attempt to connect database
log "Attempt to connect database."
DB_CONN="mysql -u$mysql_usr -p$mysql_pwd -h $mysql_ip -P $mysql_port -D $mysql_db_name --default-character-set=utf8"
$DB_CONN </dev/null >/dev/null 2>&1
if [ "$?" == "0" ]; then
  log "Connect successfully."
else
  logerr "Login failed, make sure input the correct username, password, port, ip and database name."
  exit 1
fi

for i in `find . -type f | grep sql$ | grep -v init.sql$ | grep -v /patch/`
do
    $DB_CONN < $i
done
log "step 1 done"
for i in `find . -type f | grep sql$ | grep init.sql$`
do
    $DB_CONN < $i
done
log "step 2 done"


#!/usr/bin/env bash
#
# The database schema initial script with table existence check

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
mysql_ip=${5:-localhost}

if [ "${mysql_usr}" == "" ] || [ "${mysql_pwd}" == "" ];then
  while :
  do
    log "MySQL username:"
    while :
    do
      read  mysql_usr
      if [ x"$mysql_usr" = "x" ];then
        logerr "Username cannot be empty"
      else
        break 1
      fi
    done

    log "Password:"
    while :
    do
      read -s mysql_pwd
      if [ x"$mysql_pwd" = "x" ];then
        logerr "Password cannot be empty"
      else
        break 1
      fi
    done

    log "Port [default 3306]:"
    read mysql_port
    if [ x"$mysql_port" = "x" ];then
      mysql_port="3306"
    fi

    log "Host [default localhost]:"
    read mysql_ip
    if [ -z "$mysql_ip" ]; then
      mysql_ip="localhost"
    fi

    log "Database name [default ticp]:"
    read mysql_db_name
    if [ x"$mysql_db_name" = "x" ];then
      mysql_db_name="ticp"
    fi
    break 1
  done
fi

TMP_CNF=$(mktemp)
cat > "$TMP_CNF" << EOF
[client]
user=$mysql_usr
password=$mysql_pwd
host=$mysql_ip
port=$mysql_port
database=$mysql_db_name
default-character-set=utf8
EOF

MYSQL_CMD=(mysql --defaults-extra-file="$TMP_CNF")

if "${MYSQL_CMD[@]}" </dev/null >/dev/null 2>&1; then
  log "Connection successful"
else
  logerr "Connection failed"
  rm -f "$TMP_CNF"
  exit 1
fi

TABLE_EXISTS=$("${MYSQL_CMD[@]}" -N -B -e "SHOW TABLES LIKE 'resource';")

if [ "$TABLE_EXISTS" != "resource" ]; then
  SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
  for i in $(find "$SCRIPT_DIR" -type f | grep '\.sql$' | grep -v '^init.sql$' | grep -v '/patch/' | grep -v 'rbac_init.sql' | grep -v 'user_init.sql'); do
    if ! "${MYSQL_CMD[@]}" < "$i"; then
      logerr "Failed to execute $i"
      rm -f "$TMP_CNF"
      exit 1
    fi
  done

  if [ -f "$SCRIPT_DIR/rbac_init.sql" ]; then
    if ! "${MYSQL_CMD[@]}" < "$SCRIPT_DIR/rbac_init.sql"; then
      logerr "Failed to execute rbac_init.sql"
      rm -f "$TMP_CNF"
      exit 1
    fi
  else
    logerr "rbac_init.sql not found at $SCRIPT_DIR/rbac_init.sql"
    rm -f "$TMP_CNF"
    exit 1
  fi

  if [ -f "$SCRIPT_DIR/user_init.sql" ]; then
    if ! "${MYSQL_CMD[@]}" < "$SCRIPT_DIR/user_init.sql"; then
      logerr "Failed to execute user_init.sql"
      rm -f "$TMP_CNF"
      exit 1
    fi
  else
    logerr "user_init.sql not found at $SCRIPT_DIR/user_init.sql"
    rm -f "$TMP_CNF"
    exit 1
  fi
  
  log "Initialization completed"
else
  log "Skipping initialization"
fi

# 清理临时配置文件
rm -f "$TMP_CNF"

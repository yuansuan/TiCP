%define WORK_DIR /workspace
%define PREFIX  /opt/yuansuan

Name: psp
Version: %{_version}
Summary: "LambdaCal PSP v%{_version}"
Release: %{_current_time}
Vendor: "LambdaCal Inc. 2016,2021."
License: "Copyright 2016-2021"
Group: Applications/Server
Distribution: Linux
URL: "http://www.yuansuan.cn"
Prefix: %{PREFIX}
AutoReqProv: no

%description
YuanSuan PSP Installer

%files
%{PREFIX}/psp/fe
%{PREFIX}/psp/bin
%{PREFIX}/psp/schema
%config(noreplace)  %{PREFIX}/psp/config
%config(noreplace)  %{PREFIX}/psp/.env

%install
mkdir -p $RPM_BUILD_ROOT%{PREFIX}/psp
/bin/cp -rf %{WORK_DIR}/dist/.  $RPM_BUILD_ROOT%{PREFIX}/psp


%pre

YS_TOP=${RPM_INSTALL_PREFIX}

# Upgrade the package
if [[ $1 -gt 1 ]] ; then

    if [[ -f "${YS_TOP}/psp/config/profile" ]] ; then
        source "${YS_TOP}/psp/config/profile"

        if [[ -f "${YS_TOP}/psp/bin/ysadmin" ]] ; then
            # Stop all services
            ysadmin stop all
        fi

    fi

fi


%post

YS_TOP=${RPM_INSTALL_PREFIX}

#--------------------------------------------------------------------------
#  Name: getFileOwner
#
#  Synopsis: getFileOwner $file
#
#  Description:
#  This function gets ownership and group of the specified file, and then remembers it.
#--------------------------------------------------------------------------
getFileOwner()
{
    FILEOWNER=`ls -l $1 | awk ' { print $3 } ' `
    FILEGROUP=`ls -l $1 | awk ' { print $4 } ' `
    export FILEOWNER FILEGROUP
}

#--------------------------------------------------------------------------
#  Name: setFileOwner
#
#  Synopsis: setFileOwner $file
#
#  Description: This function sets ownership and group of the specified file according to the ownership information
#--------------------------------------------------------------------------
setFileOwner()
{
    if [ -z $FILEOWNER ]; then
        return
    fi
    chown -f $FILEOWNER:$FILEGROUP $1
}

#--------------------------------------------------------------------------
#  Name: getFilePerms
#
#  Synopsis: getFilePerms $file
#
#  Description: This function gets mode of the specified file, and then remembers it.
#--------------------------------------------------------------------------
getFilePerms()
{
    file=`basename $1`
    dir=`dirname $1`
    FILEPERM=`find $dir -name $file -maxdepth 1 -printf "%m\n"`
    export FILEPERM
}

#--------------------------------------------------------------------------
#  Name: setFilePerms
#
#  Synopsis: setFilePerms $file
#
#  Description:
#       This function sets mode of the specified file
#       according to the mode collected earlier by
#       getFileOwner (see above).
#--------------------------------------------------------------------------
setFilePerms()
{
    if [ -z $FILEPERM ]; then
        return
    fi
    chmod -f $FILEPERM $1
}

#-----------------------------------------
# Name: replace_in_file
#
# Synopsis: replace_in_file
#
# Environment Variables:
#    None
#
# Description:
#    This function finds and replaces pattern in file
#
# Parameters:
#    $1: filename
#    $2: old pattern
#    $3: new pattern
#
# Return Value:
#      None
#------------------------------------------------
replace_in_file ()
{
    filename="$1"
    if [ -f $filename ] ; then
        getFileOwner $filename
        getFilePerms $filename

        mv $filename $filename.tmp

        sed "s?$2?$3?g" < $filename.tmp > $filename

        setFileOwner $filename
        setFilePerms $filename

        rm -f $filename.tmp
    fi
} # replace_in_file

# Replace @YS_TOP@ for profile
if [ -f "$YS_TOP/psp/config/profile" ] ; then
    replace_in_file "$YS_TOP/psp/config/profile" @YS_TOP@ "$YS_TOP"
fi

# Replace @YS_TOP@ for supervisord.conf
if [ -f "$YS_TOP/psp/config/supervisor/supervisord.conf" ] ; then
    replace_in_file "$YS_TOP/psp/config/supervisor/supervisord.conf" @YS_TOP@ "$YS_TOP"
fi

# Replace @YS_TOP@ for nginx.conf
if [ -f "$YS_TOP/psp/config/nginx/nginx.conf" ] ; then
    replace_in_file "$YS_TOP/psp/config/nginx/nginx.conf" @YS_TOP@ "$YS_TOP"
fi

# Replace @YS_TOP@ for frontend.conf
if [ -f "$YS_TOP/psp/config/nginx/frontend.conf" ] ; then
    replace_in_file "$YS_TOP/psp/config/nginx/frontend.conf" @YS_TOP@ "$YS_TOP"
fi

# Replace @YS_TOP@ for ysadmin
if [ -f "$YS_TOP/psp/bin/ysadmin" ] ; then
    replace_in_file "$YS_TOP/psp/bin/ysadmin" @YS_TOP@ "$YS_TOP"
fi

# Replace @YS_TOP@ for prometheus config files
if [ -f "$YS_TOP/psp/config/prometheus/prom/config/prometheus.yml" ] ; then
    replace_in_file "$YS_TOP/psp/config/prometheus/prom/config/prometheus.yml" @YS_TOP@ "$YS_TOP"
fi

# Replace @YS_TOP@ for kafka server.properties files
if [ -f "$YS_TOP/psp/config/kafka/server.properties" ] ; then
    replace_in_file "$YS_TOP/psp/config/kafka/server.properties" @YS_TOP@ "$YS_TOP"
fi

# Replace @YS_TOP@ for redis redis.conf file
if [ -f "$YS_TOP/psp/config/redis/redis.conf" ] ; then
    replace_in_file "$YS_TOP/psp/config/redis/redis.conf" @YS_TOP@ "$YS_TOP"
fi

# Replace @YS_TOP@ for psp service files
if [ -f "$YS_TOP/psp/config/service/psp.service" ] ; then
    replace_in_file "$YS_TOP/psp/config/service/psp.service" @YS_TOP@ "$YS_TOP"
fi
if [ -f "$YS_TOP/psp/bin/psp.service.sh" ] ; then
    replace_in_file "$YS_TOP/psp/bin/psp.service.sh" @YS_TOP@ "$YS_TOP"
fi

# Initialize the log directory if not exist
if [[ ! -d "$YS_TOP/psp/log" ]] ; then
    mkdir "$YS_TOP/psp/log"
fi

# Initialize the work and supervisor directory if not exist
if [[ ! -d "$YS_TOP/psp/work" ]] ; then
    mkdir "$YS_TOP/psp/work"
fi
if [[ ! -d "$YS_TOP/psp/work/supervisor" ]] ; then
    mkdir "$YS_TOP/psp/work/supervisor"
fi

# Initialize the boltdb directory if not exist
if [[ ! -d "$YS_TOP/psp/work/boltdb" ]] ; then
    mkdir "$YS_TOP/psp/work/boltdb"
fi

# Initialize the nginx directory if not exist
if [[ ! -d "$YS_TOP/psp/work/nginx" ]] ; then
    mkdir "$YS_TOP/psp/work/nginx"
fi

# Initialize the prometheus directory if not exist
if [[ ! -d "$YS_TOP/psp/work/prometheus" ]] ; then
    mkdir "$YS_TOP/psp/work/prometheus"
fi

# Initialize the kafka logs directory if not exist
if [[ ! -d "$YS_TOP/psp/work/kafka-logs" ]] ; then
    mkdir "$YS_TOP/psp/work/kafka-logs"
fi

# Initialize the redis db directory if not exist
if [[ ! -d "$YS_TOP/psp/work/redis" ]] ; then
    mkdir "$YS_TOP/psp/work/redis"
fi

# Get the CentOS/RHEL version number
OS_VERION=`cat /etc/redhat-release|sed -r 's/.* ([0-9]+)\..*/\1/'`

if [[ ${OS_VERION} -ge 7 ]] ; then
    # CentOS/RHEL 7.x
    # Add it into system service management
    \cp -f ${YS_TOP}/psp/config/service/psp.service /usr/lib/systemd/system/
    systemctl enable psp.service
else
    # CentOS/RHEL 6.x
    \cp -f ${YS_TOP}/psp/bin/psp.service.sh /etc/rc.d/init.d/
    cd /etc/rc.d/init.d/; chkconfig --add psp.service.sh
    chkconfig psp.service.sh on
fi

%preun

YS_TOP=${RPM_INSTALL_PREFIX}

# Remove the package
if [[ $1 -eq 0 ]] ; then

    if [[ -f "${YS_TOP}/psp/config/profile" ]] ; then
        source "${YS_TOP}/psp/config/profile"

        if [[ -f "${YS_TOP}/psp/bin/ysadmin" ]] ; then
            # Stop all services
            ${YS_TOP}/psp/bin/ysadmin stop all

            echo "All services has been stopped."
        fi
    fi

    # Stop the psp service for systemd
    if [[ -f /usr/lib/systemd/system/psp.service ]] ; then
        systemctl stop psp.service
    fi

fi

%postun

YS_TOP=${RPM_INSTALL_PREFIX}

if [[ $1 -eq 0 ]] ; then
    \rm -f /usr/lib/systemd/system/psp.service
fi

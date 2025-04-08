%define WORK_DIR /workspace
%define PREFIX  /opt/yuansuan

Name: psp-agent
Version: %{_version}
Summary: "LambdaCal Agent v%{_version}"
Release: %{_current_time}
Vendor: "LambdaCal Inc. 2016,2021."
License: "Copyright 2016-2021"
Group: Applications/Server
Distribution: Linux
URL: "http://www.yuansuan.cn"
Prefix: %{PREFIX}
AutoReqProv: no

%description
YuanSuan Agent Installer

%files
%{PREFIX}/agent

%install
mkdir -p $RPM_BUILD_ROOT%{PREFIX}/agent
mkdir -p $RPM_BUILD_ROOT%{PREFIX}/agent/bin

/bin/cp -rf %{WORK_DIR}/agent/psp_agent  $RPM_BUILD_ROOT%{PREFIX}/agent/bin
/bin/cp -rf %{WORK_DIR}/agent/config  $RPM_BUILD_ROOT%{PREFIX}/agent/


%pre

YS_TOP=${RPM_INSTALL_PREFIX}

# Upgrade the package
if [[ $1 -gt 1 ]] ; then
    if [[ -f "/usr/lib/systemd/system/psp-agent.service" ]] ; then
        systemctl stop psp-agent.service
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



# Replace @YS_TOP@ for agent service files
if [ -f "$YS_TOP/agent/config/service/psp-agent.service" ] ; then
    replace_in_file "$YS_TOP/agent/config/service/psp-agent.service" @YS_TOP@ "$YS_TOP"
fi

# CentOS/RHEL 7.x
# Add it into system service management
\cp -f ${YS_TOP}/agent/config/service/psp-agent.service /usr/lib/systemd/system/
systemctl enable psp-agent.service

%preun

YS_TOP=${RPM_INSTALL_PREFIX}

# Remove the package
if [[ $1 -eq 0 ]] ; then
    if [[ $1 -gt 1 ]] ; then
        if [[ -f "/usr/lib/systemd/system/psp-agent.service" ]] ; then
            systemctl stop psp-agent.service
        fi
    fi
fi

%postun

YS_TOP=${RPM_INSTALL_PREFIX}

if [[ $1 -eq 0 ]] ; then
    \rm -f /usr/lib/systemd/system/psp-agent.service
fi

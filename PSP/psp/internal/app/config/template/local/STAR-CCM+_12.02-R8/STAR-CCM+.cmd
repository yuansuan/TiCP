starccm="/opt/Siemens/14.02.010-R8/STAR-CCM+14.02.010-R8/star/bin/starccm+"
export CDLMD_LICENSE_FILE="29000@115.159.149.167"
hostlist=$(IFS=',' read -ra arr <<< "$YS_NODELIST"; for item in "${arr[@]}"; do echo -n "$item:$YS_CPUS_ON_NODE,"; done | sed 's/,$/\n/'); $starccm -rsh ssh -power -on $hostlist -mpi openmpi -batch run $YS_MAIN_FILE
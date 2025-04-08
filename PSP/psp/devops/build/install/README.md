
------------------------------------------------------------

0. PSP can only be installed on CentOS/RHEL 6.x and 7.x versions

   And before run it, please make sure the you can use your yum repo to download packages normally.

1. Get the PSP tar.gz package and put it to somewhere.

2. Install the unzip if not.

   > yum install -y unzip

3. Uncompress the tar.gz package and the tree structure as follows:

########### Tree
├──psp3.0_linux_x86_64
│   ├── README
│   ├── psp
│   │   ├── psp-3.0*.rpm
│   │   ├── pspinstall
│   │   ├── pspuninstall
│   │   ├── install.conf
│   ├── 3rd_party
│   │   ├── meld3.tar.gz
│   │   ├── setuptools.tar.gz
│   │   ├── supervisor.tar.gz
│   │   ├── nginx.tar.gz
│   │   ├── node.tar.gz
│   │   ├── prometheus.tar.gz
│   │   ├── kafka.tgz
│   │   ├── jre.tar.gz


4. Go into pbspro directory to install the PBSPro if necessary.

5. Go into ldap directory to install the openldap server and client if necessary.

6. Before installing the PSP, please make sure that Python are installed locally.

   If you want to use local DB, please make sure that MySQL or MariaDB server/client are installed locally.

7. Go into psp directory to install/upgade the PSP by editing the install.conf specifying the port, installation directory and DB parameters

   If you will install the c3po, please modify the nginx config in config also.

   Install the PSP
   > ./pspinstall -c install.conf

   Upgrade the PSP
   > ./pspinstall -u install.conf

   Note: If there are more than one RPM package found, use the highest version to install.


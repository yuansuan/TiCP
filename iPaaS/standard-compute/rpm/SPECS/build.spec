Name:     standard-compute
Version:  1.0
Release:  release
Summary:  standard-compute
Summary(zh_CN):  标准化计算
License:  GPLv3+
URL:      http://phabricator.intern.yuansuan.cn/source/standard-compute/
Source:  http://phabricator.intern.yuansuan.cn/source/standard-compute/

%description
standard-compute

%description -l zh_CN
标准化计算

%prep
%build
cd standard-compute
make clean
make linux-bin sccli
BUILD_ENV=dev make devops

%install

cd standard-compute

mkdir -p %{buildroot}/opt/standard-compute
cp standard-compute %{buildroot}/opt/standard-compute/standard-compute
cp devops %{buildroot}/opt/standard-compute/devops
cp sccli %{buildroot}/opt/standard-compute/sccli
mkdir -p %{buildroot}/opt/standard-compute/config

cp -r config/local.yml %{buildroot}/opt/standard-compute/config/local.yml.example
cp -r config/local_custom.yml %{buildroot}/opt/standard-compute/config/local_custom.yml.example
mkdir -p %{buildroot}/etc/systemd/system/
cp -r rpm/standard-compute.service %{buildroot}/etc/systemd/system/standard-compute.service.example

%files	
/opt/standard-compute/config/local.yml.example
/opt/standard-compute/config/local_custom.yml.example
/opt/standard-compute/devops
/opt/standard-compute/sccli
/opt/standard-compute/standard-compute
/etc/systemd/system/standard-compute.service.example
%doc
%clean


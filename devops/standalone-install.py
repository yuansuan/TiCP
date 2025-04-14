import os
import subprocess
import yaml
import time
import re

CONFIG_FILES = {
    'docker_compose': [
        'docker-compose-base.yml',
        'docker-compose-cloud-base.yml',
        'docker-compose-ipaas.yml',
        'docker-compose-psp.yml',
        'docker-compose-ipaas-hpc.yml'
    ],
    'ysadmin': 'ysadmin/config.yaml',
    'sc': 'config/ipaas/standard-compute/config.yaml'
}

def get_file_path(devops_path, file_key):
    if isinstance(CONFIG_FILES[file_key], list):
        return [os.path.join(devops_path, file) for file in CONFIG_FILES[file_key]]
    return os.path.join(devops_path, CONFIG_FILES[file_key])

def get_devops_path():
    current_path = os.getcwd()
    return current_path

def update_docker_compose_paths(devops_path,user_config):
    docker_compose_files = get_file_path(devops_path, 'docker_compose')
    volume_mounts = get_volume_mounts(devops_path, user_config)
    environment_mounts = get_environment_mounts(user_config)
    for file_path in docker_compose_files:
        if not os.path.exists(file_path):
            print(f"Warning: {file_path} 不存在，跳过")
            continue
        with open(file_path, 'r') as f:
            compose_content = yaml.safe_load(f)
        services = compose_content.get('services', {})
        for service_name, service in services.items():
            if service_name in environment_mounts:
                service['environment'] = environment_mounts[service_name]
            if service_name in volume_mounts:
                service['volumes'] = volume_mounts[service_name]
        with open(file_path, 'w') as f:
            yaml.dump(compose_content, f, default_flow_style=False, indent=2)

def get_environment_mounts(user_config):
    new_dsn = create_dsn(user_config['mysql'])
    return {
        'iamserver': [
            f"YS_MYSQL_DEFAULT_DSN={new_dsn}",
            "YS_MODE=prod"
        ],
        'mysql': [
            f"MYSQL_DATABASE={user_config['mysql']['database']}",
            f"MYSQL_USER={user_config['mysql']['user']}",
            f"MYSQL_PASSWORD={user_config['mysql']['password']}",
            f"MYSQL_ROOT_PASSWORD={user_config['mysql']['root_password']}",
        ]
    }

def get_volume_mounts(devops_path, user_config):
    return {
        'prometheus': [
            f"{devops_path}/config/base/prometheus.yml:/etc/prometheus/prometheus.yml",
            "/home/data/prometheus:/prometheus"
        ],
        'alertmanager': [
            f"{devops_path}/config/base/alertmanager.yml:/etc/alertmanager/alertmanager.yml",
            f"{devops_path}/config/base/emailalarm.tmpl:/etc/alertmanager/emailalarm.tmpl",
            "/home/data/alertmanager:/alertmanager"
        ],
        'idgen': [
            f"{devops_path}/config/base/idgen:/workspace/idgen/config",
            "/home/data/idgen/logs:/workspace/idgen/log"
        ],
        'account_bill': [
            f"{devops_path}/config/cloud-base/account_bill:/workspace/account_bill/config",
            "/home/data/account_bill/logs:/workspace/account_bill/log"
        ],
        'hydra_lcp': [
            f"{devops_path}/config/cloud-base/hydra_lcp:/workspace/hydra_lcp/config",
            "/home/data/hydra_lcp/logs:/workspace/hydra_lcp/log"
        ],
        'iamserver': [
            f"{devops_path}/config/cloud-base/iamserver:/workspace/iamserver/config",
            "/home/data/iamserver/logs:/workspace/iamserver/log"
        ],
        'nginx': [
            f"{devops_path}/config/ipaas/nginx/nginx.conf:/etc/nginx/nginx.conf",
            f"{devops_path}/config/ipaas/nginx/conf.d/root.conf:/etc/nginx/conf.d/root.conf",
            "/home/data/nginx/logs:/var/log/nginx"
        ],
        'storage': [
            f"{devops_path}/config/ipaas/storage:/workspace/storage/config",
            f"{user_config['storage']}:{user_config['storage']}",
            "/home/data/storage/logs:/workspace/storage/log"
        ],
        'standard-compute': [
            "/home/data/standard-compute/logs:/workspace/standard-compute/log",
            "/home/apps:/home/apps",
            "/tmp:/tmp",
            f"{devops_path}/config/ipaas/standard-compute:/workspace/standard-compute/config",
            "/usr:/usr",
            "/var/run/munge:/var/run/munge",
            "/etc/slurm:/etc/slurm",
            f"{user_config['storage']}:{user_config['storage']}"
        ],
        'job': [
            "/home/data/job/logs:/workspace/job/log",
            f"{devops_path}/config/ipaas/job:/workspace/job/config"
        ],
        'license': [
            f"{devops_path}/config/ipaas/license:/workspace/license/config",
            "/home/data/license/logs:/workspace/license/log"
        ],
        'psp-be': [
            f"{devops_path}/config/psp:/opt/yuansuan/psp/config",
            "/home/data/psp-be/logs:/opt/yuansuan/psp/logs/"
        ],
        'frontend': [
            f"{devops_path}/config/psp:/opt/yuansuan/psp/config"
        ]
    }

def read_user_config():
    with open('user_config.yml', 'r') as f:
        user_config = yaml.safe_load(f)
    return user_config

def update_prod_config(user_config, devops_path):
    config_path = os.path.join(devops_path, "config")
    for root, dirs, files in os.walk(config_path):
        for file in files:
            if file == 'prod.yml':
                file_path = os.path.join(root, file)
                update_mysql(file_path, user_config)

def update_mysql(file_path, user_config):
    with open(file_path, 'r') as f:
        config = yaml.safe_load(f)
    middleware=config['app']['middleware']
    if 'mysql' in middleware and 'default' in middleware['mysql']:
        if '/psp/prod.yml' in file_path:
            middleware['mysql']['default']['dsn'] = create_dsn_portal(user_config['mysql'])
        else:
            middleware['mysql']['default']['dsn'] = create_dsn(user_config['mysql'])
        config['app']['middleware']=middleware
    with open(file_path, 'w') as f:
        yaml.dump(config, f, default_flow_style=False)

def create_dsn(mysql_config):
    return f"{mysql_config['user']}:{mysql_config['password']}@tcp({mysql_config['host']}:3306)/{mysql_config['database']}?charset=utf8&parseTime=true&loc=Local"

def create_dsn_sc(mysql_config):
    return f"{mysql_config['user']}:{mysql_config['password']}@tcp({mysql_config['host']}:3306)/{mysql_config['database']}?charset=utf8&parseTime=true&loc=Local&multiStatements=true"

def create_dsn_portal(mysql_config):
    return f"{mysql_config['user']}:{mysql_config['password']}@tcp({mysql_config['host']}:3306)/ticp_portal?charset=utf8&parseTime=true&loc=Local"

def update_ysadmin_config_file(ak_info, devops_path):
    file_path = get_file_path(devops_path, 'ysadmin')
    update_ysadmin(file_path, ak_info)

def update_ysadmin(file_path, ak_info):
    access_key_id = ak_info['AccessKeyId']
    access_key_secret = ak_info['AccessKeySecret']
    ys_id = ak_info['YSId']
    with open(file_path, 'r') as f:
        config = yaml.safe_load(f)
    config['environments']['dev']['compute_access_key_id'] = access_key_id
    config['environments']['dev']['compute_access_key_secret'] = access_key_secret
    config['environments']['dev']['compute_ys_id'] = ys_id
    with open(file_path, 'w') as f:
        yaml.dump(config, f, default_flow_style=False)

def update_iam_endpoint(devops_path, user_config):
    file_path = os.path.join(devops_path, 'ysadmin', 'config.yaml')
    server_ip = user_config['server_ip']
    iam_ip = server_ip['iam']
    with open(file_path, 'r') as f:
        config = yaml.safe_load(f)
    config['environments']['dev']['iam_endpoint'] = f"http://{iam_ip}:8899"
    with open(file_path, 'w') as f:
        yaml.dump(config, f, default_flow_style=False)

def update_custom_config_files(ak_info, devops_path, user_config, userId):
    config_path = os.path.join(devops_path, "config")
    for root, dirs, files in os.walk(config_path):
        for file in files:
            if file == 'prod_custom.yml': # 替换AK
                file_path = os.path.join(root, file)
                update_custom_config_file(devops_path, file_path, ak_info, user_config,userId)

def update_custom_config_file(devops_path, file_path, ak_info, user_config, userId):
    access_key_id = ak_info['AccessKeyId']
    access_key_secret = ak_info['AccessKeySecret']
    ys_id = ak_info['YSId']
    server_ip = user_config['server_ip']
    hpc_ip = server_ip['hpc']
    iam_ip = server_ip['iam']
    with open(file_path, 'r') as f:
        config = yaml.safe_load(f)
    if 'ipaas/job' in file_path:
        config['self_ys_id'] = ys_id
        config['ak'] = access_key_id
        config['as'] = access_key_secret
        config['zones']['az-yuansuan']['hpc_endpoint'] = f"http://{hpc_ip}:8001"
        config['zones']['az-yuansuan']['storage_endpoint'] = f"http://{hpc_ip}:8001"
    elif 'ipaas/storage' in file_path:
        config['access_key_id'] = access_key_id
        config['access_key_secret'] = access_key_secret
        config['iam_server_url'] = f"http://{iam_ip}:8899"
    elif '/psp/prod_custom.yml' in file_path:
        config['openapi']['local']['settings']['app_key'] = access_key_id
        config['openapi']['local']['settings']['app_secret'] = access_key_secret
        config['openapi']['local']['settings']['user_id'] = ys_id
        config['openapi']['local']['settings']['api_endpoint'] = f"http://{iam_ip}:8899"
        config['openapi']['local']['settings']['hpc_endpoint'] = f"http://{hpc_ip}:8001"
        config['storage']['local_root_path'] = os.path.join(user_config['storage'], userId)
        config['system']['alert_manager']['alert_manager_config_path'] = os.path.join(devops_path, "config/base")
        config['system']['alert_manager']['alert_manager_url'] = f"http://{iam_ip}:9093"
    elif 'cloud-base/iamserver' in file_path:
        config['ys_api_server']['job'] = f"http://{iam_ip}:8893"
        config['ys_api_server']['account_bill'] = f"http://{iam_ip}:8891"
        config['ys_api_server']['lic_manager'] = f"http://{iam_ip}:8894"
    with open(file_path, 'w') as f:
        yaml.dump(config, f, default_flow_style=False)

def update_sc_config_file(ak_info, devops_path, user_config):
    file_path = get_file_path(devops_path, 'sc')
    update_sc(file_path, ak_info, user_config)

def update_sc(file_path, ak_info, user_config):
    access_key_id = ak_info['AccessKeyId']
    access_key_secret = ak_info['AccessKeySecret']
    ys_id = ak_info['YSId']
    server_ip = user_config['server_ip']
    hpc_ip = server_ip['hpc']
    iam_ip = server_ip['iam']
    with open(file_path, 'r') as f:
        config = yaml.safe_load(f)
    config['iam']['app_key'] = access_key_id
    config['iam']['app_secret'] = access_key_secret
    config['iam']['ys_id'] = ys_id
    config['iam']['endpoint'] = f"http://{iam_ip}:8899"
    config['hpc_storage_address'] = f"http://{hpc_ip}:8001"
    new_dsn = create_dsn_sc(user_config['mysql'])
    config['database']['dsn'] = new_dsn
    with open(file_path, 'w') as f:
        yaml.dump(config, f, default_flow_style=False)

def update_AK(user_config,devops_path):
    update_iam_endpoint(devops_path, user_config)
    userId =add_user(devops_path)
    result =create_access_keys(devops_path,userId)
    update_custom_config_files(result, devops_path, user_config, userId)
    update_ysadmin_config_file(result, devops_path)
    update_sc_config_file(result, devops_path, user_config)
    update_nginx_config_file(devops_path, user_config)

def update_nginx_config_file(devops_path, user_config):
    file_path = os.path.join(devops_path, 'config', 'ipaas', 'nginx', 'conf.d', 'root.conf')
    if not os.path.exists(file_path):
        print(f"Warning: {file_path} 不存在")
        return
    with open(file_path, 'r') as f:
        content = f.read()
    new_content = re.sub(
        r'set \$backend_ip\s+[\d\.]+;',
        f"set $backend_ip {user_config['server_ip']['hpc']};",
        content
    )
    with open(file_path, 'w') as f:
        f.write(new_content)

def add_user(devops_path):
    ysadmin_path = os.path.join(devops_path, "ysadmin", "ysadmin")
    run_path = os.path.join(devops_path, "ysadmin")
    with open(os.path.join(run_path, 'user_param.json'), 'r') as f:
        user_params = yaml.safe_load(f)
    phone_to_check = user_params.get('Phone')
    list_command = [ysadmin_path, 'iam', 'list', 'users']
    try:
        list_result = subprocess.run(list_command, capture_output=True, text=True, cwd=run_path, check=True)
        users_info = yaml.safe_load(list_result.stdout)
        for user in users_info.get('users', []):
            if user.get('phone') == phone_to_check:
                return user.get('ysid')
    except subprocess.CalledProcessError as e:
        print(f"Error executing list command: {e.stderr}")
        return None
    command = [ysadmin_path, 'iam', 'add', 'user', '-F', 'user_param.json']
    try:
        result = subprocess.run(command, capture_output=True, text=True, cwd=run_path, check=True)
        user_info = yaml.safe_load(result.stdout)
    except subprocess.CalledProcessError as e:
        print(f"Error executing command: {e.stderr}")
        return None
    return user_info['userId']

def create_access_keys(devops_path, userId):
    ysadmin_path = os.path.join(devops_path, "ysadmin", "ysadmin")
    run_path = os.path.join(devops_path, "ysadmin")
    command = [ysadmin_path, 'iam', 'add', 'secret', '-I', userId, '-T', 'YS_test']
    try:
        result = subprocess.run(command, capture_output=True, text=True, cwd=run_path, check=True)
        ak_info = yaml.safe_load(result.stdout)
    except subprocess.CalledProcessError as e:
        print(f"Error executing command: {e.stderr}")
        return None
    return ak_info

def set_directory_permissions():
    try:
        subprocess.run(['sudo', 'chown', '-R', '1001:1001', '/home/data/kafka'], check=True)
        subprocess.run(['sudo', 'chmod', '-R', '755', '/home/data/kafka'], check=True)
        subprocess.run(['sudo', 'chown', '-R', '65534:65534', '/home/data/prometheus'], check=True)
        subprocess.run(['sudo', 'chmod', '-R', '755', '/home/data/prometheus'], check=True)
    except subprocess.CalledProcessError as e:
        print(f"设置目录权限时出错: {e}")
        return False
    return True

def start_base_services(devops_path, user_config):
    # 设置目录权限
    if not set_directory_permissions():
        return

    if user_config['deploy_switch']['base']:
        subprocess.run(['docker', 'compose', '-f', os.path.join(devops_path, 'docker-compose-base.yml'), 'up', '-d'])
    time.sleep(1)
    if user_config['deploy_switch']['cloud-base']:
        subprocess.run(['docker', 'compose', '-f', os.path.join(devops_path, 'docker-compose-cloud-base.yml'), 'up', '-d'])
        printed_message = False
        while True:
            result = subprocess.run(
                ['docker', 'compose', '-f', os.path.join(devops_path, 'docker-compose-cloud-base.yml'), 'ps', '--services', '--filter', 'status=running'],
                capture_output=True, text=True
            )
            running_services = result.stdout.split()
            if 'hydra_lcp' in running_services and 'iamserver' in running_services:
                time.sleep(2)
                break
            if not printed_message:
                print("启动中...")
                printed_message = True
            time.sleep(2)
    time.sleep(5)

def start_other_services(devops_path, user_config):
    if user_config['deploy_switch']['ipaas']:
        subprocess.run(['docker','compose', '-f', os.path.join(devops_path, 'docker-compose-ipaas.yml'), 'up', '-d'])
    if user_config['deploy_switch']['ipaas-hpc']:
        subprocess.run(['docker','compose', '-f', os.path.join(devops_path, 'docker-compose-ipaas-hpc.yml'), 'up', '-d'])
    if user_config['deploy_switch']['psp']:
        subprocess.run(['docker','compose', '-f', os.path.join(devops_path, 'docker-compose-psp.yml'), 'up', '-d'])

def start_agent_services(devops_path, user_config):
    if user_config['deploy_switch']['agent']:
        try:
            # 安装 RPM 包
            subprocess.run([
                'rpm',
                '-ivh',
                '--prefix',
                '/opt/yuansuan',
                os.path.join(devops_path, 'agent', 'psp-agent-4.*.*.rpm'),
                '--nodeps'
            ], check=True)
            
            # 启动服务
            subprocess.run([   
                'systemctl',
                'start',
                'psp-agent'
            ], check=True)
        except subprocess.CalledProcessError as e:
            print(f"Error starting agent services: {e}")

def init_mysql(devops_path, user_config):
    init_db_sh = os.path.join(devops_path, "ticp_portal_sql", "init_db.sh")
    max_retries = 5
    retry_interval = 5
    for attempt in range(max_retries):
        try:
            subprocess.run([
                init_db_sh,
                "root",
                user_config['mysql']['root_password'],
                "3306",
                "ticp_portal",
                user_config['mysql']['host']
            ], check=True)
            return
        except subprocess.CalledProcessError as e:
            if attempt < max_retries - 1:
                time.sleep(retry_interval)
            else:
                print(f"数据库初始化失败: {e}")
                raise

def update_prometheus_config(devops_path, user_config):
    prometheus_file = os.path.join(devops_path, "config/base/prometheus.yml")
    if not os.path.exists(prometheus_file):
        print(f"Prometheus 配置文件 {prometheus_file} 不存在")
        return
    prometheus_config = user_config.get('prometheus')
    if not prometheus_config or 'ip' not in prometheus_config:
        print("user_config.yml 中缺少 prometheus.ip 配置")
        return
    ip = prometheus_config['ip']
    with open(prometheus_file, 'r') as f:
        config = yaml.safe_load(f)
    config['scrape_configs'] = [
        {
            'job_name': 'local pc',
            'static_configs': [
                {
                    'targets': [f'{ip}:9100'],
                    'labels': {
                        'instance': 'local-pc'
                    }
                }
            ]
        }
    ]
    with open(prometheus_file, 'w') as f:
        yaml.dump(config, f, default_flow_style=False, indent=2)

def update_init_sql(devops_path, user_config):
    init_sql_path = os.path.join(devops_path, "init.sql")
    if not os.path.exists(init_sql_path):
        print(f"Warning: {init_sql_path} 不存在")
        return
    with open(init_sql_path, 'r') as f:
        content = f.read()
    lines = content.split('\n')
    new_lines = []
    for line in lines:
        if 'CREATE DATABASE' in line and 'ticp_portal' not in line:
            line = f"CREATE DATABASE IF NOT EXISTS {user_config['mysql']['database']} DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;"
        elif 'GRANT ALL PRIVILEGES' in line:
            if 'ticp_portal' in line:
                line = f"GRANT ALL PRIVILEGES ON ticp_portal.* TO '{user_config['mysql']['user']}'@'%';"
            else:
                line = f"GRANT ALL PRIVILEGES ON {user_config['mysql']['database']}.* TO '{user_config['mysql']['user']}'@'%';"
        elif line.strip().startswith('USE '):
            line = f"USE {user_config['mysql']['database']};"
        new_lines.append(line)
    new_content = '\n'.join(new_lines)
    with open(init_sql_path, 'w') as f:
        f.write(new_content)

def main():
    devops_path = get_devops_path()
    user_config = read_user_config()
    update_init_sql(devops_path, user_config)
    update_prometheus_config(devops_path, user_config)
    update_docker_compose_paths(devops_path,user_config)
    update_prod_config(user_config, devops_path)
    start_base_services(devops_path, user_config)
    init_mysql(devops_path, user_config)
    update_AK(user_config,devops_path)
    start_other_services(devops_path, user_config)
    start_agent_services(devops_path, user_config)

if __name__ == '__main__':
    main()
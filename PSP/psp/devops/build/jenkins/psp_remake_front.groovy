#!groovy

import org.jenkinsci.plugins.workflow.steps.FlowInterruptedException

node("psp_land") {
    oper_env = env.OPER_ENV
    branch_name = env.BRANCH_NAME
    deploy_type = env.DEPLOY_TYPE
    docker_path = env.DOCKER_PATH
    repoURL = "ssh://vcssh@phabricator.intern.yuansuan.cn/diffusion/119/psp.git"
    repoPath = "psp"
    tarName = "fe.tar.gz"


    stage("代码更新") {
        try {
            echo "代码更新..."
            sh """
                    if [ -d "${repoPath}" ]; then
                        rm -rf ${repoPath}
                    fi
                    git clone ${repoURL} ${repoPath}
                    cd ${repoPath}
                    git checkout ${branch_name}
            """
            echo "代码更新完成..."
        } catch (Exception err) {
            echo err.toString()
            exit(1)
        }
    }

        stage("前端编译") {
        try {
            echo "前端编译..."
            sh"""
                export PATH=$PATH:/opt/yuansuan/3rd_party/node/bin
                npm config set registry https://registry.npm.taobao.org

                cd ${repoPath}/web/desktop
                sh build.sh

                cd ../../dist
                tar -cvzf ${tarName} fe/
            """
            echo "前端编译完成..."
        } catch (Exception err) {
            echo err.toString()
            exit(1)
        }
    }

    stage("远程部署") {
        sshPublisher(
	        publishers: [
	        	sshPublisherDesc(
		        	configName: env.REMOTE_SERVER,
		        	transfers: [
		        		sshTransfer(
		        			excludes: '',
		        			execCommand: '''

                                cd /opt/yuansuan
                                if [ -e fe.tar.gz ]; then
                                    if [ ${deploy_type} == "docker" ]; then
                                        mv fe.tar.gz ${docker_path}
                                        cd ${docker_path}
                                        rm -rf fe_bak
                                        mv fe fe_bak
                                        tar -xvzf fe.tar.gz
                                        chown -R root:root fe
                                        docker restart frontend
                                    else
                                        mv fe.tar.gz /opt/yuansuan/psp/
                                        cd /opt/yuansuan/psp/
                                        rm -rf fe_bak
                                        mv fe fe_bak
                                        tar -xvzf fe.tar.gz
                                        chown -R root:root fe
                                        source config/profile
                                        ysadmin restart frontend
                                    fi
                                    rm -f fe.tar.gz
                                fi

							''',
							execTimeout: 120000,
							flatten: false,
							makeEmptyDirs: false,
							noDefaultExcludes: false,
							patternSeparator: '[, ]+',
							remoteDirectory: '/opt/yuansuan',
							remoteDirectorySDF: false,
							removePrefix: 'psp/dist',
							sourceFiles: 'psp/dist/fe.tar.gz'
						)
					],
					usePromotionTimestamp: false,
					useWorkspaceInPromotion: false,
					verbose: true
				)
			]
		)

	}




}
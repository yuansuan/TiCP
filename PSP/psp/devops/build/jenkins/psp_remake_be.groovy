#!groovy

import org.jenkinsci.plugins.workflow.steps.FlowInterruptedException

node("psp_land") {
    oper_env = env.OPER_ENV
    branch_name = env.BRANCH_NAME
    deploy_type = env.DEPLOY_TYPE
    docker_path = env.DOCKER_PATH
    replace_config_file = env.REPLACE_CONFIG_FILE

    repoURL = "ssh://vcssh@phabricator.intern.yuansuan.cn/diffusion/119/psp.git"
    repoPath = "psp"
    goRoot = "/usr/local/go"


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

    stage("后端编译") {
        try {
            echo "后端编译..."
            sh"""
                sudo rm -rf ~/.cache/ys/psp/mod/git.yuansuan.cn
                RootPath=$PWD
                export GOROOT=${goRoot}
                export PATH="\${GOROOT}/bin:$PATH"
                export GO_PATH=\${RootPath}
                export GOPATH=\${RootPath}/gopath
                export GOPROXY=https://goproxy.cn

                cd ${repoPath}/cmd
                make docker-build-psp
                mv pspd pspd.new
                cd config
                tar -cvzf config.tar.gz prod.yml prod_custom.yml
            """
            echo "后端编译完成..."
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
                                if [ -e config.tar.gz ]; then
                                    if [ ${replace_config_file} = "false" ]; then
                                        rm config.tar.gz
                                    else
                                        chown root:root config.tar.gz
                                        if [ ${deploy_type} == "docker" ]; then
                                            mv config.tar.gz ${docker_path}/config
                                            cd ${docker_path}/config
                                        else
                                            mv config.tar.gz /opt/yuansuan/psp/config/
                                            cd /opt/yuansuan/psp/config
                                        fi
                                    fi
                                    rm -f prod.yml.bak
                                    rm -f prod_custom.yml.bak
                                    mv prod.yml prod.yml.bak
                                    mv prod_custom.yml prod_custom.yml.bak
                                    tar -xvf config.tar.gz
                                    chown root:root prod.yml
                                    chown root:root prod_custom.yml
                                    rm -f config.tar.gz
                                fi


                            ''',
                            execTimeout: 120000,
                            flatten: false,
                            makeEmptyDirs: false,
                            noDefaultExcludes: false,
                            patternSeparator: '[, ]+',
                            remoteDirectory: '/opt/yuansuan',
                            remoteDirectorySDF: false,
                            removePrefix: 'psp/cmd/config',
                            sourceFiles: 'psp/cmd/config/config.tar.gz'
                        ),
		        		sshTransfer(
		        			excludes: '',
		        			execCommand: '''
		        			    cd /opt/yuansuan
		        			    if [ -e pspd.new ]; then
		        			        chown -R root:root pspd.new
                                    chmod 755 pspd.new
                                    if [ ${deploy_type} == "docker" ]; then
                                        mv pspd.new ${docker_path}/bin
                                        cd ${docker_path}
                                        rm -f pspd.bak
                                        mv pspd pspd.bak
                                        mv pspd.new pspd
                                        if [ ${replace_config_file} = "false" ]; then
                                            docker restart psp-be
                                        fi

                                    else
                                        mv pspd.new /opt/yuansuan/psp/bin/
                                        cd /opt/yuansuan/psp/bin
                                        rm -f pspd.bak
                                        mv pspd pspd.bak
                                        mv pspd.new pspd
                                        source ../config/profile
                                        if [ ${replace_config_file} = "false" ]; then
                                            ysadmin restart psp
                                        fi
                                    fi
		        			    fi

							''',
							execTimeout: 120000,
							flatten: false,
							makeEmptyDirs: false,
							noDefaultExcludes: false,
							patternSeparator: '[, ]+',
							remoteDirectory: '/opt/yuansuan',
							remoteDirectorySDF: false,
							removePrefix: 'psp/cmd',
							sourceFiles: 'psp/cmd/pspd.new'
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
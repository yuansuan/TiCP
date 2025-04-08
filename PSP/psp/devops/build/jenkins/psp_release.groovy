#!groovy

import org.jenkinsci.plugins.workflow.steps.FlowInterruptedException

node("psp_land") {
    version = env.VERSION
    oper_env = env.OPER_ENV
    branch_name = env.BRANCH_NAME
    buildName = "psp_remakev${version}_release_build"
    repoURL = "ssh://vcssh@phabricator.intern.yuansuan.cn/diffusion/119/psp.git"
    repoPath = "/var/lib/jenkins/workspace/${buildName}/$BUILD_ID"
    workspace = "/var/lib/jenkins/workspace"
    tarPkgPath = "${workspace}/deploy_resource"
    spacePath = "/var/lib/jenkins/workspace/${buildName}"
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

    stage("前端编译") {
        try {
            echo "前端编译..."
            sh"""
                export PATH=$PATH:/opt/yuansuan/3rd_party/node/bin
                npm config set registry https://registry.npm.taobao.org

                cd ${repoPath}/web/desktop
                sh build.sh
                rm -rf ${spacePath}/dist/fe
                cp -r ${repoPath}/dist/fe ${spacePath}/dist
            """
            echo "前端编译完成..."
        } catch (Exception err) {
            echo err.toString()
            exit(1)
        }
    }

    stage("后端编译-PSP") {
        try {
            echo "后端编译..."
            sh """
            export GOROOT=${goRoot}
            export PATH="\${GOROOT}/bin:$PATH"
            export ROOT_PATH=${repoPath}
            export ROOTPATH=${repoPath}
            export GO_PATH=${spacePath}
            export GOPATH=${spacePath}/gopath
            export GOPROXY=https://goproxy.cn

            cd ${repoPath}

            if [[ "${oper_env}" == "CentOS" ]]; then
                image=registry.intern.yuansuan.cn/psp-app-build:latest
                buildInstallPackagePath=/workspace/devops/build/build_rpm.sh
            else
                image=registry.intern.yuansuan.cn/psp-app-build-ubt:latest
                buildInstallPackagePath=/workspace/devops/build/build_deb.sh
            fi

            docker run --user root -e HOME=/tmp/ -e XDG_CACHE_HOME=/tmp/.cache -e CGO_ENABLED=1 -e GOPATH=/workspace/gopath -e GOPROXY=https://goproxy.cn -v ${spacePath}/gopath/pkg:/workspace/gopath/pkg -v /root/.ssh:/root/.ssh -v \$(pwd):/workspace \${image} sh \${buildInstallPackagePath} master ${version}

            """
            echo "后端编译-PSP完成..."
        } catch (Exception err) {
            echo err.toString()
            exit(1)
        }
    }

    stage("后端编译-AGENT") {
        try {
            echo "后端编译..."
            sh """
            export GOROOT=${goRoot}
            export PATH="\${GOROOT}/bin:$PATH"
            export ROOT_PATH=${repoPath}
            export ROOTPATH=${repoPath}
            export GO_PATH=${spacePath}
            export GOPATH=${spacePath}/gopath
            export GOPROXY=https://goproxy.cn

            cd ${repoPath}

            image=registry.intern.yuansuan.cn/psp-app-build:latest
            buildInstallPackagePath=/workspace/agent/build/build_rpm.sh

            docker run --user root -e HOME=/tmp/ -e XDG_CACHE_HOME=/tmp/.cache -e CGO_ENABLED=1 -e GOPATH=/workspace/gopath -e GOPROXY=https://goproxy.cn -v ${spacePath}/gopath/pkg:/workspace/gopath/pkg -v /root/.ssh:/root/.ssh -v \$(pwd):/workspace \${image} sh \${buildInstallPackagePath} master ${version}

            """
            echo "后端编译-AGENT完成..."
        } catch (Exception err) {
            echo err.toString()
            exit(1)
        }
    }


    stage("打包上传") {
        try {
            echo "拷贝安装包到nas上,更新nas上tar包"
            sh """
                sharePath="/psp/psp"
                buildPath="\${sharePath}/build"
                build_date=`date +%Y-%m-%d`
                save_folder_num=1
                releasePkgPath="\${sharePath}/release/v${version}/\${build_date}"

                if [[ "${oper_env}" == "CentOS" ]]; then
                    pspInstallPath=${repoPath}/devops/build/install/pspinstall
                    pspUninstallPath=${repoPath}/devops/build/install/pspuninstall
                    pspInstallPackagePath=${repoPath}/dist/psp-${version}*.rpm
                    agentInstallPackagePath=${repoPath}/dist/agent-${version}*.rpm
                    tarName=psp${version}_linux_x86_64
                else
                    pspInstallPath=${repoPath}/devops/build/install/pspinstall-ubt
                    pspUninstallPath=${repoPath}/devops/build/install/pspuninstall-ubt
                    pspInstallPackagePath=${repoPath}/dist/psp-${version}*.deb
                    agentInstallPackagePath=${repoPath}/dist/agent-${version}*.deb
                    tarName=psp${version}_linux_x86_64-ubt
                fi

                mkdir -p \${releasePkgPath}

                # Build tar pacakge and move to NAS
                cd ~
                rm -rf \${tarName}
                mkdir -p  \${tarName}

                # File README
                cp -rf ${repoPath}/devops/build/install/README.md  \${tarName}/

                # Directory psp
                mkdir -p \${tarName}/psp
                # Directory agent
                mkdir -p \${tarName}/agent

                cp \${pspInstallPath}  \${tarName}/psp
                cp \${pspUninstallPath}  \${tarName}/psp
                cp ${repoPath}/devops/build/install/install.conf  \${tarName}/psp
                cp \${pspInstallPackagePath} \${tarName}/psp
                cp \${agentInstallPackagePath} \${tarName}/agent

                # Directory 3rd_party
                cp -rf ${tarPkgPath}/3rd_party  \${tarName}/

                # Compress
                tar -zcvf  \${tarName}.tar.gz  \${tarName}/

                # Delete the old psp
                sudo rm -rf \${releasePkgPath}/psp*

                # Move to NAS
                cp  \${tarName}.tar.gz \${releasePkgPath}

                sudo rm -rf  \${tarName}.tar.gz  \${tarName}
            """
            echo "拷贝安装包到nas上，更新nas上tar包 done"
        } catch (Exception err) {
            echo err.toString()
            exit(1)
        }
    }
}
#!groovy

import org.jenkinsci.plugins.workflow.steps.FlowInterruptedException

node("psp_land") {
    version = env.VERSION
    oper_env = env.OPER_ENV
    branch_name = env.BRANCH_NAME
    buildName = "psp_remake_agent_v${version}_release_build"
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

    stage("后端编译") {
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
            echo "后端编译完成..."
        } catch (Exception err) {
            echo err.toString()
            exit(1)
        }
    }



}
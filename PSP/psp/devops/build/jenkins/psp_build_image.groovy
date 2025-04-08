#!groovy

import org.jenkinsci.plugins.workflow.steps.FlowInterruptedException

node("psp_land") {
    version = env.VERSION
    branch_name = env.BRANCH_NAME
    environment = env.ENVIRONMENT
    buildName = "psp_image_v${version}_release_build"
    repoURL = "ssh://vcssh@phabricator.intern.yuansuan.cn/diffusion/119/psp.git"
    repoPath = "/var/lib/jenkins/workspace/${buildName}/$BUILD_ID"
    workspace = "/var/lib/jenkins/workspace"
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

    stage("后端编译") {
        try {
            echo "后端编译..."
            sh """
            export GOROOT=${goRoot}
            export PATH="\${GOROOT}/bin:$PATH"
            export GOPATH=${spacePath}
            export ROOT_PATH=${repoPath}
            export ROOTPATH=${repoPath}
            export GO_PATH=${spacePath}
            export GOPATH=${spacePath}/gopath
            export GOPROXY=https://goproxy.cn

            cd ${repoPath}/cmd
            make docker-build-psp
            """
            echo "后端编译完成..."
        } catch (Exception err) {
            echo err.toString()
            exit(1)
        }
    }

    stage("打包镜像") {
        try {
            echo "拷贝安装包到nas上,更新nas上tar包"
            sh """
                sharePath="/psp/psp"
                build_date=`date +%Y-%m-%d`
                releasePkgPath="\${sharePath}/release_docker/v${version}/\${build_date}"
                tarPkgPath="\${sharePath}/deploy_resource/3rd_party"
                tarName=psp_docker_v${version}.tar.gz

                cd ${repoPath}

                cp -f \${tarPkgPath}/nginx.tar.gz ./docker/package/frontend/
                cp -f \${tarPkgPath}/node.tar.gz ./docker/package/frontend/

                mkdir -p ${repoPath}/docker/package/config/
                mkdir -p ${repoPath}/docker/package/util/

                cp -rf ${repoPath}/devops/build/schema ${repoPath}/docker/package/
                cp -rf ${repoPath}/cmd/config/cert ${repoPath}/docker/package/config/
                cp -rf ${repoPath}/cmd/config/nginx ${repoPath}/docker/package/config/
                cp -rf ${repoPath}/internal/app/config/template ${repoPath}/docker/package/config/
                cp -f ${repoPath}/internal/user/config/license.yaml ${repoPath}/docker/package/config/
                cp -f ${repoPath}/cmd/config/prod.yml ${repoPath}/docker/package/config/
                cp -f ${repoPath}/cmd/config/prod_custom.yml ${repoPath}/docker/package/config/
                cp -f ${repoPath}/cmd/config/psp.conf ${repoPath}/docker/package/config/
                cp -f ${repoPath}/cmd/config/rbac_model.conf ${repoPath}/docker/package/config/
                cp -f ${repoPath}/internal/notice/config/sysconfig.yaml ${repoPath}/docker/package/config/
                cp -f ${repoPath}/devops/build/bin/encrypt ${repoPath}/docker/package/util/


                find ${repoPath}/docker/package/config/nginx -type f -exec sed -i 's/@YS_TOP@/\\/opt\\/yuansuan/g' {} +
                sed -i 's#/opt/yuansuan/psp/certs/cert.pem;#/opt/yuansuan/psp/config/cert/cert.pem;#g' ${repoPath}/docker/package/config/nginx/frontend.conf
                sed -i 's#/opt/yuansuan/psp/certs/cert.key;#/opt/yuansuan/psp/config/cert/cert.key;#g' ${repoPath}/docker/package/config/nginx/frontend.conf
                sed -i "s/127.0.0.1/psp/g" ${repoPath}/docker/package/config/nginx/frontend.conf
                chmod +x ${repoPath}/docker/package/util/encrypt

                make rm-psp-image IMAGE_TAG=v${version}
                make export-psp-image IMAGE_TAG=v${version}
                make rm-frontend-image IMAGE_TAG=v${version}
                make export-frontend-image IMAGE_TAG=v${version}

                cd docker
                rm -rf package/psp
                rm -rf package/frontend
                mv package/docker-compose-${environment}.yml package/docker-compose.yml
                rm -f package/docker-compose-*.yml

                if [ ${environment} == "dev" ];then
                    mkdir -p package/bin
                    mkdir -p package/fe
                fi

                mv package/ psp_docker/
                tar -cvzf \${tarName} psp_docker/

                mkdir -p \${releasePkgPath}
                cd \${releasePkgPath}
                rm -f \${tarName}
                cp ${repoPath}/docker/\${tarName} .

            """
            echo "打包镜像完成..."
        } catch (Exception err) {
            echo err.toString()
            exit(1)
        }
    }
}
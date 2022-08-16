@Library(['jenkins-shared']) _
pipeline {
    options {
        ansiColor('xterm')
      }
    agent {
        docker {
            args '-u root'
            image '763356858129.dkr.ecr.eu-central-1.amazonaws.com/nsof/terraform-provider:master'
            label 'pfptprod-ubuntu-bionic-py3'
            registryUrl 'https://763356858129.dkr.ecr.eu-central-1.amazonaws.com/'
            registryCredentialsId 'ecr:eu-central-1:Jenkins_Dev'
            }
    }
      environment {
          PFPTMETA_ORG_SHORTNAME = 'tftests'
          VERBOSE = 'true'
          WORKING_DIR = '/home/ubuntu/workspace/terraform-provider-pfptmeta'
        }

    stages {
        stage('Prepare') {
                steps {
                    sh "mkdir ${env.WORKING_DIR}"
                    sh "cp -r ./ ${env.WORKING_DIR}"
                }
            }
        stage("Tests") {
              steps {
                  withCredentials([
                        usernamePassword(credentialsId: 'META_API_KEY_TERRAFORM_PROVIDER', usernameVariable: 'PFPTMETA_API_KEY', passwordVariable: 'PFPTMETA_API_SECRET')]){
                      sh "cd ${env.WORKING_DIR} && make tests"
                }
            }
        }
    }
    post {
        always {
            sh "cp ${env.WORKING_DIR}/report.xml ."
            junit 'report.xml'
            sendNotifications currentBuild.result
            cleanWs(cleanWhenNotBuilt: false,
                                deleteDirs: true,
                                disableDeferredWipeout: true,
                                notFailBuild: true)
        }
    }
}
pipeline {
    options {
        ansiColor('xterm')
      }
    agent {
        docker {
            image '763356858129.dkr.ecr.eu-central-1.amazonaws.com/nsof/terraform-provider:master'
            args '-u root'
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
            sendNotifications currentBuild.result
        }
    }
}
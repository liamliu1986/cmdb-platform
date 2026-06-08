pipeline {
    agent any

    environment {
        HARBOR_URL = 'idc-test-harbor.neuedu.com'
        HARBOR_CREDENTIALS = 'harbor-credentials'
        K8S_CREDENTIALS = 'kubeconfig-credentials'
        NAMESPACE = 'idc-test'
        HELM_CHART_PATH = 'k8s/cmdb'
        HELM_RELEASE_NAME = 'cmdb-staging'
    }

    stages {
        stage('Checkout') {
            steps {
                echo 'Checking out source code...'
                checkout scm
            }
        }

        stage('Increment Version') {
            steps {
                script {
                    // Read and increment version for API
                    def apiVersionContent = readFile 'cmdb-api/version'
                    def apiVersionParts = apiVersionContent.trim().split('\\.')
                    def apiMajor = apiVersionParts[0].toInteger()
                    def apiFeature = apiVersionParts[1].toInteger()
                    def apiBugfix = apiVersionParts[2].toInteger() + 1
                    env.API_VERSION = "staging-v${apiMajor}.${apiFeature}.${apiBugfix}"
                    env.API_VERSION_FILE = "staging-v${apiMajor}.${apiFeature}.${apiBugfix}"
                    sh "echo ${apiMajor}.${apiFeature}.${apiBugfix} > cmdb-api/version"

                    // Read and increment version for UI
                    def uiVersionContent = readFile 'cmdb-ui/version'
                    def uiVersionParts = uiVersionContent.trim().split('\\.')
                    def uiMajor = uiVersionParts[0].toInteger()
                    def uiFeature = uiVersionParts[1].toInteger()
                    def uiBugfix = uiVersionParts[2].toInteger() + 1
                    env.UI_VERSION = "staging-v${uiMajor}.${uiFeature}.${uiBugfix}"
                    sh "echo ${uiMajor}.${uiFeature}.${uiBugfix} > cmdb-ui/version"

                    echo "API Version: ${env.API_VERSION}"
                    echo "UI Version: ${env.UI_VERSION}"
                }
            }
        }

        stage('Build Backend') {
            steps {
                echo 'Building cmdb-api...'
                dir('cmdb-api') {
                    sh '''
                        go mod download
                        go build -o cmdb-api .
                    '''
                }
            }
        }

        stage('Build Frontend') {
            steps {
                echo 'Building cmdb-ui...'
                dir('cmdb-ui') {
                    sh '''
                        npm install
                        npm run build
                    '''
                }
            }
        }

        stage('Docker Build & Push API') {
            steps {
                script {
                    echo "Building cmdb-api:${env.API_VERSION}"
                    dir('cmdb-api') {
                        sh """
                            docker build -t ${HARBOR_URL}/cmdb/cmdb-api:${env.API_VERSION} .
                            docker push ${HARBOR_URL}/cmdb/cmdb-api:${env.API_VERSION}
                        """
                    }
                }
            }
        }

        stage('Docker Build & Push UI') {
            steps {
                script {
                    echo "Building cmdb-ui:${env.UI_VERSION}"
                    dir('cmdb-ui') {
                        sh """
                            docker build -t ${HARBOR_URL}/cmdb/cmdb-ui:${env.UI_VERSION} .
                            docker push ${HARBOR_URL}/cmdb/cmdb-ui:${env.UI_VERSION}
                        """
                    }
                }
            }
        }

        stage('Helm Deploy to K8s') {
            steps {
                script {
                    withCredentials([file(credentialsId: env.K8S_CREDENTIALS, variable: 'KUBECONFIG_FILE')]) {
                        sh """
                            export KUBECONFIG=\${KUBECONFIG_FILE}
                            helm upgrade --install ${HELM_RELEASE_NAME} ${HELM_CHART_PATH} \\
                                -f ${HELM_CHART_PATH}/values-staging.yaml \\
                                --set api.image=${HARBOR_URL}/cmdb/cmdb-api \\
                                --set api.tag=${env.API_VERSION} \\
                                --set ui.image=${HARBOR_URL}/cmdb/cmdb-ui \\
                                --set ui.tag=${env.UI_VERSION} \\
                                --namespace ${NAMESPACE} \\
                                --wait \\
                                --timeout 5m
                        """
                    }
                }
            }
        }

        stage('Commit Version Update') {
            steps {
                script {
                    sh '''
                        git config user.email "jenkins@cmdb.local"
                        git config user.name "Jenkins"
                        git add cmdb-api/version cmdb-ui/version
                        git commit -m "chore: auto-increment version [skip ci]"
                        git push
                    '''
                }
            }
        }
    }

    post {
        success {
            echo 'Pipeline completed successfully!'
            // Add notification here (e.g., DingTalk, Enterprise WeChat)
        }
        failure {
            echo 'Pipeline failed!'
            // Add failure notification here
        }
        always {
            cleanWs()
        }
    }
}
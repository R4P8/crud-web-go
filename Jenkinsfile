pipeline {
     agent {
        docker {
            image 'rizqirafa8/golang-with-docker'
            args '-v /var/run/docker.sock:/var/run/docker.sock'
        }
    } 
    
    environment {
        DOCKER_IMAGE_NAME = 'crud-web-go'
        DOCKER_IMAGE_TAG = '1.23.0'
        DOCKER_USERNAME = 'rizqirafa8'
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Go Build') {
            steps {
                sh 'go mod tidy'
                sh 'go build -o app'
            }
        }

        stage('Test') {
            steps {
                sh 'go test ./...'
            }
        }

        stage('Build Docker Image') {
            steps {
                sh "docker build -t ${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG} ."
            }
        }

        stage('Push to Docker Hub') {
            steps {
                withCredentials([usernamePassword(credentialsId: 'docker_cred',
                                                  usernameVariable: 'DOCKERHUB_USERNAME',
                                                  passwordVariable: 'DOCKERHUB_PASSWORD')]) {
                    sh '''
                        echo "$DOCKERHUB_PASSWORD" | docker login -u "$DOCKERHUB_USERNAME" --password-stdin
                        docker tag ${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG} $DOCKERHUB_USERNAME/${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG}
                        docker push $DOCKERHUB_USERNAME/${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG}
                        docker logout
                    '''
                }
            }
        }
    }
}

pipeline {
	agent {
		label 'go-slave'
	}
    environment {
        CI = 'true'
    }
    stages {
        stage('Build') {
            steps {
                echo 'Build stage OK'
            }
        }
        stage('Test') {
            steps {
                sh 'go test'
            }
        }
        stage('Deploy') {
            steps {
                echo 'Deploy stage OK'
            }
        }
    }
}
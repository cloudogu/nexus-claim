#!groovy

node('docker') {

  repositoryOwner = 'cloudogu'
  repositoryName = 'nexus-claim'
  project = "github.com/${repositoryOwner}/${repositoryName}"
  githubCredentialsId = 'sonarqube-gh'

  stage('Checkout') {
    checkout scm
  }


  docker.image('cloudogu/golang:1.10.2').inside("--volume ${WORKSPACE}:/go/src/${project}") {

    stage('Build') {
      make 'clean'
      make 'build'
      archiveArtifacts 'target/**/*.tar.gz'
    }

    stage('Unit Test') {
      make 'unit-test'
      junit allowEmptyResults: true, testResults: 'target/*-tests.xml'
    }

    stage('Static Analysis') {
      def commitSha = sh(returnStdout: true, script: 'git rev-parse HEAD').trim()
      withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: githubCredentialsId, usernameVariable: 'USERNAME', passwordVariable: 'REVIEWDOG_GITHUB_API_TOKEN']]) {
        withEnv(["CI_PULL_REQUEST=${env.CHANGE_ID}", "CI_COMMIT=${commitSha}", "CI_REPO_OWNER=${repositoryOwner}", "CI_REPO_NAME=${repositoryName}"]) {
          make 'static-analysis'
        }
      }
    }
  }

}

String repositoryOwner
String repositoryName
String project
String githubCredentialsId

void make(goal) {
  sh "cd /go/src/${project} && make ${goal}"
}

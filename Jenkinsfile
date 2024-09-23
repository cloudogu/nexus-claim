#!groovy

@Library('github.com/cloudogu/ces-build-lib@2.4.0')
import com.cloudogu.ces.cesbuildlib.*

goVersion = "1.23"

node('docker') {

  repositoryOwner = 'cloudogu'
  repositoryName = 'nexus-claim'
  project = "github.com/${repositoryOwner}/${repositoryName}"
  githubCredentialsId = 'sonarqube-gh'

  stage('Checkout') {
    checkout scm
  }

  new Docker(this)
    .image("golang:${goVersion}")
    .mountJenkinsUser()
    .inside("--volume ${WORKSPACE}:/go/src/${project} -w /go/src/${project}")  {

    withCredentials([
        [$class: 'UsernamePasswordMultiBinding', credentialsId: 'sonarqube-gh', usernameVariable: 'USERNAME', passwordVariable: 'REVIEWDOG_GITHUB_API_TOKEN']
      ]) {
        sh 'git config --global url."https://$USERNAME:$REVIEWDOG_GITHUB_API_TOKEN@github.com".insteadOf "https://github.com"'
      }

      stage('Build') {
        make 'clean'
        make ''
        archiveArtifacts 'target/**/*.tar.gz'
      }

      stage('Unit Test') {
        make 'unit-test'
        junit allowEmptyResults: true, testResults: 'target/*-tests.xml'
      }

      stage('Static Analysis') {
        def commitSha = sh(returnStdout: true, script: 'git rev-parse HEAD').trim()
        withCredentials([
          [$class: 'UsernamePasswordMultiBinding', credentialsId: githubCredentialsId, usernameVariable: 'USERNAME', passwordVariable: 'REVIEWDOG_GITHUB_API_TOKEN']
          ]) {
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

#!groovy

@Library('github.com/cloudogu/ces-build-lib@2.4.0')
import com.cloudogu.ces.cesbuildlib.*

goVersion = "1.23"

// Configuration of branches
productionReleaseBranch = "master"
developmentBranch = "develop"
currentBranch = "${env.BRANCH_NAME}"

node('docker') {

  repositoryOwner = 'cloudogu'
  repositoryName = 'nexus-claim'
  project = "github.com/${repositoryOwner}/${repositoryName}"
  githubCredentialsId = 'sonarqube-gh'

  stage('Checkout') {
    checkout scm
    make 'clean'
  }

  new Docker(this)
    .image("golang:${goVersion}")
    .mountJenkinsUser()
    .inside("--volume ${WORKSPACE}:/go/src/${project} -w /go/src/${project}")  {

      stage('Build') {
        make 'clean'
        make ''
        archiveArtifacts 'target/**/*.tar.gz'
      }

      stage('Unit Test') {
        make 'unit-test'
        junit allowEmptyResults: true, testResults: 'target/*-tests.xml'
      }

      stage("Review dog analysis") {
        stageStaticAnalysisReviewDog()
      }
  }

  stage('SonarQube') {
    stageStaticAnalysisSonarQube()
  }

}

String repositoryOwner
String repositoryName
String project
String githubCredentialsId

void make(goal) {
  sh "cd /go/src/${project} && make ${goal}"
}

void gitWithCredentials(String command) {
    withCredentials([usernamePassword(credentialsId: 'cesmarvin', usernameVariable: 'GIT_AUTH_USR', passwordVariable: 'GIT_AUTH_PSW')]) {
        sh(
            script: "git -c credential.helper=\"!f() { echo username='\$GIT_AUTH_USR'; echo password='\$GIT_AUTH_PSW'; }; f\" " + command,
            returnStdout: true
        )
    }
}

void stageStaticAnalysisReviewDog() {
    def commitSha = sh(returnStdout: true, script: 'git rev-parse HEAD').trim()

    withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: 'sonarqube-gh', usernameVariable: 'USERNAME', passwordVariable: 'REVIEWDOG_GITHUB_API_TOKEN']]) {
        withEnv(["CI_PULL_REQUEST=${env.CHANGE_ID}", "CI_COMMIT=${commitSha}", "CI_REPO_OWNER=${repositoryOwner}", "CI_REPO_NAME=${repositoryName}"]) {
            make 'static-analysis-ci'
        }
    }
}

void stageStaticAnalysisSonarQube() {
    def scannerHome = tool name: 'sonar-scanner', type: 'hudson.plugins.sonar.SonarRunnerInstallation'
    withSonarQubeEnv {
        sh "git config 'remote.origin.fetch' '+refs/heads/*:refs/remotes/origin/*'"
        gitWithCredentials("fetch --all")

        if (currentBranch == productionReleaseBranch) {
            echo "This branch has been detected as the production branch."
            sh "${scannerHome}/bin/sonar-scanner -Dsonar.branch.name=${env.BRANCH_NAME}"
        } else if (currentBranch == developmentBranch) {
            echo "This branch has been detected as the development branch."
            sh "${scannerHome}/bin/sonar-scanner -Dsonar.branch.name=${env.BRANCH_NAME}"
        } else if (env.CHANGE_TARGET) {
            echo "This branch has been detected as a pull request."
            sh "${scannerHome}/bin/sonar-scanner -Dsonar.pullrequest.key=${env.CHANGE_ID} -Dsonar.pullrequest.branch=${env.CHANGE_BRANCH} -Dsonar.pullrequest.base=${developmentBranch}"
        } else if (currentBranch.startsWith("feature/")) {
            echo "This branch has been detected as a feature branch."
            sh "${scannerHome}/bin/sonar-scanner -Dsonar.branch.name=${env.BRANCH_NAME}"
        } else {
            echo "This branch has been detected as a miscellaneous branch."
            sh "${scannerHome}/bin/sonar-scanner -Dsonar.branch.name=${env.BRANCH_NAME} "
        }
    }
    timeout(time: 2, unit: 'MINUTES') { // Needed when there is no webhook for example
        def qGate = waitForQualityGate()
        if (qGate.status != 'OK') {
            unstable("Pipeline unstable due to SonarQube quality gate failure")
        }
    }
}

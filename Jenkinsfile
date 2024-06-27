#!/usr/bin/env groovy

env.GIT_DEFAULT_BRANCH = 'ng/tlsmonitor-chart-auto-release'
final Boolean isDefaultBranch = env.BRANCH_IS_PRIMARY == 'true'
// var to handle deployment from PR
final String DEPLOY_PR_COMMENT = '/deploy-PR'
String dockerTagNameForPR = ''
Boolean isDeployingPR = false

if (!isDefaultBranch) stopOlderBuilds()

pipeline {
    agent { label 'standard' }

    triggers {
        issueCommentTrigger(DEPLOY_PR_COMMENT)
    }

    options {
        buildDiscarder(logRotator(numToKeepStr: '20'))
        timeout(time: 30, unit: 'MINUTES')
        parallelsAlwaysFailFast()
        ansiColor('xterm')
    }

    parameters {
        booleanParam(
            name: 'PUSH',
            defaultValue: isDefaultBranch,
            description: 'Whether or not to push the built container image to the registry',
        )
    }

    environment {
        GIT_SSH_COMMAND = 'ssh -o StrictHostKeyChecking=no'
        GOPRIVATE = 'github.com/blablacar'
    }

    stages {
        stage('Check if triggered') {
            steps {
                script {
                    def timestamp = sh(script: "date '+%Y%m%d%H%M%S'", returnStdout: true)
                    def triggerCause = currentBuild.rawBuild.getCause(org.jenkinsci.plugins.pipeline.github.trigger.IssueCommentCause)
                    if (triggerCause && triggerCause.comment == DEPLOY_PR_COMMENT) {
                        echo 'Build triggered from PR comment'
                        isDeployingPR = true
                        // generate docker tag
                        dockerTagNameForPR = "${BRANCH_NAME}-${BUILD_ID}.${timestamp}".replace('/', '_').toLowerCase().trim()
                        echo "Docker tag will be ${dockerTagNameForPR}"
                    }
                }
            }
        }

        stage('Build app & lint chart') {
            parallel {
                stage('Lint Helm chart') {
                    steps {
                        dir("_infra") {
                            sh 'bbc charts lint'
                        }
                    }
                }

                stage('Build app & image') {
                    steps {
                        // if CUSTOM_TAG is empty, use the bbc versionning
                        // dockerTagNameForPR is set only is the case of isDeployingPR
                        sh "CUSTOM_TAG=${dockerTagNameForPR} REGISTRY_TARGET=eu.gcr.io/bbc-registry bbc dockerfile build"
                    }
                }

                stage('Security scan') {
                    steps {
                        securityScan(team: 'core-infrastructure')
                    }
                }
            }
        }

        stage('Push') {
            when {
                expression { params.PUSH || isDeployingPR }
            }
            parallel {
                stage('Push chart') {
                    steps {
                        dir("_infra") {
                          sh 'bbc charts release'
                        }
                    }
                }
                stage('Push docker') {
                    steps {
                        sh 'bbc dockerfile push'
                    }
                }
            }
        }
    }
}

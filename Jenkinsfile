pipeline {
    agent any

    environment {
        DOCKER_CREDS = credentials('dockerhub-login')
        DOCKER_USER  = 'kingdil'
        KUBECONFIG   = "C:\\Users\\Rifky Fadhillah A\\.kube\\config"
    }

    stages {
        stage('1. Checkout repo') {
            steps {
                echo 'Mengambil kode dari GitHub...'
                git branch: 'main', url: 'https://github.com/RIFKY26/SatSet.git'
            }
        }

        stage('Persiapan: Login Docker Hub') {
            steps {
                // UBAH 1: sh menjadi bat, dan $ menjadi % untuk Windows
                bat 'docker login -u %DOCKER_CREDS_USR% -p %DOCKER_CREDS_PSW%'
            }
        }

        stage('Proses Microservices (Unit Test -> Vet -> Build -> Func Test -> Push)') {
            steps {
                script {
                    def services = [
                        'auth_service', 'driver-service', 'location-service', 
                        'matching-service', 'notification', 'order-service', 
                        'payment-service', 'promo_service', 'rating', 'user_service'
                    ]

                    for (int i = 0; i < services.size(); i++) {
                        def svc = services[i]
                        def imageName = "${DOCKER_USER}/satset-${svc.replace('_', '-')}:latest"

                        echo "=================================================="
                        echo "MEMPROSES SERVICE: ${svc.toUpperCase()}"
                        echo "=================================================="

                        dir(svc) {
                            // UBAH 2: Semua sh diubah menjadi bat
                            echo "--> Menjalankan Unit Tests..."
                            bat returnStatus: true, script: 'go test -v ./...'

                            echo "--> Menjalankan Lint/Vet..."
                            bat returnStatus: true, script: 'go vet ./...'

                            echo "--> Merakit Docker Image..."
                            bat "docker build -t ${imageName} ."

                            echo "--> Menjalankan Functional Tests..."
                            bat returnStatus: true, script: 'go test -v ./... -tags=functional'

                            echo "--> Mengunggah ke Docker Hub..."
                            bat "docker push ${imageName}"
                        }
                    }
                }
            }
        }

        stage('7. Deploy di kubernetes') {
            steps {
                echo "Menerapkan konfigurasi Kubernetes..."
                // UBAH 3: sh menjadi bat
                bat returnStatus: true, script: '''
                    kubectl apply -f ./auth_service/auth-k8s.yaml
                    kubectl apply -f ./driver-service/driver-k8s.yaml
                    kubectl apply -f ./location-service/location-k8s.yaml
                    kubectl apply -f ./matching-service/matching-k8s.yaml
                    kubectl apply -f ./order-service/order-k8s.yaml
                    kubectl apply -f ./payment-service/payment-k8s.yaml
                    kubectl apply -f ./promo-service/promo-k8s.yaml
                    kubectl apply -f ./rating/rating-k8s.yaml
                    kubectl apply -f ./user_service/user-k8s.yaml
                '''
            }
        }

        stage('8. Verify') {
            steps {
                echo "Memverifikasi Pods yang berjalan di Kubernetes..."
                // UBAH 4: sh menjadi bat
                bat 'kubectl get pods'
            }
        }
    }

    post {
        always {
            echo "Pipeline SatSet selesai dieksekusi!"
            // UBAH 5: sh menjadi bat
            bat returnStatus: true, script: 'docker system prune -f'
        }
    }
}
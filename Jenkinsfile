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
                            // 2. Unit Tests
                            echo "--> Menjalankan Unit Tests..."
                            // Menggunakan returnStatus: true karena dosen bilang test akan failed
                            bat returnStatus: true, script: 'go test -v ./...'

                            // 3. Lint/Vet
                            echo "--> Menjalankan Lint/Vet..."
                            bat returnStatus: true, script: 'go vet ./...'

                            // 4. Build Image (lokal)
                            echo "--> Merakit Docker Image..."
                            bat "docker build -t ${imageName} ."

                            // 5. Functional Tests
                            echo "--> Menjalankan Functional Tests..."
                            // Simulasi pemanggilan functional test
                            bat returnStatus: true, script: 'go test -v ./... -tags=functional'

                            // 6. Push image
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
                // Menjalankan get pods langsung di dalam Jenkins sesuai permintaan dosen
                bat 'kubectl get pods'
            }
        }
    }

    post {
        always {
            echo "Pipeline SatSet selesai dieksekusi!"
            bat returnStatus: true, script: 'docker system prune -f'
        }
    }
}
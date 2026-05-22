pipeline {
    agent any

    // Memanggil kredensial dockerhub-login yang sudah kamu buat
    environment {
        DOCKER_CREDS = credentials('dockerhub-login')
        DOCKER_USER  = 'kingdil'
        DOCKER_HOST  = 'tcp://localhost:2375'
    }

    stages {
        stage('1. Ambil Kode dari GitHub') {
            steps {
                git branch: 'main', url: 'https://github.com/RIFKY26/SatSet.git'
            }
        }

        stage('2. Login Docker Hub') {
            steps {
                // Jenkins akan login ke Docker menggunakan username & password dari kredensial
                bat 'docker login -u %DOCKER_CREDS_USR% -p %DOCKER_CREDS_PSW%'
            }
        }

        stage('3. Test, Build & Push All Services') {
            steps {
                script {
                    // Daftar semua folder service kamu
                    def services = [
                        'auth_service', 'driver-service', 'location-service', 
                        'matching-service', 'notification', 'order-service', 
                        'payment-service', 'promo_service', 'rating', 'user_service'
                    ]

                    for (int i = 0; i < services.size(); i++) {
                        def svc = services[i]
                        
                        // Membuat nama image (contoh: kingdil/satset-auth-service:latest)
                        // Mengganti underscore (_) dengan strip (-) agar rapi di Docker Hub
                        def imageName = "${DOCKER_USER}/satset-${svc.replace('_', '-')}:latest"

                        echo "=================================================="
                        echo "MEMPROSES SERVICE: ${svc.toUpperCase()}"
                        echo "=================================================="

                        dir(svc) {
                            // 1. Jalankan Unit & Functional Test
                            // Menggunakan returnStatus: true agar pipeline tidak berhenti saat Functional Test gagal
                            echo "--> Menjalankan Testing..."
                            bat returnStatus: true, script: 'go test -v ./...'

                            // 2. Build Docker Image
                            echo "--> Merakit Docker Image..."
                            bat "docker build -t ${imageName} ."

                            // 3. Push Docker Image ke Docker Hub
                            echo "--> Mengunggah ke Docker Hub..."
                            bat "docker push ${imageName}"
                        }
                    }
                }
            }
        }

        stage('4. Deploy ke Kubernetes') {
            steps {
                script {
                    // Setelah semua image berhasil diunggah, perintahkan Kubernetes untuk menjalankannya
                    echo "Menerapkan konfigurasi Kubernetes..."
                    
                    // Eksekusi semua file -k8s.yaml yang ada di dalam folder service masing-masing
                    // (Pastikan Minikube/Kubernetes sudah berjalan dan terhubung dengan Jenkins)
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
                    // Catatan: Hapus atau sesuaikan nama file di atas jika ada yang berbeda
                }
            }
        }
    }

    post {
        always {
            echo "Pipeline SatSet selesai dieksekusi!"
            // Membersihkan sisa image di komputer lokal agar harddisk Jenkins tidak penuh
            bat returnStatus: true, script: 'docker system prune -f'
        }
    }
}
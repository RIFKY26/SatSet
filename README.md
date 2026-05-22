# SatSet Backend Microservices

Proyek ini adalah implementasi arsitektur *Microservices* untuk aplikasi **SatSet**, yang di-deploy menggunakan sistem kontainerisasi dan orkestrasi modern. Proyek ini dibangun untuk memenuhi tugas implementasi *Cloud Computing* dan *Continuous Integration/Continuous Deployment* (CI/CD).

## Teknologi yang Digunakan
* **Bahasa Pemrograman:** Golang (Go)
* **Containerization:** Docker & Docker Hub
* **Orchestration:** Kubernetes (Minikube)
* **CI/CD Pipeline:** Jenkins
* **Version Control:** Git & GitHub

## Arsitektur Layanan (Microservices)
Sistem ini terdiri dari 10 layanan independen yang berjalan secara terpisah di dalam *pod* Kubernetes:
1. `auth-service`
2. `driver-service`
3. `location-service`
4. `matching-service`
5. `notification-service`
6. `order-service`
7. `payment-service`
8. `promo-service`
9. `rating-service`
10. `user-service`

## CI/CD Pipeline (Jenkins)
Proyek ini mengadopsi otomatisasi penuh menggunakan Jenkins (metode *Monorepo Pipeline*). Script otomatisasi dapat dilihat pada file `Jenkinsfile` di *root directory*. Alur kerjanya meliputi:
1. **Pull Repository:** Mengambil kode terbaru dari GitHub.
2. **Testing & Build:** Melakukan pengujian dan merakit 10 *Docker Image* secara iteratif.
3. **Push to Registry:** Mengunggah *image* ke Docker Hub (`kingdil/satset-*`).
4. **Deploy to K8s:** Mengaplikasikan konfigurasi `.yaml` untuk melakukan *rolling update* ke dalam *cluster* Kubernetes lokal (Minikube).

---

## Bukti Eksekusi & Pengujian (Proof of Work)

Berikut adalah dokumentasi keberhasilan *build* dan *deployment* dari proyek ini:

### 1. Keberhasilan CI/CD Pipeline (Jenkins)
<img width="1879" height="910" alt="image" src="https://github.com/user-attachments/assets/217d9afe-2fa6-4ae9-81d9-aca82e291ba8" />

### 2. Ketersediaan Image di Docker Hub
<img width="1472" height="496" alt="image" src="https://github.com/user-attachments/assets/9722c8f5-3148-4dfe-86f1-9573c33afc06" />


### 3. Kubernetes Pods Status (Running)
<img width="1081" height="491" alt="image" src="https://github.com/user-attachments/assets/72080f4e-1e42-4264-bbd6-80f4b5c634dc" />


### 4. Uji Coba Port-Forwarding (Test Drive)
Aplikasi telah diuji coba dengan membuka terowongan ke *user-service* (`kubectl port-forward svc/user-service 8085:80`).

<img width="517" height="306" alt="image" src="https://github.com/user-attachments/assets/63dcaccb-4a33-42f3-ad3f-4eec62591fc8" />

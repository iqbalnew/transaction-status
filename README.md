# qcash-template-service

**QCash Template Service** sebuah template dengan framework GRPC yang sudah disederhanakan agar lebih mudah digunakan dalam pengembangan, template ini dikhususkan untuk service (~bukan engine~) atau backend yang langsung berhubungan dengan frontend

Template ini masih dapat berkembang lagi agar **QCash** memiliki standard code yang sama untuk seluruh feature, dan memudahkan backend developer dalam kolaborasi

## Installation
Untuk menjalankan aplikasi **QCash Template Service** gunakan command dibawah ini:
- Pastikan sudah mengikuti step 6 untuk instalasi GRPC: https://confluence.bri.co.id/x/KWKEE
- Rubah atau sesuaikan file .proto pada path /proto
- Sesuaikan isi **www\swagger-ui\index.html** dengan api yang akan tersedia pada apliksai
- Execute pada file generate.bat (Windows) atau generate.sh (Linux/WSL)
- Rename folder **example.vscode** menjadi **.vscode**
- Sesuaikan seluruh environment variable pada **.env**
- Apabila sudah selesai cleanup atau code jangan lupa jalankan `go mod tidy` untuk mengupdate library atau dependency yang ada pada project
- Gunakan command: `make run` untuk menjalankan aplikasi, atau `F5` untuk debugging
- Untuk mengecek apakah aplikasi sudah berjalan atau belum dapet mengakses dibrowser dengan path: `<host>:<port>/v1/template/docs/`
  `<host>` : Apabila dinyalakan pada laptop bisa gunakan **localhost**
  `<port>` : Gunakan port **HTTP** bukan **GRPC** (--port2 pada command run)
- Untuk command lainnya dapat dilihat pada **Makefile**

## Tips & Tricks
Untuk penggunaan template ini ada beberapa tips dan tricks yang dapat diikuti:
- Apabila ingin menambahkan atau membuat api, table, model, dan kebutuhan fungsi lainnya harus menggunakan file .proto yang akan ter-generate otomatis
- Pisahkan file .proto sesuai dengan kebutuhannya, misal: **_api**, **_core**, **_db**, dan **_payload**
- **swagger.json** akan tergenerate otomatis sesuai dengan file **_api.proto**
- Tidak perlu mengubah apapun pada folder **www** (kecuali path **www\swagger-ui\index.html**), karena hanya digunakan oleh swagger, datanya akan based on **swagger.json**
- Source code ada pada folder **server**, berikut penjelasan **package** atau
    - path **/api** digunakan sebagai logic utama / http handler utama aplikasi
        * **/grpc** digunakan untuk api dengan framework grpc
            - **/handler** digunakan untuk logic-logic dari service lain (ex: auth, workflow, dll)
        * **/http** digunakan untuk api http native
    - path **/config** digunakan untuk configurasi dari environment variable
    - path **/constant** digunakan untuk value-value yang reusable, hardcode, atau key
    - path **/db** digunakan untuk per-database-an (postgresql, mongo, redis, dkk), jadi seluruh repository ada disana, pisahkan query setiap table atau setiap menu pada feature dengan file sendiri (ex: menu.go, inquiry.go)
    - path **/interceptors** digunakan untuk embed middleware grpc mulai dari Auth sampai Logging
    - path **/lib** digunakan apabila memerlukan 3rd party apps atau library
    - path **/model** digunakan apabila memerlukan struct atau model database
    - path **/pb** hasil generate dari .proto jangan pernah mengedit langsung file pb tersebut
    - path **/service** digunakan apabila kita memerlukan koneksi ke service server GRPC yang lain
        * **/stubs** digunakan untuk menyimpan file protobuf dari service lain dari folder **/pb**, selalu update file protobuf dari service-service lainnya apabila ada pengembangan (confirm ke devops untuk branch terbaru)
    - path **/utils** digunakan apabila kita memiliki function utility atau reusable function
- **main.go** sesuaikan value-value configurasi untuk pengembangan fitur
- Selalu update **.env.example** agar Squad Lead dapat menambahkan value-value environment variable pada Helm Chart pada repository deployment [addons-deployment](https://bitbucket.bri.co.id/scm/addons-ops/addons-deployment.git)

## Tech

 **QCash Template Service** menggunakan beberapa open source project dan aplikasi:

- [GRPC]
- [GRPC Gateway]
- [Mux]
- [Swagger]
- [Makefile]
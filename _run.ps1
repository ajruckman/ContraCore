.\_setup.ps1

cd ..\..\coredns\coredns\
go build -o ..\..\ajruckman\ContraCore\coredns.exe
cd ..\..\ajruckman\ContraCore\

.\coredns.exe

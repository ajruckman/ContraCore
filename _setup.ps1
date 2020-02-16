Write-Output "Set options..."
$env:GOPRIVATE="github.com/ajruckman/ContraCore"
$env:GOOS="windows"
$env:GOARCH="amd64"
Write-Output ""

#Write-Output "Update go-autorest..."
#cd "$(go env GOPATH)\src\github.com\coredns\coredns"
#go get github.com/Azure/go-autorest@latest
#go get github.com/Azure/go-autorest/autorest/adal@latest
#Write-Output ""

Write-Output "Copy module..."
Remove-Item -Force -Recurse "$(go env GOPATH)\pkg\mod\github.com\ajruckman\!contra!core@v0.0.1"
New-Item -Type Directory "$(go env GOPATH)\pkg\mod\github.com\ajruckman\!contra!core@v0.0.1" | Out-Null
Copy-Item "$(go env GOPATH)\src\github.com\ajruckman\ContraCore\*" -Exclude ('.git', '.idea') -Force -Recurse "$(go env GOPATH)\pkg\mod\github.com\ajruckman\!contra!core@v0.0.1"
Write-Output ""

cd "$(go env GOPATH)\src\github.com\ajruckman\ContraCore"

#go build -o coredns.exe

#$ErrorActionPreference = 'SilentlyContinue'

#& .\coredns.exe -conf "$(go env GOPATH)\src\github.com\ajruckman\ContraCore\Corefile"

#Stop-Process -Id $PID -Force

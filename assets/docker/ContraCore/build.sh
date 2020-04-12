cd /opt/

wget https://github.com/coredns/coredns/archive/v1.6.9.tar.gz -O coredns.tgz
tar zxvf coredns.tgz
mv ./coredns-1.6.9/ ./coredns/

cd ./coredns/

sed -i '/forward:forward/i \
contracore:github.com/ajruckman/ContraCore/pkg/plugin' plugin.cfg

cd ..

git clone https://github.com/ajruckman/ContraCore ./contracore/

cd ./contracore/

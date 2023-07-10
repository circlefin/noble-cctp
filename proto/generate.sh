cd proto
buf generate
cd ..

cp -r github.com/circlefin/noble-cctp/* ./
rm -rf github.com

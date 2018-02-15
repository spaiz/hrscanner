### HR-Scanner
A simple domain names scanner.
The app knows to scan websites and extract X-Recruiting header,
if presents.

### Usage
```bash
mkdir -p ${GOPATH}/src/github.com/spaiz/
cd ${GOPATH}/src/github.com/spaiz/
git clone git@github.com:spaiz/hrscanner.git .
cd hrscanner
```

Install dependency manager use in the project

```bash
go get -u github.com/kardianos/govendor
```

Install dependencies:

```bash
govendor sync
```

Install app locally:

```bash
go install
```

Use this helper to build statically compiled binary on the MacOS (DOcker needed)
```bash
./bin/build
```

To run the app:

```bash
cd ${GOPATH}/src/github.com/spaiz/hrscanner/data/
unzip dns-servers.txt.zip && rm dns-servers.txt.zip
unzip uniq-domains.txt.zip && uniq-domains.txt.zip
cd ${GOPATH}/src/github.com/spaiz/hrscanner/
hrscanner -workers=500
```

U can previously compiled binary inside the Docker locally or upload the image to AWS ECR and run it as a Task

```bash
docker build -t hrscanner -f Dockerfile .
docker run -it hrscanner
# or with arguments
docker run -it hrscanner hrscanner -workers=500
```

### Important
Increase ulimit on your system if you want to scan huge number of websites

U can always download updated list of DNS servers from:

https://public-dns.info//nameservers.csv
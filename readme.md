# Requirements

- go 1.24.4
- air(live reload)
  go install github.com/air-verse/air@latest
- make
  sudo apt install make
- mockgen
  go install go.uber.org/mock/mockgen@latest
- migrate
  go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

mocks folder filled by mockgen

- check if mockgen installed
  mockgen -version

- if it didn't
  go install go.uber.org/mock/mockgen@latest

note: make sure $GOPATH/bin is in your $PATH, which is usually $HOME/go/bin
  if it's not
  echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.bashrc && source ~/.bashrc

# HOW TO TEST

- testing repository require postgresql
  make compose-test-up

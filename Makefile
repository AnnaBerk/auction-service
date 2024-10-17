
BUFBUILD := docker run --rm -v $(CURDIR):/workspace -w /workspace local/buf:latest

PROTO_DIR=api

.PHONY: proto bufbuild


proto:: bufbuild
	$(BUFBUILD) generate --template buf.gen.yaml -v $(PROTO_DIR)

bufbuild:
	$(BUFBUILD) lint --path $(PROTO_DIR)


mfd-xml:
	@mfd-generator xml -c "postgres://auction_user:auction_password@localhost:5432/postgres?sslmode=disable" -m ./docs/auction.mfd -n "auction:auction,user,bid,lot"

mfd-model:
	@mfd-generator model -m ./docs/auction.mfd -p db -o ./internal/infrastructure/repo


	genna model \
      --conn="postgres://auction_user:auction_password@localhost:5432/postgres?sslmode=disable"\
      --output=./internal/domain/models \
      --package=models


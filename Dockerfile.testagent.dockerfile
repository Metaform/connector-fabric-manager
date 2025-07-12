#  Copyright (c) 2025 Metaform Systems, Inc
#
#  This program and the accompanying materials are made available under the
#  terms of the Apache License, Version 2.0 which is available at
#  https://www.apache.org/licenses/LICENSE-2.0
#
#  SPDX-License-Identifier: Apache-2.0
#
#  Contributors:
#       Metaform Systems, Inc. - initial API and implementation
#

FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .

# Build the server binary
RUN CGO_ENABLED=0 go build -o bin/testagent ./pmanager/agent/testagent/main.go

# Production stage
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /app/bin/testagent /testagent

ENTRYPOINT ["/testagent"]
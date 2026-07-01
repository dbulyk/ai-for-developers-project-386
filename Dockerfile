# Stage 1: build the React frontend.
FROM node:22-alpine AS frontend
WORKDIR /repo
COPY frontend/package.json frontend/package-lock.json ./frontend/
RUN npm ci --prefix frontend --legacy-peer-deps
COPY frontend/ ./frontend/
RUN npm run build --prefix frontend

# Stage 2: build the Go backend with the SPA embedded.
FROM golang:1.25-alpine AS backend
WORKDIR /repo
COPY backend/go.mod backend/go.sum ./backend/
RUN cd backend && go mod download
COPY backend/ ./backend/
COPY --from=frontend /repo/frontend/dist ./frontend/dist
RUN mkdir -p backend/internal/assets/dist && \
    cp -R frontend/dist/. backend/internal/assets/dist/ && \
    mkdir -p backend/internal/assets/dist/assets && \
    touch backend/internal/assets/dist/assets/.gitkeep && \
    cd backend && \
    CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags='-s -w' -o /out/server ./cmd/server

# Stage 3: minimal runtime image.
FROM gcr.io/distroless/static:nonroot
COPY --from=backend /out/server /server
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/server"]
